package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lhespanhol/weather-by-cep/internal/models"
	"github.com/lhespanhol/weather-by-cep/internal/services"
)

func TestWeatherHandler_GetWeatherByCEP_InvalidCEP(t *testing.T) {
	tests := []struct {
		name           string
		cep            string
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:           "CEP too short",
			cep:            "0131010",
			expectedStatus: http.StatusUnprocessableEntity,
			expectedMsg:    "invalid zipcode",
		},
		{
			name:           "CEP too long",
			cep:            "013101000",
			expectedStatus: http.StatusUnprocessableEntity,
			expectedMsg:    "invalid zipcode",
		},
		{
			name:           "CEP with letters",
			cep:            "0131010a",
			expectedStatus: http.StatusUnprocessableEntity,
			expectedMsg:    "invalid zipcode",
		},
		{
			name:           "empty CEP",
			cep:            "",
			expectedStatus: http.StatusUnprocessableEntity,
			expectedMsg:    "invalid zipcode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viaCEPService := services.NewViaCEPService()
			weatherService := services.NewWeatherService("test-key")
			handler := NewWeatherHandler(viaCEPService, weatherService)

			req := httptest.NewRequest(http.MethodGet, "/weather/"+tt.cep, nil)
			rec := httptest.NewRecorder()

			handler.GetWeatherByCEP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			var response models.ErrorResponse
			if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if response.Message != tt.expectedMsg {
				t.Errorf("expected message %q, got %q", tt.expectedMsg, response.Message)
			}
		})
	}
}

func TestWeatherHandler_GetWeatherByCEP_CEPWithDash(t *testing.T) {
	// Setup mock servers
	viaCEPServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"cep": "01310-100",
			"localidade": "São Paulo",
			"uf": "SP"
		}`))
	}))
	defer viaCEPServer.Close()

	weatherServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"current": {"temp_c": 25.0}}`))
	}))
	defer weatherServer.Close()

	// Create services with mock servers
	viaCEPService := services.NewViaCEPServiceWithClient(viaCEPServer.URL, viaCEPServer.Client())
	weatherService := services.NewWeatherServiceWithClient(weatherServer.URL, "test-key", weatherServer.Client())

	handler := NewWeatherHandler(viaCEPService, weatherService)

	// Test with CEP containing dash
	req := httptest.NewRequest(http.MethodGet, "/weather/01310-100", nil)
	rec := httptest.NewRecorder()

	handler.GetWeatherByCEP(rec, req)

	// Should successfully process CEP with dash
	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d. Body: %s", http.StatusOK, rec.Code, rec.Body.String())
	}
}

func TestWeatherHandler_GetWeatherByCEP_Success(t *testing.T) {
	// Setup mock servers
	viaCEPServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"cep": "01310-100",
			"localidade": "São Paulo",
			"uf": "SP"
		}`))
	}))
	defer viaCEPServer.Close()

	weatherServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"current": {"temp_c": 28.5}}`))
	}))
	defer weatherServer.Close()

	// Create services with mock servers
	viaCEPService := services.NewViaCEPServiceWithClient(viaCEPServer.URL, viaCEPServer.Client())
	weatherService := services.NewWeatherServiceWithClient(weatherServer.URL, "test-key", weatherServer.Client())

	handler := NewWeatherHandler(viaCEPService, weatherService)

	req := httptest.NewRequest(http.MethodGet, "/weather/01310100", nil)
	rec := httptest.NewRecorder()

	handler.GetWeatherByCEP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d. Body: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	var response models.WeatherResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Verify temperatures
	if response.TempC != 28.5 {
		t.Errorf("expected TempC 28.5, got %v", response.TempC)
	}

	expectedF := 28.5*1.8 + 32 // 83.3
	if response.TempF != expectedF {
		t.Errorf("expected TempF %v, got %v", expectedF, response.TempF)
	}

	expectedK := 28.5 + 273 // 301.5
	if response.TempK != expectedK {
		t.Errorf("expected TempK %v, got %v", expectedK, response.TempK)
	}
}

func TestWeatherHandler_GetWeatherByCEP_NotFound(t *testing.T) {
	// Setup mock servers - viaCEP returns error for not found
	viaCEPServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"erro": true}`))
	}))
	defer viaCEPServer.Close()

	weatherServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer weatherServer.Close()

	// Create services with mock servers
	viaCEPService := services.NewViaCEPServiceWithClient(viaCEPServer.URL, viaCEPServer.Client())
	weatherService := services.NewWeatherServiceWithClient(weatherServer.URL, "test-key", weatherServer.Client())

	handler := NewWeatherHandler(viaCEPService, weatherService)

	req := httptest.NewRequest(http.MethodGet, "/weather/99999999", nil)
	rec := httptest.NewRecorder()

	handler.GetWeatherByCEP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}

	var response models.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Message != "can not find zipcode" {
		t.Errorf("expected message 'can not find zipcode', got %q", response.Message)
	}
}
