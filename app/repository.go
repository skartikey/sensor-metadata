package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

// Repository represents the interface for interacting with the database.
type Repository interface {
	CreateSensorMetadata(sensorMetadata *SensorMetadata) error
	GetSensorMetadataByName(name string) (*SensorMetadata, error)
	UpdateSensorMetadata(sensorMetadata *SensorMetadata) error
	GetNearestSensorMetadata(latitude, longitude string) (*SensorMetadata, error)
	GetNearestSensorByCityMetadata(city string) (string, error)
}

// PostgresRepository represents the PostgreSQL repository implementation.
type PostgresRepository struct {
	Db *sql.DB
}

// NewPostgresRepository creates a new instance of the PostgreSQL repository.
func NewPostgresRepository() (*PostgresRepository, error) {
	// Read the environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	// Construct the connection string
	connStr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", host, port, dbName, user, password)

	// Connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Check the database connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{Db: db}, nil
}

// CreateSensorMetadata creates a new sensor metadata entry in the database.
func (r *PostgresRepository) CreateSensorMetadata(sensorMetadata *SensorMetadata) error {
	// Prepare the SQL statement
	stmt, err := r.Db.Prepare("INSERT INTO sensor_metadata (name, location_latitude, location_longitude, tags) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(sensorMetadata.Name, sensorMetadata.Location.Latitude, sensorMetadata.Location.Longitude, pq.Array(sensorMetadata.Tags))
	if err != nil {
		return err
	}

	return nil
}

// GetSensorMetadataByName retrieves sensor metadata from the database by name.
func (r *PostgresRepository) GetSensorMetadataByName(name string) (*SensorMetadata, error) {
	// Prepare the SQL statement
	stmt, err := r.Db.Prepare("SELECT id, name, location_latitude, location_longitude, tags FROM sensor_metadata WHERE name = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the SQL statement
	row := stmt.QueryRow(name)

	// Initialize a SensorMetadata struct to store the result
	var sensorMetadata SensorMetadata

	// Scan the result into the SensorMetadata struct
	err = row.Scan(&sensorMetadata.ID, &sensorMetadata.Name, &sensorMetadata.Location.Latitude, &sensorMetadata.Location.Longitude, pq.Array(&sensorMetadata.Tags))
	if err != nil {
		if err == sql.ErrNoRows {
			// Return nil and a custom error if no rows are found
			return nil, fmt.Errorf("sensor metadata not found")
		}
		return nil, err
	}

	return &sensorMetadata, nil
}

// UpdateSensorMetadata updates an existing sensor metadata entry in the database.
func (r *PostgresRepository) UpdateSensorMetadata(sensorMetadata *SensorMetadata) error {
	// Prepare the SQL statement
	stmt, err := r.Db.Prepare("UPDATE sensor_metadata SET name = $1, location_latitude = $2, location_longitude = $3, tags = $4 WHERE id = $5")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(sensorMetadata.Name, sensorMetadata.Location.Latitude, sensorMetadata.Location.Longitude, pq.Array(sensorMetadata.Tags), sensorMetadata.ID)
	if err != nil {
		return err
	}

	return nil
}

// GetNearestSensorMetadata retrieves the nearest sensor metadata from the database based on location.
func (r *PostgresRepository) GetNearestSensorMetadata(latitude, longitude string) (*SensorMetadata, error) {
	// Prepare the SQL statement
	stmt, err := r.Db.Prepare("SELECT id, name, location_latitude, location_longitude, tags, earth_distance(ll_to_earth($1, $2), ll_to_earth(location_latitude, location_longitude)) AS distance FROM sensor_metadata ORDER BY distance LIMIT 1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the SQL statement
	row := stmt.QueryRow(latitude, longitude)

	// Initialize a SensorMetadata struct to store the result
	var sensorMetadata SensorMetadata

	// Scan the result into the SensorMetadata struct
	err = row.Scan(&sensorMetadata.ID, &sensorMetadata.Name, &sensorMetadata.Location.Latitude, &sensorMetadata.Location.Longitude, pq.Array(&sensorMetadata.Tags), &sensorMetadata.Distance)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return nil and a custom error if no rows are found
			return nil, fmt.Errorf("sensor metadata not found")
		}
		return nil, err
	}

	return &sensorMetadata, nil
}

func (r *PostgresRepository) GetNearestSensorByCityMetadata(city string) (string, error) {
	apiKey := os.Getenv("API_KEY")
	encodedCity := url.QueryEscape(city)
	url := fmt.Sprintf("https://api.mapbox.com/geocoding/v5/mapbox.places/%s.json?access_token=%s", encodedCity, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response MapboxResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	if len(response.Features) > 0 {
		coordinates := response.Features[0].Geometry.Coordinates
		latitude := coordinates[1]
		longitude := coordinates[0]
		return fmt.Sprintf("%f,%f", latitude, longitude), nil
	}

	return "", fmt.Errorf("no sensor found for the city")

}
