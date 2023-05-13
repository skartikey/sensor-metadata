package app

import (
	"database/sql/driver"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/skartikey/sensor-metadata/app"
	"github.com/stretchr/testify/assert"
)

func TestPostgresRepository_CreateSensorMetadata(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	repo := &app.PostgresRepository{Db: mockDB}

	sensor := &app.SensorMetadata{
		Name: "Sensor1",
		Location: app.Location{
			Latitude:  123.456,
			Longitude: 789.012,
		},
		Tags: []string{"tag1", "tag2"},
	}

	expectedQuery := "INSERT INTO sensor_metadata (name, location_latitude, location_longitude, tags) VALUES ($1, $2, $3, $4)"

	mock.ExpectExec(expectedQuery).
		WithArgs(sqlmock.AnyArg(), 123.456, 789.012, AnyEmptyArray()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.CreateSensorMetadata(sensor)

	assert.NoError(t, mock.ExpectationsWereMet())
	assert.NoError(t, err)
}

func TestPostgresRepository_GetSensorMetadataByName(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	repo := &app.PostgresRepository{Db: mockDB}

	expectedSensor := &app.SensorMetadata{
		ID:   1,
		Name: "Sensor1",
		Location: app.Location{
			Latitude:  123.456,
			Longitude: 789.012,
		},
		Tags: []string{"tag1", "tag2"},
	}

	expectedQuery := "SELECT id, name, location_latitude, location_longitude, tags FROM sensor_metadata WHERE name = $1"
	expectedArgs := []driver.Value{"Sensor1"}

	mock.ExpectPrepare(expectedQuery).ExpectQuery().WithArgs(expectedArgs...).WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "location_latitude", "location_longitude", "tags"}).
			AddRow(expectedSensor.ID, expectedSensor.Name, expectedSensor.Location.Latitude, expectedSensor.Location.Longitude, pq.Array(expectedSensor.Tags)),
	)

	sensor, err := repo.GetSensorMetadataByName("Sensor1")

	assert.NoError(t, mock.ExpectationsWereMet())
	assert.NoError(t, err)
	assert.Equal(t, expectedSensor, sensor)
}

func TestPostgresRepository_UpdateSensorMetadata(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	repo := &app.PostgresRepository{Db: mockDB}

	sensor := &app.SensorMetadata{
		ID:   1,
		Name: "Sensor1",
		Location: app.Location{
			Latitude:  123.456,
			Longitude: 789.012,
		},
		Tags: []string{"tag1", "tag2"},
	}

	expectedQuery := "UPDATE sensor_metadata SET name = $1, location_latitude = $2, location_longitude = $3, tags = $4 WHERE id = $5"
	expectedArgs := []driver.Value{"Sensor1", 123.456, 789.012, pq.Array([]string{"tag1", "tag2"}), 1}

	mock.ExpectPrepare(expectedQuery).ExpectExec().WithArgs(expectedArgs...).WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.UpdateSensorMetadata(sensor)

	assert.NoError(t, mock.ExpectationsWereMet())
	assert.NoError(t, err)
}

func TestPostgresRepository_GetNearestSensorMetadata(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	repo := &app.PostgresRepository{Db: mockDB}

	expectedSensor := &app.SensorMetadata{
		ID:   1,
		Name: "Sensor1",
		Location: app.Location{
			Latitude:  123.456,
			Longitude: 789.012,
		},
		Tags: []string{"tag1", "tag2"},
	}

	expectedQuery := "SELECT id, name, location_latitude, location_longitude, tags, earth_distance(ll_to_earth($1, $2), ll_to_earth(location_latitude, location_longitude)) AS distance FROM sensor_metadata ORDER BY distance LIMIT 1"
	expectedArgs := []driver.Value{"123.456", "789.012"}

	mock.ExpectPrepare(expectedQuery).ExpectQuery().WithArgs(expectedArgs...).WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "location_latitude", "location_longitude", "tags"}).
			AddRow(expectedSensor.ID, expectedSensor.Name, expectedSensor.Location.Latitude, expectedSensor.Location.Longitude, pq.Array(expectedSensor.Tags)),
	)

	sensor, err := repo.GetNearestSensorMetadata("123.456", "789.012")

	assert.NoError(t, mock.ExpectationsWereMet())
	assert.NoError(t, err)
	assert.Equal(t, expectedSensor, sensor)
}

func TestNewPostgresRepository(t *testing.T) {
	// Set the required environment variables for the test
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpassword")
	repo, err := app.NewPostgresRepository()

	assert.NoError(t, err)
	assert.NotNil(t, repo.Db)
}

func AnyEmptyArray() interface{} {
	return sqlmock.AnyArg()
}
