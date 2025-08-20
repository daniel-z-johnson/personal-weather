# Personal Weather

A simple Go web application for tracking weather in multiple cities using the OpenWeatherMap API.

## Features

- Add cities and track their current weather
- Delete cities from your tracking list
- Automatic weather data refresh (30-minute cache)
- Environment variable configuration support
- Health check endpoint
- Input validation and error handling

## Setup

1. Get an API key from [OpenWeatherMap](https://openweathermap.org/api)

2. Create a `config.json` file based on `example-config.json`:
   ```json
   {
       "weatherAPI": {
           "key": "your-api-key-here"
       },
       "server": {
           "port": "1117"
       },
       "database": {
           "path": "w.db"
       }
   }
   ```

3. Or use environment variables:
   ```bash
   export WEATHER_API_KEY="your-api-key-here"
   export PORT="1117"
   export DATABASE_PATH="w.db"
   ```

## Running

### With Go:
```bash
go build .
./personal-weather
```

### With Docker:
```bash
docker build -t personal-weather .
docker run -p 1117:1117 -e WEATHER_API_KEY="your-key" personal-weather
```

## Endpoints

- `GET /` - Main weather dashboard
- `GET /cities` - Add/search cities page
- `POST /cities` - Search for cities
- `POST /addCity` - Add a city to tracking
- `POST /deleteCity` - Remove a city from tracking
- `GET /health` - Health check endpoint

## Testing

```bash
go test ./...
```

## Improvements Made

- Fixed typos and naming inconsistencies
- Added proper error handling and input validation
- Added city deletion functionality
- Enhanced configuration with environment variable support
- Added health check endpoint
- Improved UI with error display and delete buttons
- Added basic test coverage
- Enhanced security with input sanitization