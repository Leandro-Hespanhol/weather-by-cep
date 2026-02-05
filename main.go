package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/lhespanhol/weather-by-cep/internal/handlers"
	"github.com/lhespanhol/weather-by-cep/internal/services"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Get WeatherAPI key from environment variable
	weatherAPIKey := os.Getenv("WEATHER_API_KEY")
	if weatherAPIKey == "" {
		weatherAPIKey = "2ed235e3ee27451a89231449260502"
		// log.Fatal("WEATHER_API_KEY environment variable is required")
	}

	// Initialize services
	viaCEPService := services.NewViaCEPService()
	weatherService := services.NewWeatherService(weatherAPIKey)

	// Initialize handlers
	weatherHandler := handlers.NewWeatherHandler(viaCEPService, weatherService)

	// Setup routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("weather-by-cep service is running"))
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/weather/", weatherHandler.GetWeatherByCEP)

	// Bind to 0.0.0.0 explicitly for container environments
	addr := fmt.Sprintf("0.0.0.0:%s", port)
	log.Printf("Server listening on %s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
