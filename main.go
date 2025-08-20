package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/daniel-z-johnson/personal-weather/config"
	"github.com/daniel-z-johnson/personal-weather/controllers"
	"github.com/daniel-z-johnson/personal-weather/models"
	"github.com/daniel-z-johnson/personal-weather/templates"
	"github.com/daniel-z-johnson/personal-weather/views"
	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	
	conf, err := config.LoadConfig("config.json")
	if err != nil {
		// Try to create a basic config from environment variables if file doesn't exist
		conf = &config.Config{}
		if apiKey := os.Getenv("WEATHER_API_KEY"); apiKey != "" {
			conf.WeatherAPI.Key = apiKey
		} else {
			// for now since app won't work without a config, just panic
			panic(fmt.Errorf("config file not found and WEATHER_API_KEY not set: %w", err))
		}
		if port := os.Getenv("PORT"); port != "" {
			conf.Server.Port = port
		} else {
			conf.Server.Port = "1117"
		}
		if dbPath := os.Getenv("DATABASE_PATH"); dbPath != "" {
			conf.Database.Path = dbPath
		} else {
			conf.Database.Path = "w.db"
		}
	}
	logger.Info("Configuration loaded", "config", conf.String())
	
	db, err := sql.Open("sqlite3", conf.Database.Path)
	if err != nil {
		logger.Error("Failed to open database", slog.Any("error", err))
		panic(fmt.Errorf("Failed to open database: %w", err))
	}
	goose.SetLogger(&SlogGooseLogger{Logger: logger})
	goose.SetDialect("sqlite3")
	if err := goose.Up(db, "migrations"); err != nil {
		logger.Error("Failed to up migrations", slog.Any("error", err))
		panic(fmt.Errorf("Failed to up migrations: %w", err))
	}
	logger.Info("Personal Weather start")
	weatherAPI := &models.OpenWeatherAPI{Logger: logger, APIKey: conf.WeatherAPI.Key}
	weatherService := &models.WeatherService{DB: db, Logger: logger}
	weatherController, err := controllers.NewWeather(logger, weatherAPI, weatherService)
	if err != nil {
		// just fail at startup if something goes wrong at this point
		panic(err)
	}
	weatherController.Templates.Main =
		views.Must(views.ParseFS(templates.FS, logger, "main-layout.gohtml", "main-page.gohtml"))
	weatherController.Templates.Cities =
		views.Must(views.ParseFS(templates.FS, logger, "main-layout.gohtml", "add-city.gohtml"))

	r := chi.NewRouter()
	r.Get("/", weatherController.Main)
	r.Get("/cities", weatherController.Cities)
	r.Post("/cities", weatherController.FindCities)
	r.Post("/addCity", weatherController.AddCity)
	r.Post("/deleteCity", weatherController.DeleteCity)
	r.Get("/health", weatherController.Health)

	if err := http.ListenAndServe(":"+conf.Server.Port, r); err != nil {
		logger.Error("Failed to start server", slog.Any("error", err))
		panic(fmt.Errorf("Failed to start server: %w", err))
	}
}

type SlogGooseLogger struct {
	Logger *slog.Logger
}

func (l *SlogGooseLogger) Printf(format string, v ...interface{}) {
	l.Logger.Info("goose", "msg", fmt.Sprintf(format, v...))
}

func (l *SlogGooseLogger) Fatalf(format string, v ...interface{}) {
	l.Logger.Error("goose fatal", "msg", fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (l *SlogGooseLogger) Error(v ...interface{}) {
	l.Logger.Error("goose error", "msg", fmt.Sprint(v...))
}
