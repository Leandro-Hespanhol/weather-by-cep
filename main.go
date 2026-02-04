package main

import (
	"log"
	"net/http"
	"os"

	"github.com/lhespanhol/weather-by-cep/internal/handlers"
	"github.com/lhespanhol/weather-by-cep/internal/services"
)

func main() {
	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Get WeatherAPI key from environment variable
	weatherAPIKey := os.Getenv("WEATHER_API_KEY")
	if weatherAPIKey == "" {
		log.Fatal("WEATHER_API_KEY environment variable is required")
	}

	// Initialize services
	viaCEPService := services.NewViaCEPService()
	weatherService := services.NewWeatherService(weatherAPIKey)

	// Initialize handlers
	weatherHandler := handlers.NewWeatherHandler(viaCEPService, weatherService)

	// Setup routes
	mux := http.NewServeMux()
	mux.HandleFunc("/weather/", weatherHandler.GetWeatherByCEP)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start server
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
