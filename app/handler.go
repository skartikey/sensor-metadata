package app

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// Handler represents the HTTP handlers for the API endpoints.
type Handler struct {
	repo      Repository
	validator *validator.Validate
}

// NewHandler creates a new instance of the Handler.
func NewHandler(repo Repository) *Handler {
	return &Handler{
		repo:      repo,
		validator: validator.New(),
	}
}

// ErrorResponse represents the structure of error responses.
type ErrorResponse struct {
	Message string `json:"message"`
}

// CreateSensorMetadata handles the HTTP POST request to create sensor metadata.
func (h *Handler) CreateSensorMetadata(w http.ResponseWriter, r *http.Request) {
	var sensorMetadata SensorMetadata
	err := json.NewDecoder(r.Body).Decode(&sensorMetadata)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate the input
	if err := h.validator.Struct(sensorMetadata); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Save the sensor metadata
	err = h.repo.CreateSensorMetadata(&sensorMetadata)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to create sensor metadata")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetSensorMetadata handles the HTTP GET request to retrieve sensor metadata by name.
func (h *Handler) GetSensorMetadata(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		sendErrorResponse(w, http.StatusBadRequest, "Missing 'name' parameter")
		return
	}

	sensorMetadata, err := h.repo.GetSensorMetadataByName(name)
	if err != nil {
		sendErrorResponse(w, http.StatusNotFound, "Sensor metadata not found")
		return
	}

	jsonResponse(w, http.StatusOK, sensorMetadata)
}

// UpdateSensorMetadata handles the HTTP PUT request to update sensor metadata.
func (h *Handler) UpdateSensorMetadata(w http.ResponseWriter, r *http.Request) {
	var sensorMetadata SensorMetadata
	err := json.NewDecoder(r.Body).Decode(&sensorMetadata)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate the input
	if err := h.validator.Struct(sensorMetadata); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Update the sensor metadata
	err = h.repo.UpdateSensorMetadata(&sensorMetadata)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to update sensor metadata")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetNearestSensor handles the HTTP GET request to find the sensor nearest to a given location.
func (h *Handler) GetNearestSensorMetadata(w http.ResponseWriter, r *http.Request) {
	latitude := r.URL.Query().Get("latitude")
	longitude := r.URL.Query().Get("longitude")
	if latitude == "" || longitude == "" {
		sendErrorResponse(w, http.StatusBadRequest, "Missing 'latitude' or 'longitude' parameter")
		return
	}

	// Query the nearest sensor
	sensorMetadata, err := h.repo.GetNearestSensorMetadata(latitude, longitude)
	if err != nil {
		sendErrorResponse(w, http.StatusNotFound, "No nearest sensor found")
		return
	}

	jsonResponse(w, http.StatusOK, sensorMetadata)
}

// Helper function to send JSON response with appropriate status code.
func jsonResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		// Handle encoding error
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Helper function to send error response with appropriate status code.
func sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	errorResponse := ErrorResponse{
		Message: message,
	}
	jsonResponse(w, statusCode, errorResponse)
}
