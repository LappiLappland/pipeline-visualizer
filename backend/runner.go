package main

import (
	"context"
	"log"
	"math/rand/v2"
	"sync"
	"time"

	"gorm.io/gorm"
)

type RunnerJob struct {
	jobWorker    *JobWorker
	dependencies []*RunnerJob
	dependents   []*RunnerJob
	status       statusOption
	cancel       context.CancelFunc
	plannedBy    *User
}

type Runner struct {
	pipeline        PipelineWorker
	mu              sync.Mutex
	allJobs         map[uint]*RunnerJob
	plannedJobs     map[uint]*RunnerJob
	isDead          bool
	timeoutDuration time.Duration
	notifier        RunnerNotifier
}

var runners map[uint]*Runner

func initRunners() {
	runners = make(map[uint]*Runner)
}

func NewRunner(pipeline *PipelineWorker, notifier RunnerNotifier) *Runner {
	runner := Runner{
		pipeline:        *pipeline,
		allJobs:         make(map[uint]*RunnerJob),
		plannedJobs:     make(map[uint]*RunnerJob),
		isDead:          false,
		timeoutDuration: time.Minute * time.Duration(25),
		notifier:        notifier,
	}

	for _, jobWorker := range pipeline.JobWorkers {
		runnerJob := RunnerJob{
			jobWorker:    &jobWorker,
			status:       statusOption(jobWorker.StatusId),
			dependencies: []*RunnerJob{},
			dependents:   []*RunnerJob{},
			cancel:       nil,
			plannedBy:    nil,
		}

		runner.allJobs[jobWorker.ID] = &runnerJob
	}

	for _, runnerJob := range runner.allJobs {
		visited := map[uint]bool{runnerJob.jobWorker.ID: true}
		runnerJob.dependencies = runner.collectDependencies(runnerJob, visited)
		for _, dep := range runnerJob.dependencies {
			dep.dependents = append(dep.dependents, runnerJob)
		}
	}

	go runner.Start()

	return &runner
}

func GetRunner(pipelineWorkerId uint) *Runner {
	runner, ok := runners[pipelineWorkerId]
	if !ok {
		notifier := &RunnerNotifierHub{pipelineId: pipelineWorkerId}
		var pipelineWorker PipelineWorker
		db.Preload("JobWorkers").Preload("User").Preload("JobWorkers.StageWorkers", func(db *gorm.DB) *gorm.DB {
			return db.Order("id ASC")
		}).Preload("JobWorkers.StageWorkers.Stage").Preload("JobWorkers.Job").Preload("JobWorkers.Job.Dependencies").Preload("JobWorkers.Job.Dependencies.JobWorker").First(&pipelineWorker, pipelineWorkerId)
		runner = NewRunner(&pipelineWorker, notifier)
		runners[pipelineWorkerId] = runner
	}
	return runner
}

