package controllers

import (
	"log/slog"
	"net/http"
)

type Weather struct {
	logger    *slog.Logger
	Templates struct {
		Main Template
	}
}

func NewWeather(logger *slog.Logger) (*Weather, error) {
	return &Weather{logger: logger}, nil
}

func (weather *Weather) Main(w http.ResponseWriter, r *http.Request) {
	weather.Templates.Main.Execute(w, r, nil)
}
