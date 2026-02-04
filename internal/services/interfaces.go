package services

import (
	"context"
	"net/http"
	"time"

	"github.com/lhespanhol/weather-by-cep/internal/models"
)

// CEPService defines the interface for CEP lookup services
type CEPService interface {
	GetLocation(ctx context.Context, cep string) (*models.ViaCEPResponse, error)
}

// WeatherServiceInterface defines the interface for weather services
type WeatherServiceInterface interface {
	GetTemperature(ctx context.Context, city string) (*models.WeatherAPIResponse, error)
}

// NewViaCEPServiceWithClient creates a new ViaCEP service with custom base URL and client
func NewViaCEPServiceWithClient(baseURL string, client *http.Client) *ViaCEPService {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	return &ViaCEPService{
		baseURL:    baseURL,
		httpClient: client,
	}
}

// NewWeatherServiceWithClient creates a new Weather service with custom base URL and client
func NewWeatherServiceWithClient(baseURL, apiKey string, client *http.Client) *WeatherService {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	return &WeatherService{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: client,
	}
}
