package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type StageRunData struct {
	ID         uint       `json:"id"`
	Title      string     `json:"title"`
	Status     string     `json:"status"`
	StartedAt  *time.Time `json:"startedAt"`
	FinishedAt *time.Time `json:"finishedAt"`
}

type JobRunData struct {
	Stages []StageRunData `json:"stages"`
}

func getJobRun(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	//id := vars["id"]
	id := vars["job"]

	var jobWorker JobWorker
	db.Preload("StageWorkers", func(db *gorm.DB) *gorm.DB {
		return db.Order("id ASC")
	}).Preload("StageWorkers.Stage").Preload("StageWorkers.Status").First(&jobWorker, id)

	var stagesToSend []StageRunData
	for _, stageWorker := range jobWorker.StageWorkers {

		var data StageRunData
		data.ID = stageWorker.ID
		data.Title = stageWorker.Stage.Name
		data.Status = stageWorker.Status.Name
		data.StartedAt = stageWorker.StartedAt
		data.FinishedAt = stageWorker.FinishedAt

		stagesToSend = append(stagesToSend, data)
	}
	var jobRunData JobRunData
	jobRunData.Stages = stagesToSend

	writer.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(writer).Encode(jobRunData)
	if err != nil {
		log.Println(err)
	}
}
