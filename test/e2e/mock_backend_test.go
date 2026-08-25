//go:build e2e

package e2e

import (
	"encoding/json"
	"net/http"
	"strings"
)

// newMockBackendHandler stands in for the real BrasilAPI CEP and Open-Meteo
// weather upstreams so the containerized app under test gets deterministic,
// offline responses instead of hitting the live internet.
func newMockBackendHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/cep/", handleMockCep)
	mux.HandleFunc("/weather", handleMockWeather)
	return mux
}

func handleMockCep(w http.ResponseWriter, r *http.Request) {
	cep := strings.TrimPrefix(r.URL.Path, "/cep/")

	if cep == notFoundZipcode {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	location := map[string]any{
		"cep":          happyZipcode,
		"state":        "SP",
		"city":         "Sao Paulo",
		"neighborhood": "Se",
		"street":       "Praca da Se",
		"service":      "mock",
		"location": map[string]any{
			"coordinates": map[string]string{
				"longitude": "-46.633309",
				"latitude":  "-23.550520",
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(location)
}

func handleMockWeather(w http.ResponseWriter, r *http.Request) {
	weather := map[string]any{
		"current_weather": map[string]any{
			"temperature": mockTempC,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(weather)
}
