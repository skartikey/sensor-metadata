-- 2_add_tags_column.up.sql

-- Add the tags column to the sensor_metadata table
ALTER TABLE sensor_metadata
ADD COLUMN tags VARCHAR(255)[] DEFAULT ARRAY[]::VARCHAR(255)[];
