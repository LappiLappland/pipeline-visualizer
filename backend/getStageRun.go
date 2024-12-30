package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type StageLogData struct {
	Line      uint      `json:"line"`
	Type      string    `json:"type"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
}

type StageLogRunData struct {
	Logs []StageLogData `json:"logs"`
}

func getStageRun(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	//id := vars["id"]
	//id := vars["job"]
	id := vars["stage"]

	var stageWorker StageWorker
	db.Preload("Logs", func(db *gorm.DB) *gorm.DB {
		return db.Order("line ASC")
	}).First(&stageWorker, id)

	var logsToSend []StageLogData
	for _, log := range stageWorker.Logs {

		var data StageLogData
		data.Line = log.Line
		data.Text = log.Text
		data.CreatedAt = log.CreatedAt
		data.Type = logTypeOptions(log.LogTypeId).String()

		logsToSend = append(logsToSend, data)
	}
	var stageLogRunData StageLogRunData
	stageLogRunData.Logs = logsToSend

	writer.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(writer).Encode(stageLogRunData)
	if err != nil {
		log.Println(err)
	}
}
