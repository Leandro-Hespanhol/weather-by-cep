package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/lhespanhol/weather-by-cep/internal/models"
)

// WeatherService handles weather lookups
type WeatherService struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewWeatherService creates a new weather service
func NewWeatherService(apiKey string) *WeatherService {
	return &WeatherService{
		baseURL: "https://api.weatherapi.com/v1",
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetTemperature fetches the current temperature for a location
func (s *WeatherService) GetTemperature(ctx context.Context, city string) (*models.WeatherAPIResponse, error) {
	// Encode the city name to handle special characters
	encodedCity := url.QueryEscape(city)
	reqURL := fmt.Sprintf("%s/current.json?key=%s&q=%s", s.baseURL, s.apiKey, encodedCity)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weatherAPI returned status %d", resp.StatusCode)
	}

	var weatherResp models.WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &weatherResp, nil
}

// ConvertCelsiusToFahrenheit converts Celsius to Fahrenheit
// Formula: F = C * 1.8 + 32
func ConvertCelsiusToFahrenheit(celsius float64) float64 {
	return celsius*1.8 + 32
}

// ConvertCelsiusToKelvin converts Celsius to Kelvin
// Formula: K = C + 273
func ConvertCelsiusToKelvin(celsius float64) float64 {
	return celsius + 273
}
