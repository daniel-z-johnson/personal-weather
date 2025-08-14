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
const baseTemperatureURL = "https://api.openweathermap.org/data/3.0/onecall"

type OpenWeatherAPI struct {
	APIKey string
	Logger *slog.Logger
}

type GeoLocation struct {
	ID         int                `json:"_"`
	Name       string             `json:"name"`
	LocalNames map[string]*string `json:"local_names"`
	Latitude   float64            `json:"lat"`
	Longitude  float64            `json:"lon"`
	Country    string             `json:"country"`
	State      string             `json:"state,omitempty"`
}

type TemperatureData struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
	Current   struct {
		Temp float64 `json:"temp"`
	}
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

func (ows *OpenWeatherAPI) GetTemperature(lat, lon float64) (float64, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	uri, err := url.Parse(baseTemperatureURL)
	if err != nil {
		ows.Logger.Error("Failed to parse GetTemperature API URL",
			slog.String("url", baseTemperatureURL), slog.Any("error", err))
		return 0, fmt.Errorf("failed to parse GetTemperature API URL: %w", err)
	}
	values := uri.Query()
	values.Set("lat", fmt.Sprintf("%f", lat))
	values.Set("lon", fmt.Sprintf("%f", lon))
	values.Set("appid", ows.APIKey)
	values.Set("units", "imperial")
	values.Set("exclude", "minutely,hourly,daily,alerts")
	uri.RawQuery = values.Encode()
	resp, err := client.Get(uri.String())
	if err != nil {
		ows.Logger.Error("Request failed", slog.String("error", err.Error()))
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		ows.Logger.Error("Did not get a 200 OK response from GetTemperature API",
			slog.String("status", resp.Status))
		return 0, fmt.Errorf("status code %d Error: %w", resp.StatusCode, err)
	}
	var tempData TemperatureData
	err = json.NewDecoder(resp.Body).Decode(&tempData)
	if err != nil {
		ows.Logger.Error("failed to decode response body", slog.String("error", err.Error()))
		return 0, err
	}
	return tempData.Current.Temp, nil
}
