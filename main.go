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
	db, err := sql.Open("sqlite3", "w.db")
	if err != nil {
		slog.Error("Failed to open database", slog.Any("error", err))
		panic(fmt.Errorf("Failed to open database: %w", err))
	}
	goose.SetDialect("sqlite3")
	goose.SetLogger(&SlogGooseLogger{Logger: logger})
	if err := goose.Up(db, "migrations"); err != nil {
		slog.Error("Failed to up migrations", slog.Any("error", err))
		panic(fmt.Errorf("Failed to up migrations: %w", err))
	}
	logger.Info("Personal Weather start")
	conf, err := config.LoadConfig("config.json")
	if err != nil {
		// for now since app won't work without a config, just panic
		panic(err)
	}
	logger.Info(fmt.Sprintf(conf.String()))
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

	if err := http.ListenAndServe(":1117", r); err != nil {
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
