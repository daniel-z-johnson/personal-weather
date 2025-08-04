package controllers

import (
	"fmt"
	"github.com/daniel-z-johnson/personal-weather/models"
	"log/slog"
	"net/http"
	"strconv"
)

type Weather struct {
	logger         *slog.Logger
	openWeatherAPI *models.OpenWeatherAPI
	Templates      struct {
		Main   Template
		Cities Template
	}
}

type LocationPageData struct {
	City      string
	State     string
	Country   string
	Latitude  float64
	Longitude float64
}

func NewWeather(logger *slog.Logger, openWeatherAPI *models.OpenWeatherAPI) (*Weather, error) {
	return &Weather{logger: logger, openWeatherAPI: openWeatherAPI}, nil
}

func (weather *Weather) Main(w http.ResponseWriter, r *http.Request) {
	weather.Templates.Main.Execute(w, r, nil)
}

func (weather *Weather) Cities(w http.ResponseWriter, r *http.Request) {
	weather.Templates.Cities.Execute(w, r, nil)
}

func (weather *Weather) AddCity(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		Form LocationPageData
	}
	var data Data
	err := r.ParseForm()
	if err != nil {
		weather.logger.Error("Failed to parse form", slog.Any("error", err))
		weather.Templates.Cities.Execute(w, r, nil, fmt.Errorf("Server issue try again later"))
		return
	}
	data.Form.City = r.FormValue("city")
	data.Form.State = r.FormValue("state")
	data.Form.Country = r.FormValue("country")
	long, err := strconv.ParseFloat(r.FormValue("longitude"), 64)
	if err != nil {
		weather.logger.Error("Failed to parse longitude", slog.Any("error", err))
		weather.Templates.Cities.Execute(w, r, nil, fmt.Errorf("Invalid longitude value"))
		return
	}
	lat, err := strconv.ParseFloat(r.FormValue("latitude"), 64)
	if err != nil {
		weather.logger.Error("Failed to parse latitude", slog.Any("error", err))
		weather.Templates.Cities.Execute(w, r, nil, fmt.Errorf("Invalid latitude value"))
		return
	}
	data.Form.Longitude = long
	data.Form.Latitude = lat
	weather.Templates.Main.Execute(w, r, nil)
}

func (weather *Weather) FindCities(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		Form      LocationPageData
		Locations []LocationPageData
	}
	err := r.ParseForm()
	if err != nil {
		weather.logger.Error("Failed to parse form", slog.Any("error", err))
		weather.Templates.Cities.Execute(w, r, nil, fmt.Errorf("Server issue try again later"))
		return
	}
	var data Data
	data.Form.City = r.FormValue("city")
	data.Form.State = r.FormValue("state")
	data.Form.Country = r.FormValue("country")
	locations, err := weather.openWeatherAPI.GetCityCoordinates(data.Form.City, data.Form.State, data.Form.Country)
	if err != nil {
		weather.logger.Error("Failed to parse form", slog.Any("error", err))
		weather.Templates.Cities.Execute(w, r, nil, fmt.Errorf("Server issue try again later"))
		return
	}
	data.Locations = make([]LocationPageData, 0)
	for _, loc := range locations {
		data.Locations = append(data.Locations, LocationPageData{
			City:      loc.Name,
			State:     loc.State,
			Country:   loc.Country,
			Latitude:  loc.Latitude,
			Longitude: loc.Longitude,
		})
	}
	weather.Templates.Cities.Execute(w, r, &data)
}
