package models

import (
	"database/sql"
	"log/slog"
)

type WeatherService struct {
	DB     *sql.DB
	Logger *slog.Logger
}

func (ws *WeatherService) SaveLocation(city, state, country string, latitude, longitude float64) error {
	query := `INSERT INTO locations (city, state, country, latitude, longitude) VALUES (?, ?, ?, ?, ?)`
	_, err := ws.DB.Exec(query, city, state, country, latitude, longitude)
	if err != nil {
		ws.Logger.Error("Failed to save location", slog.String("city", city), slog.String("state", state),
			slog.String("country", country), slog.String("error", err.Error()))
		return err
	}
	ws.Logger.Info("Location saved successfully",
		slog.String("city", city), slog.String("state", state), slog.String("country", country))
	return nil
}
