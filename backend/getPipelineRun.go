package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/imacks/bitflags-go"
	"gorm.io/gorm"
)

type PipelineData struct {
	ID           uint   `json:"id"`
	JobId        uint   `json:"jobId"`
	Title        string `json:"title"`
	Status       string `json:"status"`
	Dependencies []uint `json:"dependencies"`
}

type PipelineRunStatistics struct {
	AverageDuration *uint      `json:"averageDuration,omitempty"`
	TotalDuration   *uint      `json:"totalDuration,omitempty"`
	TotalErrors     uint       `json:"totalErrors"`
	StartedBy       *string    `json:"startedBy,omitempty"`
	StartedAt       *time.Time `json:"startedAt,omitempty"`
	FinishedAt      *time.Time `json:"finishedAt,omitempty"`
}

type PipelineRunJobData struct {
	PipelineData
	Progress   uint       `json:"progress"`
	StartedBy  *string    `json:"startedBy"`
	StartedAt  *time.Time `json:"startedAt"`
	FinishedAt *time.Time `json:"finishedAt"`
}

type PipelineRunData struct {
	Title          string                `json:"title"`
	Jobs           []PipelineRunJobData  `json:"jobs"`
	Statistics     PipelineRunStatistics `json:"statistics"`
	PipelineStatus string                `json:"pipelineStatus,omitempty"`
	Permissions    string                `json:"permissions"`
}

func getPipelineRun(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	idAsUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Println(err)
		return
	}

	session, _ := store.Get(r, "session-name")
	userId, _ := session.Values["id"].(uint)

	var pipelineWorker PipelineWorker
	db.Preload("User").Preload("Pipeline").Preload("Pipeline.Users", func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", userId)
	}).Preload("Pipeline.Users.Role").Preload("Pipeline.Commit").Preload("JobWorkers").Preload("JobWorkers.User").Preload("JobWorkers.Status").Preload("JobWorkers.Job").Preload("JobWorkers.Job.Dependencies").Preload("JobWorkers.StageWorkers").First(&pipelineWorker, id)

	permissions := []byte("0000")
	users := pipelineWorker.Pipeline.Users

	if len(users) != 0 {
		userPermissions := users[0].Role.Permissions
		if bitflags.Has(userPermissions, admin) {
			permissions[0] = byte('1')
		}
		if bitflags.Has(userPermissions, execute) {
			permissions[1] = byte('1')
		}
		if bitflags.Has(userPermissions, write) {
			permissions[2] = byte('1')
		}
		if bitflags.Has(userPermissions, read) {
			permissions[3] = byte('1')
		}
	}

	var stats PipelineStatisticsShort
	fetchPipelineShortStatistics(&stats, uint(idAsUint))

	var jobsToSend []PipelineRunJobData
	for _, jobWorker := range pipelineWorker.JobWorkers {

		var data PipelineRunJobData
		data.ID = jobWorker.ID
		data.JobId = jobWorker.Job.ID
		if jobWorker.Job.Name == nil {
			data.Title = jobWorker.Job.NameId
		} else {
			data.Title = *jobWorker.Job.Name
		}
		if jobWorker.User != nil {
			data.StartedBy = &jobWorker.User.Name
		}
		data.Status = jobWorker.Status.Name
		data.FinishedAt = jobWorker.FinishedAt
		data.StartedAt = jobWorker.StartedAt
		data.Dependencies = []uint{}
		for _, dep := range jobWorker.Job.Dependencies {
			data.Dependencies = append(data.Dependencies, dep.ID)
		}

		completedStages := 0.0
		for _, stage := range jobWorker.StageWorkers {
			if stage.StatusId == uint(statusCompleted) {
				completedStages++
			}
		}

		data.Progress = uint((completedStages / float64(len(jobWorker.StageWorkers))) * 100)

		jobsToSend = append(jobsToSend, data)
	}
	var pipelineRunData PipelineRunData
	pipelineRunData.Title = pipelineWorker.Pipeline.Name
	pipelineRunData.Jobs = jobsToSend
	pipelineRunData.Permissions = string(permissions)
	pipelineRunData.PipelineStatus = statusOption(pipelineWorker.StatusId).String()
	var startedBy string
	if pipelineWorker.User != nil {
		startedBy = pipelineWorker.User.Name
	}
	pipelineRunData.Statistics = PipelineRunStatistics{
		AverageDuration: &stats.AverageDuration,
		TotalErrors:     stats.TotalErrors,
		TotalDuration:   &stats.TotalDuration,
		StartedBy:       &startedBy,
		StartedAt:       pipelineWorker.StartedAt,
		FinishedAt:      pipelineWorker.FinishedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(pipelineRunData)
	if err != nil {
		log.Println(err)
	}
}
