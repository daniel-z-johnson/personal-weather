package models

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
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

func (ows *OpenWeatherAPI) GetCityCoordinates(city, state, country string) ([]*GeoLocation, error) {
	city = strings.TrimSpace(city)
	state = strings.TrimSpace(state)
	country = strings.TrimSpace(country)
	if city == "" {
		return nil, errors.New("city cannot be empty")
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
		return nil, err
	}
	values := uri.Query()
	values.Set("q", locationValue)
	values.Set("limit", "5")
	values.Set("appid", ows.APIKey)
	uri.RawQuery = values.Encode()
	resp, err := http.Get(uri.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	locations := make([]*GeoLocation, 0)
	err = json.NewDecoder(resp.Body).Decode(&locations)
	if err != nil {
		return nil, err
	}
	return locations, nil
}
