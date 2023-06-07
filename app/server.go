package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Server represents the HTTP server for the API.
type Server struct {
	handler *Handler
}

// NewServer creates a new instance of the HTTP server.
func NewServer(db *sql.DB) *Server {
	// Create a new instance of the PostgresRepository
	repo, err := NewPostgresRepository()
	if err != nil {
		log.Fatal("Error creating Postgres repository:", err)
	}
	handler := NewHandler(repo)
	return &Server{
		handler: handler,
	}
}

// Start starts the HTTP server and listens for incoming requests.
func (s *Server) Start(port string) error {
	router := mux.NewRouter()

	// Define API endpoints
	router.HandleFunc("/sensors", s.handler.CreateSensorMetadata).Methods(http.MethodPost)
	router.HandleFunc("/sensors", s.handler.GetSensorMetadata).Methods(http.MethodGet)
	router.HandleFunc("/sensors", s.handler.UpdateSensorMetadata).Methods(http.MethodPut)
	router.HandleFunc("/sensors/nearest", s.handler.GetNearestSensorMetadata).Methods(http.MethodGet)
	router.HandleFunc("/sensors", s.handler.GetNearestSensorByCityMetadata).Methods(http.MethodGet)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server started. Listening on %s", addr)
	return http.ListenAndServe(addr, router)
}
