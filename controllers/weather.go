package controllers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/daniel-z-johnson/personal-weather/models"
)

type Weather struct {
	logger         *slog.Logger
	openWeatherAPI *models.OpenWeatherAPI
	weatherService *models.WeatherService
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

type LocationTemp struct {
	City    string
	State   string
	Country string
	TempF   string
	TempC   string
}

func NewWeather(logger *slog.Logger, openWeatherAPI *models.OpenWeatherAPI, openWeatherService *models.WeatherService) (*Weather, error) {
	return &Weather{logger: logger, openWeatherAPI: openWeatherAPI, weatherService: openWeatherService}, nil
}

func (weather *Weather) Main(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		Locations []LocationTemp
		Errors    []error
	}
	expired, err := weather.weatherService.GetAllExpired()
	if err != nil {
		weather.logger.Error("Failed to get expired locations", slog.Any("error", err))
		weather.Templates.Main.Execute(w, r, nil, fmt.Errorf("server issue try again later"))
		return
	}
	for _, v := range expired {
		temp, err := weather.openWeatherAPI.GetTemperature(v.Latitude, v.Longitude)
		if err != nil {
			weather.logger.Error("Failed to get temperature for expired location", slog.Any("error", err),
				slog.String("city", v.Name), slog.String("state", v.State), slog.String("country", v.Country),
				slog.Float64("latitude", v.Latitude), slog.Float64("longitude", v.Longitude))
			continue // skip this location if we can't get the temperature
		}
		err = weather.weatherService.UpdateLocation(v.ID, temp)
		if err != nil {
			weather.logger.Error("Failed to update expired location", slog.Any("error", err),
				slog.String("city", v.Name), slog.String("state", v.State), slog.String("country", v.Country),
				slog.Float64("latitude", v.Latitude), slog.Float64("longitude", v.Longitude))
			continue // skip this location if we can't update it
		}
	}
	allLocations, err := weather.weatherService.GetAll()
	if err != nil {
		weather.logger.Error("Failed to get all locations after updating expired", slog.Any("error", err))
		weather.Templates.Main.Execute(w, r, nil, fmt.Errorf("server issue try again later"))
		return
	}
	locationTemps := make([]LocationTemp, 0)
	for _, v := range allLocations {
		var locationTemp LocationTemp
		locationTemp.City = v.City
		locationTemp.State = v.State
		locationTemp.Country = v.Country
		locationTemp.TempF = fmt.Sprintf("%.f", v.Temperature)
		locationTemp.TempC = fmt.Sprintf("%.f", (v.Temperature-32)*5/9)
		locationTemps = append(locationTemps, locationTemp)
	}

	weather.Templates.Main.Execute(w, r, &Data{Locations: locationTemps})
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
	
	// Basic input validation and sanitization
	data.Form.City = strings.TrimSpace(data.Form.City)
	data.Form.State = strings.TrimSpace(data.Form.State)
	data.Form.Country = strings.TrimSpace(data.Form.Country)
	
	if data.Form.City == "" {
		weather.logger.Error("City name is required")
		weather.Templates.Cities.Execute(w, r, nil, fmt.Errorf("City name is required"))
		return
	}
	if len(data.Form.City) > 100 {
		weather.logger.Error("City name too long", slog.String("city", data.Form.City))
		weather.Templates.Cities.Execute(w, r, nil, fmt.Errorf("City name must be 100 characters or less"))
		return
	}
	long, err := strconv.ParseFloat(r.FormValue("longitude"), 64)
	if err != nil {
		weather.logger.Error("Failed to parse longitude", slog.Any("error", err))
		weather.Templates.Cities.Execute(w, r, nil, fmt.Errorf("Invalid longitude value"))
		return
	}
	if long < -180 || long > 180 {
		weather.logger.Error("Longitude out of range", slog.Float64("longitude", long))
		weather.Templates.Cities.Execute(w, r, nil, fmt.Errorf("Longitude must be between -180 and 180"))
		return
	}
	lat, err := strconv.ParseFloat(r.FormValue("latitude"), 64)
	if err != nil {
		weather.logger.Error("Failed to parse latitude", slog.Any("error", err))
		weather.Templates.Cities.Execute(w, r, nil, fmt.Errorf("Invalid latitude value"))
		return
	}
	if lat < -90 || lat > 90 {
		weather.logger.Error("Latitude out of range", slog.Float64("latitude", lat))
		weather.Templates.Cities.Execute(w, r, nil, fmt.Errorf("Latitude must be between -90 and 90"))
		return
	}
	data.Form.Longitude = long
	data.Form.Latitude = lat
	err = weather.weatherService.SaveLocation(data.Form.City, data.Form.State, data.Form.Country, data.Form.Latitude, data.Form.Longitude)
	if err != nil {
		weather.logger.Error("Failed to save location", slog.Any("error", err))
		weather.Templates.Cities.Execute(w, r, nil, fmt.Errorf("Failed to save location: %v", err))
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
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
	
	// Basic input validation and sanitization
	data.Form.City = strings.TrimSpace(data.Form.City)
	data.Form.State = strings.TrimSpace(data.Form.State)
	data.Form.Country = strings.TrimSpace(data.Form.Country)
	
	if data.Form.City == "" {
		weather.logger.Error("City name is required for search")
		weather.Templates.Cities.Execute(w, r, nil, fmt.Errorf("City name is required"))
		return
	}
	
	locations, err := weather.openWeatherAPI.GetCityCoordinates(data.Form.City, data.Form.State, data.Form.Country)
	if err != nil {
		weather.logger.Error("Failed to get city coordinates", slog.Any("error", err))
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

func (weather *Weather) DeleteCity(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		weather.logger.Error("Failed to parse form", slog.Any("error", err))
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	
	city := strings.TrimSpace(r.FormValue("city"))
	state := strings.TrimSpace(r.FormValue("state"))
	country := strings.TrimSpace(r.FormValue("country"))
	
	if city == "" {
		weather.logger.Error("City name is required for deletion")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	
	err = weather.weatherService.DeleteLocation(city, state, country)
	if err != nil {
		weather.logger.Error("Failed to delete location", slog.Any("error", err))
	}
	
	http.Redirect(w, r, "/", http.StatusFound)
}

func (weather *Weather) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"personal-weather"}`))
}
