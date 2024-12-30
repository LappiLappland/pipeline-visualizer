package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type PipelinesData struct {
	Name        string `json:"name"`
	Id          uint   `json:"id"`
	Description string `json:"description"`
}

func getPipelineRuns(writer http.ResponseWriter, request *http.Request) {

	var pipelines []Pipeline
	db.Preload("PipelineWorkers").Find(&pipelines)

	data := make([]PipelinesData, 0)

	for _, pipeline := range pipelines {
		for _, worker := range pipeline.PipelineWorkers {
			data = append(data, PipelinesData{Name: pipeline.Name, Id: worker.ID, Description: pipeline.Description})
		}
	}

	writer.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(writer).Encode(data)
	if err != nil {
		log.Println(err)
	}
}
