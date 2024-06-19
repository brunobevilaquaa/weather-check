package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"weather-check/internal/adapters/api"
	"weather-check/internal/adapters/handlers"
	"weather-check/internal/services"
)

func main() {
	router := mux.NewRouter()

	apiClient := api.NewClient()

	weatherService := services.NewWeatherService(apiClient)

	weatherHandler := handlers.NewWeatherHandler(weatherService)

	router.HandleFunc("/api/v1/weather-check/{zipcode}", weatherHandler.CheckWeather).Methods(http.MethodGet)

	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
