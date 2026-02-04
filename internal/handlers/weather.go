package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/lhespanhol/weather-by-cep/internal/models"
	"github.com/lhespanhol/weather-by-cep/internal/services"
)

// WeatherHandler handles weather-related HTTP requests
type WeatherHandler struct {
	viaCEPService  *services.ViaCEPService
	weatherService *services.WeatherService
}

// NewWeatherHandler creates a new weather handler
func NewWeatherHandler(viaCEPService *services.ViaCEPService, weatherService *services.WeatherService) *WeatherHandler {
	return &WeatherHandler{
		viaCEPService:  viaCEPService,
		weatherService: weatherService,
	}
}

// GetWeatherByCEP handles GET /weather/{cep}
func (h *WeatherHandler) GetWeatherByCEP(w http.ResponseWriter, r *http.Request) {
	// Extract CEP from URL path
	path := strings.TrimPrefix(r.URL.Path, "/weather/")
	cep := strings.TrimSpace(path)

	// Remove any dashes from CEP (e.g., "01310-100" -> "01310100")
	cep = strings.ReplaceAll(cep, "-", "")

	// Validate CEP format
	if !services.ValidateCEP(cep) {
		h.respondWithError(w, http.StatusUnprocessableEntity, "invalid zipcode")
		return
	}

	// Fetch location from viaCEP
	location, err := h.viaCEPService.GetLocation(r.Context(), cep)
	if err != nil {
		log.Printf("Error fetching location: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if location == nil {
		h.respondWithError(w, http.StatusNotFound, "can not find zipcode")
		return
	}

	// Fetch weather for the location
	weather, err := h.weatherService.GetTemperature(r.Context(), location.Localidade)
	if err != nil {
		log.Printf("Error fetching weather: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	// Convert temperatures
	tempC := weather.Current.TempC
	tempF := services.ConvertCelsiusToFahrenheit(tempC)
	tempK := services.ConvertCelsiusToKelvin(tempC)

	// Respond with success
	response := models.WeatherResponse{
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}

	h.respondWithJSON(w, http.StatusOK, response)
}

func (h *WeatherHandler) respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ErrorResponse{Message: message})
}

func (h *WeatherHandler) respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
