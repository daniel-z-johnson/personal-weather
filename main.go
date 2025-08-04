package main

import (
	"database/sql"
	"fmt"
	"github.com/daniel-z-johnson/personal-weather/config"
	"github.com/daniel-z-johnson/personal-weather/controllers"
	"github.com/daniel-z-johnson/personal-weather/models"
	"github.com/daniel-z-johnson/personal-weather/templates"
	"github.com/daniel-z-johnson/personal-weather/views"
	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	db, err := sql.Open("sqlite3", "w.db")
	if err != nil {
		slog.Error("Failed to open database", slog.Any("error", err))
		panic(fmt.Errorf("Failed to open database: %w", err))
	}
	goose.SetDialect("sqlite3")
	if err := goose.Up(db, "migrations"); err != nil {
		slog.Error("Failed to up migrations", slog.Any("error", err))
		panic(fmt.Errorf("Failed to up migrations: %w", err))
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("Personal Weather start")
	conf, err := config.LoadConfig("config.json")
	if err != nil {
		// for now since app won't work without a config, just panic
		panic(err)
	}
	logger.Info(fmt.Sprintf(conf.String()))

	weatherService := &models.OpenWeatherAPI{Logger: logger, APIKey: conf.WeatherAPI.Key}
	weatherController, err := controllers.NewWeather(logger, weatherService)
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