func (r *Runner) Start() {
	idleTimer := time.NewTimer(r.timeoutDuration)
	defer idleTimer.Stop()

	for {
		r.mu.Lock()
		jobsToRun := r.getReadyJobs()
		r.mu.Unlock()

		if len(jobsToRun) > 0 {

			idleTimer.Reset(r.timeoutDuration)

			// Execute ready jobs in parallel
			for _, job := range jobsToRun {
				job.status = statusRunning
				go r.runJob(job)
			}
		} else {
			select {
			case <-idleTimer.C:
				log.Println("Runner has been idle for too long. Terminating.")
				delete(runners, r.pipeline.PipelineId)
				return
			default:
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

// Run pipeline from the very start, only manually
func (r *Runner) RunPipeline(userId uint) {

	log.Println("pipeline started")

	startedAt := time.Now()

	// Make sure pipeline database is clean
	db.Model(&JobWorker{}).Where("pipeline_worker_id = ?", r.pipeline.ID).Select("StatusId", "FinishedAt", "StartedAt", "UserId").Updates(JobWorker{StatusId: uint(statusPlanned), FinishedAt: nil, StartedAt: nil, UserId: &userId})
	db.Model(&PipelineWorker{}).Where("id = ?", r.pipeline.ID).Select("UserId", "StartedAt", "FinishedAt", "StatusId").Updates(PipelineWorker{
		UserId:     &userId,
		StartedAt:  &startedAt,
		FinishedAt: nil,
		StatusId:   uint(statusRunning),
	})
	stagesIds := r.getPipelineJobStagesIds()
	r.deleteAllPipelineLogs(stagesIds)
	db.Model(&StageWorker{}).Where("stage_id in ?", stagesIds).Select("StatusId", "FinishedAt", "StartedAt").Updates(StageWorker{StatusId: uint(statusPending), FinishedAt: nil, StartedAt: nil})

	var user User
	db.First(&user, userId)

	// Make all jobs planned
	r.mu.Lock()
	r.pipeline.StartedAt = &startedAt
	r.pipeline.User = &user
	r.pipeline.StatusId = uint(statusRunning)
	for _, job := range r.allJobs {
		r.plannedJobs[job.jobWorker.ID] = job
		job.plannedBy = &user
		job.status = statusPlanned
		if job.cancel != nil {
			job.cancel()
		}
	}
	// Determine which will begin as running and which will stay planned
	planned, running := r.getPipelineStartKeys()
	r.mu.Unlock()

	// Notify on changes
	statuses := statusBatch{
		Planned: planned,
		Running: running,
	}
	stats := &PipelineRunStatistics{
		AverageDuration: nil,
		TotalDuration:   nil,
		TotalErrors:     0,
		StartedAt:       &startedAt,
		StartedBy:       &user.Name,
		FinishedAt:      nil,
	}
	meta := &JobMetadata{
		Fields:     []string{"startedAt", "finishedAt", "startedBy"},
		StartedBy:  &user.Name,
		FinishedAt: nil,
		StartedAt:  nil,
	}
	go r.notifier.NotifyJobStatusBatch(&statuses, meta, statusOption(r.pipeline.StatusId), stats)
}

// Manually stop pipeline
func (r *Runner) StopPipeline() {
	if r.pipeline.StatusId != uint(statusRunning) {
		return
	}

	log.Println("pipeline stopped")

	finishedAt := time.Now()

	// Notify first, because there will be another notification about pipeline finish coming later
	statuses := statusBatch{
		Cancelled: r.getPlannedJobsKeys(),
	}
	meta := &JobMetadata{
		FinishedAt: &finishedAt,
		Fields:     []string{"finishedAt"},
	}
	go r.notifier.NotifyJobStatusBatch(&statuses, meta, statusOption(r.pipeline.StatusId), nil)

	// Anything that was running becomes cancelled
	// local changes

	r.mu.Lock()
	r.clearPlannedJobs()
	r.mu.Unlock()
	for _, job := range r.allJobs {
		if job.cancel != nil {
			// cancel will handle all the logic
			// Since all jobs would end up being removed, this will also trigger FinishPipeline
			r.mu.Lock()
			job.status = statusCancelled
			r.mu.Unlock()
			job.cancel()
		} else if job.status == statusPlanned {
			// If job is just planned, then make it cancelled
			r.mu.Lock()
			job.status = statusCancelled
			r.mu.Unlock()
		}
	}

	// Database update
	db.Model(&JobWorker{}).Where("status_id = ? OR status_id = ?", statusPlanned, statusRunning).Select("StatusId", "FinishedAt").Updates(JobWorker{StatusId: uint(statusCancelled), FinishedAt: &finishedAt})

	log.Println("pipeline stopped till ze end")
}

// Clean up for pipeline, once it's considered finished
func (r *Runner) FinishPipeline() {
	log.Println("pipeline finished")
	finishedAt := time.Now()

	r.mu.Lock()
	r.pipeline.FinishedAt = &finishedAt
	r.mu.Unlock()
	r.pipeline.StatusId = uint(r.calculatePipelineStatus())

	// Update database
	db.Model(&PipelineWorker{}).Where("id = ?", r.pipeline.ID).Updates(PipelineWorker{
		FinishedAt: &finishedAt,
		StatusId:   r.pipeline.StatusId,
	})

	// Notify about pipeline new stats
	// No changes for jobs should have happened at this point
	var stats PipelineStatisticsShort
	fetchPipelineShortStatistics(&stats, r.pipeline.PipelineId)

	go r.notifier.NotifyJobStatusBatch(&statusBatch{}, &JobMetadata{}, statusOption(r.pipeline.StatusId), &PipelineRunStatistics{
		AverageDuration: &stats.AverageDuration,
		TotalErrors:     stats.TotalErrors,
		TotalDuration:   &stats.TotalDuration,
		FinishedAt:      &finishedAt,
		StartedAt:       r.pipeline.StartedAt,
		StartedBy:       &r.pipeline.User.Name,
	})
}

func (r *Runner) ContinuePipeline(userId uint) {
	log.Println("pipeline continued")
	if r.pipeline.StatusId == uint(statusPending) || r.pipeline.StatusId == uint(statusRunning) {
		return
	}

	startedAt := time.Now()

	var user User
	db.First(&user, userId)

	db.Model(&PipelineWorker{}).Where("id = ?", r.pipeline.ID).Updates(PipelineWorker{
		FinishedAt: nil,
		StatusId:   uint(statusRunning),
	})

	db.Model(&JobWorker{}).Where("pipeline_worker_id = ? AND (status_id = ? OR status_id = ? OR status_id = ?)", r.pipeline.ID, statusCancelled, statusPending, statusFailed).Updates(PipelineWorker{
		FinishedAt: nil,
		StartedAt:  &startedAt,
		StatusId:   uint(statusPlanned),
	})

	// update local storage
	// and launch all possible jobs
	r.mu.Lock()
	r.pipeline.FinishedAt = nil
	r.pipeline.StatusId = uint(statusRunning)

	for _, job := range r.allJobs {
		if job.status == statusPending || job.status == statusFailed || job.status == statusCancelled {
			r.plannedJobs[job.jobWorker.ID] = job
			job.plannedBy = &user
			job.status = statusPlanned
		}
	}

	planned, running := r.getPipelineStartKeys()
	r.mu.Unlock()

	meta := &JobMetadata{
		FinishedAt: nil,
		StartedAt:  &startedAt,
		StartedBy:  &user.Name,
		Fields:     []string{"finishedAt", "StartedAt", "StartedBy"},
	}
	go r.notifier.NotifyJobStatusBatch(&statusBatch{
		Planned: planned,
		Running: running,
	}, meta, statusRunning, &PipelineRunStatistics{
		AverageDuration: nil,
		TotalErrors:     r.countTotalErrors(),
		TotalDuration:   nil,
		StartedAt:       r.pipeline.StartedAt,
		StartedBy:       &r.pipeline.User.Name,
		FinishedAt:      nil,
	})
}

func (r *Runner) getReadyJobs() []*RunnerJob {
	readyJobs := []*RunnerJob{}

	for _, job := range r.plannedJobs {
		allDepsDone := true

		if job.status == statusRunning {
			continue
		}

		// Check if all dependencies of this job are complete
		for _, depJob := range job.dependencies {
			if depJob.status != statusCompleted {
				allDepsDone = false
				break
			}
		}

		// If all dependencies are complete, add job to readyJobs
		if allDepsDone {
			readyJobs = append(readyJobs, job)
		}
	}

	return readyJobs
}

func (r *Runner) calculatePipelineStatus() statusOption {
	r.mu.Lock()
	defer r.mu.Unlock()

	allCompleted := true
	hasFailed := false
	hasRunning := false
	hasCancelled := false

	for _, job := range r.allJobs {
		switch job.status {
		case statusCancelled:
			hasCancelled = true
			allCompleted = false
		case statusRunning:
			hasRunning = true
			allCompleted = false
		case statusFailed:
			hasFailed = true
		case statusPending:
			allCompleted = false
		case statusPlanned:
			hasRunning = true
			allCompleted = false
		case statusCompleted:
		}
	}

	if allCompleted {
		return statusCompleted
	}

	if hasRunning {
		return statusRunning
	}

	if hasFailed {
		return statusFailed
	}

	if hasCancelled {
		return statusCancelled
	}

	return statusPending
}

func (r *Runner) runJob(job *RunnerJob) {
	ctx, cancel := context.WithCancel(context.Background())
	r.mu.Lock()
	job.cancel = cancel
	job.status = statusRunning
	r.mu.Unlock()

	defer func() {
		// Special case, probably was restarted
		if job.status == statusPlanned {
			return
		}

		r.mu.Lock()
		job.plannedBy = nil
		job.cancel = nil
		leftJobs := len(r.plannedJobs)
		r.mu.Unlock()

		// If none is planned anymore, we stop pipeline
		if leftJobs == 0 && r.pipeline.StatusId == uint(statusRunning) {
			r.FinishPipeline()
		}
	}()

	// Make sure job is clean
	db.Model(&StageWorker{}).Where("job_worker_id = ?", job.jobWorker.ID).Select("StatusId", "StartedAt", "FinishedAt").Updates(StageWorker{StatusId: uint(statusPending), StartedAt: nil, FinishedAt: nil})
	r.deleteAllJobLogs(job)

	log.Printf("Running job \"%d\"\n", job.jobWorker.ID)

	// Update local storage and database
	startedAt := time.Now()

	// Update database
	userId := uint(0)
	var userName *string
	if job.plannedBy == nil {
		log.Printf("Dangerous error on job %d - plannedBy did not exist", job.jobWorker.ID)
	} else {
		userId = job.plannedBy.ID
		userName = &job.plannedBy.Name
	}
	db.Model(&JobWorker{}).Where("id = ?", job.jobWorker.ID).Select("StatusId", "FinishedAt", "StartedAt", "UserId").Updates(JobWorker{StatusId: uint(statusRunning), StartedAt: &startedAt, FinishedAt: nil, UserId: &userId})

	// Update local storage
	job.jobWorker.FinishedAt = nil
	job.jobWorker.StartedAt = &startedAt

	// Notify that job has started
	go r.notifier.NotifyJobStatus(job.jobWorker.ID, statusRunning, statusRunning, userName, &startedAt, nil)

	// JOB STARTED
	i := 1.0
	stagesLen := float64(len(job.jobWorker.StageWorkers))
	statusToAssign := statusCompleted
loop:
	for _, stage := range job.jobWorker.StageWorkers {
	checkAgain:
		select {
		// JOB CANCELLED
		case <-ctx.Done():
			finishedAt := time.Now()

			// Special case, probably was restarted
			if job.status == statusPlanned {
				return
			}

			// Update local storage
			r.mu.Lock()
			job.status = statusCancelled
			job.jobWorker.FinishedAt = &finishedAt
			delete(r.plannedJobs, job.jobWorker.ID)
			r.mu.Unlock()

			// Update database
			db.Model(&JobWorker{}).Where("id = ?", job.jobWorker.ID).Updates(JobWorker{StatusId: uint(statusCancelled), FinishedAt: &finishedAt})

			// Notify about cancellation
			go r.notifier.NotifyJobStatus(job.jobWorker.ID, statusCancelled, r.calculatePipelineStatus(), userName, job.jobWorker.StartedAt, &finishedAt)

			log.Printf("Job \"%d\" was cancelled\n", job.jobWorker.ID)
			return
		// JOB RUNNING
		default:
			// STAGE START
			startedAt := time.Now()

			stage.StartedAt = &startedAt

			// STAGE IN PROGRESS
			stageStatus := r.runStage(stage, ctx)
			if stageStatus == statusCancelled {
				goto checkAgain
			}

			// STAGE FINISHED

			if stageStatus == statusFailed {
				statusToAssign = statusFailed
				break loop
			}

			// Notify about job progress
			progress := uint((i / stagesLen) * 100)
			go r.notifier.NotifyJobProgress(job.jobWorker.ID, progress)

			i++
		}
	}

	// JOB FINISHED
	finishedAt := time.Now()

	// Update local storage
	r.mu.Lock()
	r.allJobs[job.jobWorker.ID].status = statusToAssign
	job.jobWorker.FinishedAt = &finishedAt
	delete(r.plannedJobs, job.jobWorker.ID)
	r.mu.Unlock()

	// Update database
	db.Model(&JobWorker{}).Where("id = ?", job.jobWorker.ID).Updates(JobWorker{StatusId: uint(statusToAssign), FinishedAt: &finishedAt})

	// JOB FAILED
	if statusToAssign == statusFailed {
		jobIds := make([]uint, len(job.dependents))
		r.mu.Lock()
		// Cancel all dependants and update local storage
		for i, dep := range job.dependents {
			delete(r.plannedJobs, dep.jobWorker.ID)
			dep.status = statusCancelled
			jobIds[i] = dep.jobWorker.ID
		}
		r.mu.Unlock()

		// Notify about cancelled jobs
		meta := &JobMetadata{
			FinishedAt: &finishedAt,
			Fields:     []string{"finishedAt"},
		}
		go r.notifier.NotifyJobStatusBatch(&statusBatch{
			Cancelled: jobIds,
		}, meta, statusRunning, nil)

		// Update database
		db.Model(&JobWorker{}).Where("id in ?", jobIds).Updates(JobWorker{
			StatusId:   uint(statusCancelled),
			FinishedAt: &finishedAt,
		})
	}

	// Notify about job completion
	go r.notifier.NotifyJobStatus(job.jobWorker.ID, statusToAssign, r.calculatePipelineStatus(), userName, nil, &finishedAt)

	log.Printf("Job \"%d\" completed\n", job.jobWorker.ID)
}

func (r *Runner) addRunnableJobDeepDependencies(runnerJob *RunnerJob, user *User) {
	r.plannedJobs[runnerJob.jobWorker.ID] = runnerJob
	runnerJob.plannedBy = user
	runnerJob.status = statusPlanned
	for _, dep := range runnerJob.dependencies {
		if dep.status == statusFailed || dep.status == statusPending || dep.status == statusCancelled {
			r.addRunnableJobDeepDependencies(dep, user)
		}
	}
}

// Start job manually (also restart)
func (r *Runner) startJob(jobId uint, userId uint) {

	r.mu.Lock()
	job, ok := r.allJobs[jobId]
	r.mu.Unlock()

	startedAt := time.Now()

	if ok {
		if job.status == statusRunning || job.status == statusPlanned {
			return
		}

		// Since we start job, pipeline is definetly RUNNING now
		db.Model(&PipelineWorker{}).Where("id = ?", r.pipeline.ID).Updates(PipelineWorker{
			StatusId:   uint(statusRunning),
			FinishedAt: nil,
		})

		r.mu.Lock()
		r.pipeline.FinishedAt = nil
		r.pipeline.StatusId = uint(statusRunning)
		r.mu.Unlock()

		var user User
		db.First(&user, userId)

		// IF for some reason pipeline was not run by anyone before
		// Person that started individual job will be considered pipeline launcher
		if r.pipeline.User == nil {
			db.Model(&PipelineWorker{}).Where("id = ?", r.pipeline.ID).Update("user_id", user.ID)
			r.pipeline.User = &user
		}

		// All dependant jobs would end up being cancelled

		// Update local storage
		jobWorkerIds := make([]uint, len(job.dependents))
		for i, dep := range job.dependents {
			dep.status = statusCancelled
			// Remove them from being planned, if they are
			delete(r.plannedJobs, dep.jobWorker.ID)
			// Probably impossible ?
			if dep.cancel != nil {
				dep.cancel()
			}
			jobWorkerIds[i] = dep.jobWorker.ID
		}

		r.mu.Lock()
		// Add current job and all dependecies to plannedJobs
		r.addRunnableJobDeepDependencies(job, &user)
		r.mu.Unlock()

		planned, running := r.getPipelineStartKeys()

		// Notify about cancelled, planned and running jobs
		// AND the fact that pipeline was started again
		meta := &JobMetadata{
			FinishedAt: &startedAt,
			StartedAt:  &startedAt,
			StartedBy:  &user.Name,
			Fields:     []string{"finishedAt", "startedAt", "startedBy"},
		}
		println(jobWorkerIds)
		go r.notifier.NotifyJobStatusBatch(&statusBatch{
			Cancelled: jobWorkerIds,
			Planned:   planned,
			Running:   running,
		}, meta, statusOption(r.pipeline.StatusId), &PipelineRunStatistics{
			AverageDuration: nil,
			TotalErrors:     r.countTotalErrors(),
			TotalDuration:   nil,
			FinishedAt:      nil,
			StartedAt:       r.pipeline.StartedAt,
			StartedBy:       &r.pipeline.User.Name,
		})

	}
}

// Stop job manually
func (r *Runner) stopJob(jobId uint) {
	r.mu.Lock()
	job, ok := r.allJobs[jobId]
	r.mu.Unlock()

	if ok {
		if !(job.status == statusRunning || job.status == statusPlanned) {
			return
		}

		finishedAt := time.Now()

		// Cancel job execution
		if job.status == statusRunning {
			if job.cancel != nil {
				// Cancel would also trigger check for pipeline status from running job
				job.cancel()
			} else {
				log.Println("Dangerous error - job.cancel did not exist on job")
			}
		}

		// Update local storage
		r.mu.Lock()
		job.status = statusCancelled
		delete(r.plannedJobs, job.jobWorker.ID)
		// Since we cancelled this job, anything that depended on it, would be cancelled too
		depIds := make([]uint, 0)
		for _, dep := range job.dependents {
			delete(r.plannedJobs, dep.jobWorker.ID)
			if dep.status == statusPlanned {
				dep.status = statusCancelled
				depIds = append(depIds, dep.jobWorker.ID)
				dep.jobWorker.FinishedAt = &finishedAt
				dep.jobWorker.StartedAt = nil
			}
		}
		job.jobWorker.FinishedAt = &finishedAt
		job.jobWorker.StartedAt = nil
		r.mu.Unlock()
		r.pipeline.StatusId = uint(r.calculatePipelineStatus())

		// Update database
		depIds = append(depIds, jobId)
		db.Model(&JobWorker{}).Where("id in ?", depIds).Select("StartedAt", "FinishedAt", "StatusId").Updates(JobWorker{
			FinishedAt: &finishedAt,
			StartedAt:  nil,
			StatusId:   uint(statusCancelled),
		})

		meta := &JobMetadata{
			FinishedAt: &finishedAt,
			StartedAt:  nil,
			Fields:     []string{"finishedAt", "startedAt"},
		}
		go r.notifier.NotifyJobStatusBatch(&statusBatch{
			Cancelled: depIds,
		}, meta, statusOption(r.pipeline.StatusId), nil)

		// if we have stopped the last job available, have to trigger finishing
		if r.pipeline.StatusId != uint(statusRunning) {
			r.FinishPipeline()
		}

		//go r.notifier.NotifyJobStatus(jobId, statusCancelled, statusOption(r.pipeline.StatusId), userName, job.jobWorker.StartedAt, &finishedAt)
	}
}

func (r *Runner) runStage(stageWorker StageWorker, ctx context.Context) statusOption {
	log.Printf("Running stage \"%d\"...\n", stageWorker.ID)

	startedAt := time.Now()
	db.Model(&StageWorker{}).Where("id = ?", stageWorker.ID).Updates(StageWorker{StatusId: uint(statusRunning), StartedAt: &startedAt})

	// Notify stage has started
	go r.notifier.NotifyStageStatus(stageWorker.JobWorkerId, stageWorker.ID, &stageStatusExtra{
		StartedAt:  &startedAt,
		FinishedAt: nil,
		Status:     statusRunning.String(),
	})

	willFail := stageWorker.Stage.WillFail
	duration := stageWorker.Stage.RunsFor
	if duration == 0 {
		duration = 1 + rand.UintN(4)
	}

	times := int(duration) * 5
	failAt := times / 2
	statusToAssign := statusCompleted
loop:
	for i := 0; i < times; i++ {
	goBack:
		select {
		case <-ctx.Done():
			finishedAt := time.Now()
			db.Model(&StageWorker{}).Where("id = ?", stageWorker.ID).Updates(StageWorker{StatusId: uint(statusCancelled), FinishedAt: &finishedAt})

			log.Printf("Stage \"%d\" cancelled\n", stageWorker.ID)
			return statusCancelled
		default:

			if i%2 == 0 {
				if willFail && i >= failAt {
					statusToAssign = statusFailed
				}

				justCreated := createLog(stageWorker.ID, willFail)

				logs := make([]StageLogData, len(justCreated))
				for i, log := range justCreated {
					logs[i] = StageLogData{
						Line:      log.Line,
						Text:      log.Text,
						CreatedAt: log.CreatedAt,
						Type:      logTypeOptions(log.LogTypeId).String(),
					}
				}

				isForced := i == 0 || (times-i) <= 1
				r.notifier.NotifyStageLogs(stageWorker.JobWorkerId, stageWorker.ID, logs, isForced)
				if statusToAssign == statusFailed {
					break loop
				}
			}

			time.Sleep(200 * time.Millisecond)
			if ctx.Err() != nil {
				goto goBack
			}
		}
	}

	finishedAt := time.Now()

	// Notify about stage completion
	go r.notifier.NotifyStageStatus(stageWorker.JobWorkerId, stageWorker.ID, &stageStatusExtra{
		StartedAt:  &startedAt,
		FinishedAt: &finishedAt,
		Status:     statusToAssign.String(),
	})

	db.Model(&StageWorker{}).Where("id = ?", stageWorker.ID).Updates(StageWorker{StatusId: uint(statusToAssign), FinishedAt: &finishedAt})

	log.Printf("Stage \"%d\" completed\n", stageWorker.ID)
	return statusToAssign
}
