package app

// SensorMetadata represents the structure of sensor metadata.
type SensorMetadata struct {
	ID       int      `json:"id"`
	Name     string   `json:"name" validate:"required"`
	Location Location `json:"location" validate:"required"`
	Tags     []string `json:"tags"`
	Distance float64  `json:"distance,omitempty"`
}

// Location represents the GPS position of a sensor.
type Location struct {
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}

type MapboxFeature struct {
	Geometry struct {
		Coordinates []float64 `json:"coordinates"`
	} `json:"geometry"`
	PlaceName string `json:"place_name"`
}

type MapboxResponse struct {
	Features []MapboxFeature `json:"features"`
}
