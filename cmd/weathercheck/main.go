package main

import (
	"log"
	"net/http"
	"weather-check/internal/web"
)

func main() {
	mux := http.NewServeMux()

	webWeatherCheckHandler := web.NewWeatherCheckHandler()

	mux.HandleFunc("/weather-check", webWeatherCheckHandler.GetInfo)

	log.Println("Server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
