-- 1_initial_migration.up.sql

-- Create the table for sensor metadata
CREATE TABLE sensor_metadata (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    location_latitude FLOAT NOT NULL,
    location_longitude FLOAT NOT NULL
);

-- Create an index for faster querying based on the location
CREATE INDEX idx_sensor_metadata_location ON sensor_metadata (location_latitude, location_longitude);
