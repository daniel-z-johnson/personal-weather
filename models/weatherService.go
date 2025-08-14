package models

import (
	"database/sql"
	"errors"
	"log/slog"
	"time"
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

func (ws *WeatherService) GetAllExpired() ([]GeoLocation, error) {
	query := `SELECT id, city, state, country, latitude, longitude FROM locations WHERE expires < ?`
	dateTimeNow := time.Now().Format(time.DateTime)
	rows, err := ws.DB.Query(query, dateTimeNow)
	if err != nil {
		ws.Logger.Error("Failed to get expired locations", slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()

	var locations []GeoLocation
	for rows.Next() {
		var loc GeoLocation
		err := rows.Scan(&loc.ID, &loc.Name, &loc.State, &loc.Country, &loc.Latitude, &loc.Longitude)
		if err != nil {
			ws.Logger.Error("Failed to scan expired location row", slog.String("error", err.Error()))
			return nil, err
		}
		locations = append(locations, loc)
	}

	if err = rows.Err(); err != nil {
		ws.Logger.Error("Error iterating over expired rows", slog.String("error", err.Error()))
		return nil, err
	}
	return locations, nil
}

func (ws *WeatherService) GetAll() ([]GeoLocation, error) {
	query := `SELECT city, state, country, latitude, longitude FROM locations`
	rows, err := ws.DB.Query(query)
	if err != nil {
		ws.Logger.Error("Failed to get all locations", slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()

	var locations []GeoLocation
	for rows.Next() {
		var loc GeoLocation
		err := rows.Scan(&loc.Name, &loc.State, &loc.Country, &loc.Latitude, &loc.Longitude)
		if err != nil {
			ws.Logger.Error("Failed to scan location row", slog.String("error", err.Error()))
			return nil, err
		}
		locations = append(locations, loc)
	}

	if err = rows.Err(); err != nil {
		ws.Logger.Error("Error iterating over rows", slog.String("error", err.Error()))
		return nil, err
	}

	return locations, nil
}

func (ws *WeatherService) UpdateLocation(id int, temp float64) error {
	query := `UPDATE locations SET expires = ?, temp = ? WHERE id = ?`
	dateTimeExpires := time.Now().Add(30 * time.Minute).Format(time.DateTime)
	_, err := ws.DB.Exec(query, dateTimeExpires, temp, id)
	if err != nil {
		ws.Logger.Error("Failed to update location", slog.Int("id", id), slog.String("error", err.Error()))
		return err
	}
	ws.Logger.Info("Location updated successfully", slog.Int("id", id))
	return nil
}
