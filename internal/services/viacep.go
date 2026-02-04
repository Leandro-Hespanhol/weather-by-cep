package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/lhespanhol/weather-by-cep/internal/models"
)

// ViaCEPService handles CEP lookups
type ViaCEPService struct {
	baseURL    string
	httpClient *http.Client
}

// NewViaCEPService creates a new ViaCEP service
func NewViaCEPService() *ViaCEPService {
	return &ViaCEPService{
		baseURL: "https://viacep.com.br/ws",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// ValidateCEP validates if the CEP has the correct format (8 digits)
func ValidateCEP(cep string) bool {
	// CEP must be exactly 8 digits
	re := regexp.MustCompile(`^\d{8}$`)
	return re.MatchString(cep)
}

// GetLocation fetches the location for a given CEP
func (s *ViaCEPService) GetLocation(ctx context.Context, cep string) (*models.ViaCEPResponse, error) {
	url := fmt.Sprintf("%s/%s/json/", s.baseURL, cep)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch CEP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("viaCEP returned status %d", resp.StatusCode)
	}

	var viaCEPResp models.ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&viaCEPResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if viaCEPResp.Erro {
		return nil, nil // CEP not found
	}

	return &viaCEPResp, nil
}
