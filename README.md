# UDCO2S-SQLite API in Go

The UDCO2S-SQLite API in Go is a Docker-based system designed specifically for collecting, storing, and querying environmental sensor data from the IO-DATA [UD-CO2S](https://www.iodata.jp/product/tsushin/iot/ud-co2s/) sensor.

## Description

The UDCO2S-SQLite API in Go project provides a Dockerized microservices architecture with a Go-based API and SQLite database. The system is optimized to process data from UD-CO2S environmental sensors. The system focuses on capturing CO2 concentration, humidity levels, and temperature readings.

## Getting Started

### Prerequisites

- Docker Engine (>= 25.0.3) and Docker Compose (>= 2.24.7) installed on your system.
- An IO-DATA [UD-CO2S](https://www.iodata.jp/product/tsushin/iot/ud-co2s/) CO2 sensor connected to the system via a serial port (default `/dev/ttyACM0`).

### Installation

1. Clone this repository to your local machine.
2. Navigate to the project's root directory.
3. Run the following command to build and start the containers:

```shell
docker compose up
```

This command builds and starts the API and database services as defined in the `compose.yaml` file.

### Usage

After the containers are running, the API is available at `http://localhost:8080/sensor_data`. Retrieve sensor data by specifying `start` and `end` query parameters in ISO 8601 format.

Example request:

```http
GET http://localhost:8080/sensor_data?start=2024-03-25T00:00:00Z&end=2024-03-25T23:59:59Z
```

## API Reference

- GET /sensor_data
  - Retrieves sensor data within a specific datetime range.
  - Query Parameters:
    - `start`: Start datetime in ISO 8601 format.
    - `end`: End datetime in ISO 8601 format.
  - Responses:
    - 200: Success, with an array of sensor data objects.
    - 400: Bad Request, due to missing or invalid parameters.
    - 404: Not Found, if no data exists for the specified range.
    - 405: Method Not Allowed, if a non-GET method is used.

## Architecture

The project consists of two Docker services:

- `api`: A Go application providing an HTTP API to access sensor data from the SQLite database, accessible on port 8080.
- `db`: Manages sensor data collection and storage in the SQLite database, sharing the database file with the API service through a volume.

## Development

### Building the API

The API server is developed in Go, with dependencies managed by Go Modules. The `api` directory contains a Dockerfile outlining a multi-stage build process for compiling the application.

### Database Schema

The SQLite database stores sensor readings with fields for ID, CO2 concentration (ppm), humidity (%), temperature (Celsius), and timestamp.

## Author

- Nakurei - Initial work -
    - [GitHub](https://github.com/NakuRei)

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
