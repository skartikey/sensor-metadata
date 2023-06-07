package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/gorilla/mux"

	"github.com/skartikey/sensor-metadata/app"
)

func main() {
	// Load the environment variables from .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Create a new instance of the PostgresRepository
	repo, err := app.NewPostgresRepository()
	if err != nil {
		log.Fatal("Error creating Postgres repository:", err)
	}

	// Create a new router
	router := mux.NewRouter()

	// Create a new handler and register the routes
	handler := app.NewHandler(repo)

	// Routes
	router.HandleFunc("/sensors", handler.CreateSensorMetadata).Methods("POST")
	router.HandleFunc("/sensors", handler.GetSensorMetadata).Methods("GET")
	router.HandleFunc("/sensors/{name}", handler.UpdateSensorMetadata).Methods("PUT")
	router.HandleFunc("/sensors/nearest", handler.GetNearestSensorMetadata).Methods("GET")
	router.HandleFunc("/sensors/{city}", handler.GetNearestSensorByCityMetadata).Methods("GET")

	// Start the HTTP server
	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
