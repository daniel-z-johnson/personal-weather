package main

import (
	"fmt"
	"github.com/daniel-z-johnson/personal-weather/config"
	"github.com/daniel-z-johnson/personal-weather/controllers"
	"github.com/daniel-z-johnson/personal-weather/models"
	"github.com/daniel-z-johnson/personal-weather/templates"
	"github.com/daniel-z-johnson/personal-weather/views"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"os"
)

func main() {
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

	r := chi.NewRouter()
	r.Get("/", weatherController.Main)

	http.ListenAndServe(":1117", r)
}
