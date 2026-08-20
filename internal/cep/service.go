package cep

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var ErrZipcodeNotFound = errors.New("zipcode not found")

type connectorErrorResponse struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

type cepService struct{}

func NewCepService() *cepService {
	return &cepService{}
}

func (s *cepService) GetLocationByCep(cep string) (*CepLocation, error) {
	url := resolveConnectorUrl(cep)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrZipcodeNotFound
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("cep connector returned status %d", resp.StatusCode)
	}

	var connectorErr connectorErrorResponse
	if err := json.Unmarshal(body, &connectorErr); err == nil {
		if connectorErr.Name == "CepPromiseError" && connectorErr.Type == "service_error" {
			return nil, fmt.Errorf("%w: %s", ErrZipcodeNotFound, connectorErr.Message)
		}
	}

	var location CepLocation
	if err := json.Unmarshal(body, &location); err != nil {
		return nil, err
	}

	if location.CEP == "" {
		return nil, fmt.Errorf("invalid cep connector response")
	}

	return &location, nil
}

func resolveConnectorUrl(zipcode string) string {
	envUrl := os.Getenv("CEP_CONNECTOR_URL")
	if envUrl == "" {
		envUrl = "https://brasilapi.com.br/api/cep/v2/{cep}"
	}

	connectorBaseUrl := strings.ReplaceAll(envUrl, "{cep}", zipcode)

	return connectorBaseUrl
}
