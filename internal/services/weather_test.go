package services

import (
	"context"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
)

func almostEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}

func TestConvertCelsiusToFahrenheit(t *testing.T) {
	tests := []struct {
		name     string
		celsius  float64
		expected float64
	}{
		{
			name:     "freezing point",
			celsius:  0,
			expected: 32,
		},
		{
			name:     "boiling point",
			celsius:  100,
			expected: 212,
		},
		{
			name:     "room temperature",
			celsius:  25,
			expected: 77,
		},
		{
			name:     "negative temperature",
			celsius:  -10,
			expected: 14,
		},
		{
			name:     "body temperature",
			celsius:  37,
			expected: 98.6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertCelsiusToFahrenheit(tt.celsius)
			if !almostEqual(result, tt.expected, 0.0001) {
				t.Errorf("ConvertCelsiusToFahrenheit(%v) = %v, want %v", tt.celsius, result, tt.expected)
			}
		})
	}
}

func TestConvertCelsiusToKelvin(t *testing.T) {
	tests := []struct {
		name     string
		celsius  float64
		expected float64
	}{
		{
			name:     "freezing point",
			celsius:  0,
			expected: 273,
		},
		{
			name:     "boiling point",
			celsius:  100,
			expected: 373,
		},
		{
			name:     "room temperature",
			celsius:  25,
			expected: 298,
		},
		{
			name:     "absolute zero in practice",
			celsius:  -273,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertCelsiusToKelvin(tt.celsius)
			if result != tt.expected {
				t.Errorf("ConvertCelsiusToKelvin(%v) = %v, want %v", tt.celsius, result, tt.expected)
			}
		})
	}
}

func TestWeatherService_GetTemperature(t *testing.T) {
	t.Run("successful weather lookup", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify API key is present
			if r.URL.Query().Get("key") == "" {
				t.Error("API key not present in request")
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"current": {
					"temp_c": 25.5
				}
			}`))
		}))
		defer server.Close()

		service := &WeatherService{
			baseURL:    server.URL,
			apiKey:     "test-api-key",
			httpClient: server.Client(),
		}

		weather, err := service.GetTemperature(context.Background(), "São Paulo")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if weather == nil {
			t.Fatal("expected weather data, got nil")
		}

		if weather.Current.TempC != 25.5 {
			t.Errorf("expected temp_c 25.5, got %v", weather.Current.TempC)
		}
	})

	t.Run("server error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer server.Close()

		service := &WeatherService{
			baseURL:    server.URL,
			apiKey:     "invalid-key",
			httpClient: server.Client(),
		}

		_, err := service.GetTemperature(context.Background(), "São Paulo")
		if err == nil {
			t.Error("expected error for unauthorized response")
		}
	})

	t.Run("city with special characters", func(t *testing.T) {
		var receivedQuery string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedQuery = r.URL.Query().Get("q")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"current": {"temp_c": 20.0}}`))
		}))
		defer server.Close()

		service := &WeatherService{
			baseURL:    server.URL,
			apiKey:     "test-key",
			httpClient: server.Client(),
		}

		_, err := service.GetTemperature(context.Background(), "São Paulo")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// The city name should be URL-encoded
		if receivedQuery == "" {
			t.Error("expected city to be sent in query")
		}
	})
}
