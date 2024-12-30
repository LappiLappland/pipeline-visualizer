package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var HostAddress string

func main() {

	_ = godotenv.Load()
	HostAddress = os.Getenv("HOST_ADDRESS")

	err = DBConnection()
	if err != nil {
		log.Fatal("Database connection error", err)
	}

	r := mux.NewRouter()

	initSessionsStore()
	initCreateLog()
	initRunners()
	initHubs()

	r.HandleFunc("/api/login", loginHandler)
	r.HandleFunc("/api/logout", logoutHandler)
	r.HandleFunc("/api/user", getUser)
	r.HandleFunc("/api/pipelines", getPipelineRuns)

	r.HandleFunc("/api/pipeline/run/{id}", Adapt(
		http.HandlerFunc(getPipelineRun),
		requirePipelineAuth(read),
	).ServeHTTP)

	r.HandleFunc("/api/pipeline/run/{id}/statistics", Adapt(
		http.HandlerFunc(getPipelineRunStats),
		requirePipelineAuth(read),
	).ServeHTTP)
	r.HandleFunc("/api/pipeline/run/{id}/statistics/excel", Adapt(
		http.HandlerFunc(getPipelineRunStatsExcel),
		requirePipelineAuth(read),
	).ServeHTTP)
	r.HandleFunc("/api/pipeline/run/{id}/job/{job}", Adapt(
		http.HandlerFunc(getJobRun),
		requirePipelineAuth(read),
	).ServeHTTP)
	r.HandleFunc("/api/pipeline/run/{id}/job/{job}/stage/{stage}", Adapt(
		http.HandlerFunc(getStageRun),
		requirePipelineAuth(read),
	).ServeHTTP)

	r.HandleFunc("/api/pipeline/run/{id}/ws", Adapt(
		http.HandlerFunc(serveWs),
		requirePipelineAuth(read),
	).ServeHTTP)

	log.Println("Server is up and running")

	options := handlers.AllowedHeaders([]string{"Authorization"})
	err := http.ListenAndServe(":8000", handlers.CORS(options)(r))
	log.Print(err)
}
