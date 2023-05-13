package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/skartikey/sensor-metadata/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandlerCreateSensorMetadata(t *testing.T) {
	// Prepare mock data
	sensor := app.SensorMetadata{
		Name: "Sensor1",
		Location: app.Location{
			Latitude:  123.456,
			Longitude: 789.012,
		},
		Tags: []string{"tag1", "tag2"},
	}
	payload, _ := json.Marshal(sensor)

	// Create a request with the payload
	req, err := http.NewRequest(http.MethodPost, "/sensors", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Create a mock repository
	repo := &MockRepository{}
	repo.On("CreateSensorMetadata", mock.AnythingOfType("*app.SensorMetadata")).Return(nil)

	// Create a handler and serve the request
	handler := app.NewHandler(repo)
	http.HandlerFunc(handler.CreateSensorMetadata).ServeHTTP(rr, req)

	// Assert the status code and response
	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, "", rr.Body.String())
}

func TestGetSensorMetadata(t *testing.T) {
	// Prepare mock data
	expectedSensor := app.SensorMetadata{
		ID:   1,
		Name: "Sensor1",
		Location: app.Location{
			Latitude:  123.456,
			Longitude: 789.012,
		},
		Tags: []string{"tag1", "tag2"},
	}

	// Create a request with query parameters
	req, err := http.NewRequest(http.MethodGet, "/sensors?name=Sensor1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Create a mock repository
	repo := &MockRepository{}
	repo.On("GetSensorMetadataByName", "Sensor1").Return(&expectedSensor, nil)

	// Create a handler and serve the request
	handler := app.NewHandler(repo)
	http.HandlerFunc(handler.GetSensorMetadata).ServeHTTP(rr, req)

	// Assert the status code and response
	assert.Equal(t, http.StatusOK, rr.Code)

	var responseSensor app.SensorMetadata
	err = json.Unmarshal(rr.Body.Bytes(), &responseSensor)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedSensor, responseSensor)
}

// Define a mock repository for testing
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateSensorMetadata(sensorMetadata *app.SensorMetadata) error {
	args := m.Called(sensorMetadata)
	return args.Error(0)
}

func (m *MockRepository) GetSensorMetadataByName(name string) (*app.SensorMetadata, error) {
	args := m.Called(name)
	return args.Get(0).(*app.SensorMetadata), args.Error(1)
}

func (m *MockRepository) UpdateSensorMetadata(sensorMetadata *app.SensorMetadata) error {
	args := m.Called(sensorMetadata)
	return args.Error(0)
}

func (m *MockRepository) GetNearestSensorMetadata(latitude, longitude string) (*app.SensorMetadata, error) {
	args := m.Called(latitude, longitude)
	return args.Get(0).(*app.SensorMetadata), args.Error(1)
}
