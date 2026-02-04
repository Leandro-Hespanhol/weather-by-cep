package services

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidateCEP(t *testing.T) {
	tests := []struct {
		name     string
		cep      string
		expected bool
	}{
		{
			name:     "valid CEP with 8 digits",
			cep:      "01310100",
			expected: true,
		},
		{
			name:     "valid CEP all zeros",
			cep:      "00000000",
			expected: true,
		},
		{
			name:     "invalid CEP with less than 8 digits",
			cep:      "0131010",
			expected: false,
		},
		{
			name:     "invalid CEP with more than 8 digits",
			cep:      "013101000",
			expected: false,
		},
		{
			name:     "invalid CEP with letters",
			cep:      "0131010a",
			expected: false,
		},
		{
			name:     "invalid CEP with dash",
			cep:      "01310-100",
			expected: false,
		},
		{
			name:     "invalid CEP empty string",
			cep:      "",
			expected: false,
		},
		{
			name:     "invalid CEP with special characters",
			cep:      "01310@00",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateCEP(tt.cep)
			if result != tt.expected {
				t.Errorf("ValidateCEP(%q) = %v, want %v", tt.cep, result, tt.expected)
			}
		})
	}
}

func TestViaCEPService_GetLocation(t *testing.T) {
	t.Run("successful CEP lookup", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"cep": "01310-100",
				"logradouro": "Avenida Paulista",
				"complemento": "até 610 - lado par",
				"bairro": "Bela Vista",
				"localidade": "São Paulo",
				"uf": "SP",
				"ibge": "3550308",
				"gia": "1004",
				"ddd": "11",
				"siafi": "7107"
			}`))
		}))
		defer server.Close()

		service := &ViaCEPService{
			baseURL:    server.URL,
			httpClient: server.Client(),
		}

		location, err := service.GetLocation(context.Background(), "01310100")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if location == nil {
			t.Fatal("expected location, got nil")
		}

		if location.Localidade != "São Paulo" {
			t.Errorf("expected localidade 'São Paulo', got '%s'", location.Localidade)
		}

		if location.UF != "SP" {
			t.Errorf("expected UF 'SP', got '%s'", location.UF)
		}
	})

	t.Run("CEP not found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"erro": true}`))
		}))
		defer server.Close()

		service := &ViaCEPService{
			baseURL:    server.URL,
			httpClient: server.Client(),
		}

		location, err := service.GetLocation(context.Background(), "99999999")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if location != nil {
			t.Error("expected nil location for not found CEP")
		}
	})

	t.Run("server error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		service := &ViaCEPService{
			baseURL:    server.URL,
			httpClient: server.Client(),
		}

		_, err := service.GetLocation(context.Background(), "01310100")
		if err == nil {
			t.Error("expected error for server error response")
		}
	})
}
