package handlers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"weather-check/internal/services"
)

type CheckWeatherResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type WeatherHandler struct {
	service services.WeatherServiceInterface
}

func NewWeatherHandler(service services.WeatherServiceInterface) *WeatherHandler {
	return &WeatherHandler{service: service}
}

func (wh *WeatherHandler) CheckWeather(w http.ResponseWriter, r *http.Request) {
	zipcode := mux.Vars(r)["zipcode"]

	data, err := wh.service.CheckWeather(zipcode)
	if err != nil {
		if errors.Is(err, services.ERROR_INVALID_ZIPCODE) {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		if errors.Is(err, services.ERROR_CANNOT_FIND_ZIPCODE) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")

	res := CheckWeatherResponse{
		TempC: data.TempC,
		TempF: data.TempF,
		TempK: data.TempK,
	}

	json.NewEncoder(w).Encode(res)
}
