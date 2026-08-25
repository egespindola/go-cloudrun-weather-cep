# go-cloudrun-weather-cep

Receives a Brazilian zipcode (CEP), looks up its location, and returns the current
weather at that location in Celsius, Fahrenheit, and Kelvin.

Published on Google Cloud Run: [googlecloud-get-weather](https://weatherapi-328385283436.southamerica-east1.run.app/cep/01001000)

## Requirements

- Go 1.25+
- Docker (for running the containerized server and the e2e tests)

## Configuration

Copy `.env.sample` to `.env` and adjust as needed:

```
PORT=3000
CEP_CONNECTOR_URL=https://brasilapi.com.br/api/cep/v2/{cep}
WEATHER_CONNECTOR_URL=https://api.open-meteo.com/v1/forecast?latitude={latitude}&longitude={longitude}&current_weather=true
```

## Running locally

```
go mod tidy
make goserver
```

## Running with Docker

Start the containerized server:

```
make server up
```

Stop and remove the container:

```
make server stop
```

Stop the container, remove the container, and remove the built image:

```
make server down
```

## Automated tests

Runs the app container and executes the end-to-end test suite against it:

```
make test
```

## API

### `GET /cep/:zipcode`

`zipcode` must be an 8-digit number.

#### Response

```json
{
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.65
}
```
