package models

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const baseGeoLocatorURL = "http://api.openweathermap.org/geo/1.0/direct"

type OpenWeatherAPI struct {
	APIKey string
	Logger *slog.Logger
}

type GeoLocation struct {
	Name       string             `json:"name"`
	LocalNames map[string]*string `json:"local_names"`
	Latitude   float64            `json:"lat"`
	Longitude  float64            `json:"lon"`
	Country    string             `json:"country"`
	State      string             `json:"state,omitempty"`
}

func (ows *OpenWeatherAPI) GetCityCoordinates(city, state, country string) ([]GeoLocation, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	city = strings.TrimSpace(city)
	state = strings.TrimSpace(state)
	country = strings.TrimSpace(country)
	if city == "" {
		ows.Logger.Error("City cannot be empty")
		return nil, fmt.Errorf("city cannot be empty")
	}
	locationValue := city
	if state != "" {
		locationValue += "," + state
	}
	if country != "" {
		locationValue += "," + country
	}
	uri, err := url.Parse(baseGeoLocatorURL)
	if err != nil {
		ows.Logger.Error("Failed to parse FindCities API URL",
			slog.String("url", baseGeoLocatorURL), slog.Any("error", err))
		return nil, fmt.Errorf("failed to parse FindCities API URL: %w", err)
	}
	values := uri.Query()
	values.Set("q", locationValue)
	values.Set("limit", "5")
	values.Set("appid", ows.APIKey)
	uri.RawQuery = values.Encode()
	resp, err := client.Get(uri.String())
	if err != nil {
		ows.Logger.Error("Request failed", slog.String("error", err.Error()))
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		ows.Logger.Error("Did not get a 200 OK response from FindCities API",
			slog.String("status", resp.Status))
		return nil, fmt.Errorf("status code %d Error: %w", resp.StatusCode, err)
	}
	locations := make([]GeoLocation, 0)
	err = json.NewDecoder(resp.Body).Decode(&locations)
	if err != nil {
		ows.Logger.Error("failed to decode response body", slog.String("error", err.Error()))
		return nil, err
	}
	return locations, nil
}
