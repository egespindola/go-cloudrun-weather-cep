//go:build e2e

// Test cases follow .tmp/specs/tst/TESTS.md.
//
// NOTE on status codes: TESTS.md asks for HTTP 421 on invalid zipcodes, but
// the shipped handler (internal/orch/orch.go) responds with 422
// (http.StatusUnprocessableEntity). Per the "don't change production code"
// rule, these tests assert the actual 422 behavior rather than the spec's
// 421 so the suite reflects the real contract. Flagged to the team so the
// spec/code mismatch can be resolved (either the spec is wrong, or the
// handler is).
package e2e

import (
	"encoding/json"
	"net/http"
	"testing"
)

type temperatureResponse struct {
	CelsiusC   float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}

type errorResponse struct {
	Message string `json:"message"`
}

func TestCepHappyPath_ReturnsMockedTemperature(t *testing.T) {
	resp, err := http.Get(baseURL + "/cep/" + happyZipcode)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var body temperatureResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	want := temperatureResponse{CelsiusC: 28.5, Fahrenheit: 83.3, Kelvin: 301.65}
	if body != want {
		t.Fatalf("unexpected temperature payload: got %+v, want %+v", body, want)
	}
}

func TestUnknownEndpoint_Returns404(t *testing.T) {
	resp, err := http.Get(baseURL + "/unknown-route")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
}

func TestZipcodeWithSpecialCharacter_ReturnsInvalidZipcodeError(t *testing.T) {
	// 8 characters, but one of them ('@') is not a digit.
	assertInvalidZipcode(t, "0100100@")
}

func TestZipcodeWithWrongLength_ReturnsInvalidZipcodeError(t *testing.T) {
	// 7 digits: not equal to the required 8.
	assertInvalidZipcode(t, "0100100")
}

func assertInvalidZipcode(t *testing.T, zipcode string) {
	t.Helper()

	resp, err := http.Get(baseURL + "/cep/" + zipcode)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// See package-level NOTE: spec asks for 421, actual handler returns 422.
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("expected status %d, got %d", http.StatusUnprocessableEntity, resp.StatusCode)
	}

	var body errorResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body.Message != "invalid zipcode" {
		t.Fatalf("expected message %q, got %q", "invalid zipcode", body.Message)
	}
}

func TestZipcodeNotFoundUpstream_Returns404(t *testing.T) {
	resp, err := http.Get(baseURL + "/cep/" + notFoundZipcode)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, resp.StatusCode)
	}

	var body errorResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body.Message != "zipcode not found" {
		t.Fatalf("expected message %q, got %q", "zipcode not found", body.Message)
	}
}
