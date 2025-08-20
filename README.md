# Personal Weather

A Go web application that allows you to track current weather conditions for multiple cities around the world. Built with Go, SQLite, and the OpenWeatherMap API.

## Features

- 🌍 Add cities from anywhere in the world using city search
- 🌡️ View current temperatures in both Fahrenheit and Celsius  
- 🎨 Clean, responsive web interface with Tailwind CSS
- 🗄️ SQLite database for persistent city storage
- 🐳 Docker support for easy deployment

## Prerequisites

- Go 1.24.4 or later
- OpenWeatherMap API key (free tier available)
- SQLite3 (automatically included via Go driver)

## Setup

### 1. Clone the Repository

```bash
git clone https://github.com/daniel-z-johnson/personal-weather.git
cd personal-weather
```

### 2. Get OpenWeatherMap API Key

1. Visit [OpenWeatherMap](https://openweathermap.org/api) and sign up for a free account
2. Generate an API key from your account dashboard

### 3. Configure the Application

Create a `config.json` file in the project root:

```json
{
    "weatherAPI": {
        "key": "your_openweathermap_api_key_here"
    }
}
```

You can also copy and modify the example configuration:

```bash
cp explample-cong.json config.json
# Edit config.json with your API key
```

### 4. Install Dependencies

```bash
go mod download
```

## Running the Application

### Local Development

```bash
go run .
```

The application will start on `http://localhost:1117`

### Docker

Build and run using Docker:

```bash
# Build the image
docker build -t personal-weather .

# Run the container
docker run -p 1117:1117 personal-weather
```

**Note:** Make sure your `config.json` file is present before building the Docker image, as it gets copied during the build process.

## Usage

### Web Interface

Navigate to `http://localhost:1117` in your browser to access the web interface.

#### Adding Cities

1. Click on "Cities" in the navigation
2. Enter city name (required)
3. Optionally specify state code (US only) and country code (ISO 3166-1 alpha-2)
4. Click "Add city" to search for matching locations
5. Select the correct location from the search results
6. Click "Add This One" to save the city

#### Viewing Weather

- Visit the home page to see current temperatures for all your saved cities
- Temperatures are displayed in both Fahrenheit and Celsius
- Data is automatically refreshed when it expires

### API Endpoints

The application exposes the following HTTP endpoints:

- `GET /` - Main weather dashboard showing all saved cities
- `GET /cities` - City management page for adding new locations  
- `POST /cities` - Search for cities by name, state, and country
- `POST /addCity` - Add a selected city to your saved locations

## Database Schema

The application uses SQLite with the following tables:

### locations
- `id` (INTEGER PRIMARY KEY) - Unique identifier
- `City` (TEXT) - City name
- `State` (TEXT) - State/province (optional)
- `Country` (TEXT) - Country code
- `Latitude` (REAL) - Geographic latitude
- `Longitude` (REAL) - Geographic longitude  
- `temp` (REAL) - Current temperature in Fahrenheit
- `expires` (TEXT) - Temperature data expiration timestamp

## Development

### Project Structure

```
├── main.go                 # Application entry point
├── config/                 # Configuration loading
├── controllers/            # HTTP handlers and routing logic
├── models/                 # Data models and API integrations
├── views/                  # Template rendering utilities
├── templates/              # HTML templates
├── migrations/             # Database migration files
├── docs/                   # Documentation and diagrams
├── Dockerfile             # Container build configuration
└── README.md              # This file
```

### Key Components

- **OpenWeatherMap Integration**: Geocoding API for city search and One Call API for weather data
- **Database Migrations**: Managed with [Goose](https://github.com/pressly/goose) 
- **HTTP Router**: [Chi](https://github.com/go-chi/chi) for lightweight routing
- **Templates**: Go HTML templates with Tailwind CSS styling
- **Logging**: Structured JSON logging with Go's `slog` package

### Building and Testing

```bash
# Build the application
go build -o personal-weather

# Run database migrations manually (optional, done automatically on startup)
go run main.go  # Migrations run on application start

# Check dependencies
go mod tidy
```

### Environment Variables

The application currently loads configuration from `config.json`. Environment variable support could be added in future versions.

## API Integration

This application integrates with the OpenWeatherMap API:

- **Geocoding API**: Used to search for cities and get coordinates
- **One Call API 3.0**: Used to retrieve current weather data
- **Rate Limits**: Free tier allows 1,000 calls/day
- **Data Updates**: Temperature data expires and is automatically refreshed

## License

This project is available under the MIT License. See the LICENSE file for details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

**Note**: The maintainer may not actively review or approve pull requests. For ongoing development, it's recommended to continue work on your own fork.

## Support

For issues and questions:
1. Check existing issues on GitHub
2. Create a new issue with detailed information
3. Include configuration details (without API keys)