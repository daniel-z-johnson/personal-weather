package controllers

import (
	"github.com/daniel-z-johnson/personal-weather/models"
	"log/slog"
	"net/http"
)

type Weather struct {
	logger         *slog.Logger
	openWeatherAPI *models.OpenWeatherAPI
	Templates      struct {
		Main Template
	}
}

func NewWeather(logger *slog.Logger, openWeatherAPI *models.OpenWeatherAPI) (*Weather, error) {
	return &Weather{logger: logger, openWeatherAPI: openWeatherAPI}, nil
}

func (weather *Weather) Main(w http.ResponseWriter, r *http.Request) {
	weather.Templates.Main.Execute(w, r, nil)
}
