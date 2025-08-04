package models

import (
	"database/sql"
	"errors"
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

func (ws *WeatherService) GetLocation(city, state, country string) (*GeoLocation, error) {
	query := `SELECT latitude, longitude FROM locations WHERE city = ? AND state = ? AND country = ?`
	row := ws.DB.QueryRow(query, city, state, country)
	var location GeoLocation
	location.Name = city
	location.State = state
	location.Country = country
	err := row.Scan(&location.Latitude, &location.Longitude)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ws.Logger.Warn("No location found", slog.String("city", city), slog.String("state", state),
				slog.String("country", country))
			return nil, nil
		}
		ws.Logger.Error("Failed to get location", slog.String("error", err.Error()))
		return nil, err
	}
	return &location, nil
}
