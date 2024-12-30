package main

func (r *Runner) getPlannedJobsKeys() []uint {
	keys := []uint{}
	for k := range r.plannedJobs {
		keys = append(keys, uint(k))
	}
	return keys
}

func (r *Runner) getPipelineStartKeys() ([]uint, []uint) {
	jobsToBeRunned := r.getReadyJobs()

	plannedIds := []uint{}
	runningIds := make([]uint, len(jobsToBeRunned))
	for i, job := range jobsToBeRunned {
		runningIds[i] = job.jobWorker.ID
	}

	for k := range r.plannedJobs {
		inRunning := false
		for _, job := range jobsToBeRunned {
			if uint(k) == job.jobWorker.ID {
				inRunning = true
				break
			}
		}
		if !inRunning {
			plannedIds = append(plannedIds, uint(k))
		}
	}

	return plannedIds, runningIds
}

func (r *Runner) clearPlannedJobs() {
	for k, job := range r.plannedJobs {
		job.plannedBy = nil
		delete(r.plannedJobs, k)
	}
}

func (r *Runner) countTotalErrors() uint {
	c := 0
	for _, job := range r.allJobs {
		if job.status == statusFailed {
			c++
		}
	}
	return uint(c)
}

func (r *Runner) collectDependencies(runnerJob *RunnerJob, visited map[uint]bool) []*RunnerJob {
	var deps []*RunnerJob
	for _, dep := range runnerJob.jobWorker.Job.Dependencies {
		if visited[dep.JobWorker.ID] {
			continue
		}
		visited[dep.ID] = true

		runnerDep, ok := r.allJobs[dep.ID]
		if ok {
			deps = append(deps, runnerDep)

			deps = append(deps, r.collectDependencies(runnerDep, visited)...)
		}
	}
	return deps
}

func (r *Runner) getJobStagesIds(job *RunnerJob) []uint {
	slice := make([]uint, len(job.jobWorker.StageWorkers))
	for _, stage := range job.jobWorker.StageWorkers {
		slice = append(slice, stage.ID)
	}
	return slice
}

func (r *Runner) getPipelineJobStagesIds() []uint {
	slice := make([]uint, 0)

	for _, job := range r.allJobs {
		for _, stage := range job.jobWorker.StageWorkers {
			slice = append(slice, stage.ID)
		}
	}

	return slice
}

func (r *Runner) deleteAllJobLogs(job *RunnerJob) error {
	query := `
		DELETE FROM stage_logs WHERE stage_worker_id IN ?;
	`

	if err := db.Exec(query, r.getJobStagesIds(job)).Error; err != nil {
		return err
	}

	return nil
}

func (r *Runner) deleteAllPipelineLogs(ids []uint) error {
	query := `
		DELETE FROM stage_logs WHERE stage_worker_id IN ?;
	`

	if err := db.Exec(query, ids).Error; err != nil {
		return err
	}

	return nil
}
