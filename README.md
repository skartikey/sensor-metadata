```markdown
# Sensor Metadata API

The Sensor Metadata API is a JSON REST API that allows storing and querying sensor metadata. It provides endpoints for creating, retrieving, updating, and querying sensor metadata.

## Features

- Store sensor metadata including name, location (GPS position), and tags.
- Retrieve sensor metadata by name.
- Update sensor metadata.
- Find the nearest sensor based on a given location.

## Technologies Used

- Go programming language
- PostgreSQL database
- Gorilla Mux (HTTP router)
- pq (PostgreSQL driver)
- Docker (optional)

## Prerequisites

- Go 1.16 or higher installed
- PostgreSQL database running

## Getting Started

1. Clone the repository:

```bash
git clone https://github.com/your-username/sensor-metadata-api.git
cd sensor-metadata-api
```

2. Set up the database:

- Create a PostgreSQL database.
- Update the database connection string in `main.go` with your database details.

3. Build and run the application:

```bash
go build
./sensor-metadata-api
```

4. The API server should now be running at `http://localhost:8080`.

## API Endpoints

### Create Sensor Metadata

**URL:** `/sensors`

**Method:** `POST`

**Request Body:**

```json
{
  "name": "Sensor1",
  "location": {
    "latitude": 123.456,
    "longitude": 789.012
  },
  "tags": ["tag1", "tag2"]
}
```

**Response:**

- Status Code: `201 Created`
- Response Body: Empty

### Get Sensor Metadata

**URL:** `/sensors/{name}`

**Method:** `GET`

**Response:**

- Status Code: `200 OK`
- Response Body:

```json
{
  "id": 1,
  "name": "Sensor1",
  "location": {
    "latitude": 123.456,
    "longitude": 789.012
  },
  "tags": ["tag1", "tag2"]
}
```

### Update Sensor Metadata

**URL:** `/sensors/{name}`

**Method:** `PUT`

**Request Body:**

```json
{
  "name": "Sensor1",
  "location": {
    "latitude": 111.222,
    "longitude": 333.444
  },
  "tags": ["tag3", "tag4"]
}
```

**Response:**

- Status Code: `204 No Content`
- Response Body: Empty

### Get Nearest Sensor Metadata

**URL:** `/sensors/nearest?latitude={latitude}&longitude={longitude}`

**Method:** `GET`

**Response:**

- Status Code: `200 OK`
- Response Body:

```json
{
  "id": 2,
  "name": "Sensor2",
  "location": {
    "latitude": 222.333,
    "longitude": 444.555
  },
  "tags": ["tag5", "tag6"]
}
```

## Testing

To run the tests, use the following command:

```bash
go test ./...
```

## Docker

You can also run the application using Docker. Dockerize the application with the following steps:

1. Build the Docker image:

```bash
docker build -t sensor-metadata-api .
```

2. Run the Docker container:

```bash
docker run -p 8080:8080 sensor-metadata-api
```

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvement, please feel free to submit a pull request.

Before contributing, please ensure that you:

- Follow the existing coding style and conventions.
- Write clear commit messages.
- Test your changes thoroughly.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgements

The Sensor Metadata API is built using various open-source libraries and frameworks. We would like to acknowledge and thank the following:

- [Gorilla Mux](https://github.com/gorilla/mux) - HTTP router for Go.
- [pq](https://github.com/lib/pq) - PostgreSQL driver for Go.
- [Docker](https://www.docker.com/) - Containerization platform.

## Contact

For any inquiries or questions, please contact:

Sunil Kartikey
Email: s.kartikey@gmail.com
