//go:build e2e

// Package e2e runs black-box tests against the real application Docker
// image described in TESTS.md. It never imports production packages: the
// only contract with the app is the HTTP surface exposed by the container.
//
// The CEP and weather upstreams are swapped for a local mock HTTP server via
// CEP_CONNECTOR_URL/WEATHER_CONNECTOR_URL env overrides passed to `docker
// run`, so the container never touches the real internet and always
// produces the deterministic mock reading required by the spec.
package e2e

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

const (
	imageName     = "weather-cep:e2e-test"
	containerName = "weather-cep-e2e-test"

	happyZipcode    = "01001000"
	notFoundZipcode = "01001123"
	mockTempC       = 28.5
)

var baseURL string

func TestMain(m *testing.M) {
	os.Exit(run(m))
}

func run(m *testing.M) int {
	mockListener, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		log.Printf("failed to start mock backend listener: %v", err)
		return 1
	}
	mockPort := mockListener.Addr().(*net.TCPAddr).Port

	mockServer := &http.Server{Handler: newMockBackendHandler()}
	go func() {
		_ = mockServer.Serve(mockListener)
	}()
	defer mockServer.Close()

	// Clean up any container left over from a previous interrupted run.
	dockerCleanup()

	if err := dockerBuild(); err != nil {
		log.Printf("docker build failed: %v", err)
		return 1
	}
	defer dockerCleanup()

	hostPort, err := dockerRun(mockPort)
	if err != nil {
		log.Printf("docker run failed: %v", err)
		return 1
	}

	baseURL = fmt.Sprintf("http://127.0.0.1:%d", hostPort)

	if err := waitForServer(baseURL, 20*time.Second); err != nil {
		log.Printf("app container did not become ready: %v", err)
		printContainerLogs()
		return 1
	}

	return m.Run()
}

// dockerBuild builds the production image from the project root. go test
// sets the working directory to this package's directory, so the project
// root is always two levels up.
func dockerBuild() error {
	cmd := exec.Command("docker", "build", "-t", imageName, ".")
	cmd.Dir = "../.."
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, out)
	}
	return nil
}

func dockerRun(mockPort int) (int, error) {
	cepURL := fmt.Sprintf("http://host.docker.internal:%d/cep/{cep}", mockPort)
	weatherURL := fmt.Sprintf("http://host.docker.internal:%d/weather?latitude={latitude}&longitude={longitude}", mockPort)

	args := []string{
		"run", "-d",
		"--name", containerName,
		"--add-host=host.docker.internal:host-gateway",
		"-e", "PORT=8080",
		"-e", "CEP_CONNECTOR_URL=" + cepURL,
		"-e", "WEATHER_CONNECTOR_URL=" + weatherURL,
		"-p", "127.0.0.1::8080",
		imageName,
	}

	if out, err := exec.Command("docker", args...).CombinedOutput(); err != nil {
		return 0, fmt.Errorf("%w: %s", err, out)
	}

	portOut, err := exec.Command("docker", "port", containerName, "8080/tcp").CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("docker port failed: %w: %s", err, portOut)
	}

	return parseHostPort(string(portOut))
}

func parseHostPort(output string) (int, error) {
	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		idx := strings.LastIndex(line, ":")
		if idx == -1 {
			continue
		}
		var port int
		if _, err := fmt.Sscanf(line[idx+1:], "%d", &port); err == nil {
			return port, nil
		}
	}
	return 0, fmt.Errorf("could not parse docker port output: %q", output)
}

func waitForServer(base string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	client := &http.Client{Timeout: 2 * time.Second}

	var lastErr error
	for time.Now().Before(deadline) {
		resp, err := client.Get(base + "/cep/" + happyZipcode)
		if err == nil {
			resp.Body.Close()
			return nil
		}
		lastErr = err
		time.Sleep(300 * time.Millisecond)
	}
	return fmt.Errorf("timed out waiting for server: %w", lastErr)
}

func dockerCleanup() {
	_ = exec.Command("docker", "rm", "-f", containerName).Run()
}

func printContainerLogs() {
	out, _ := exec.Command("docker", "logs", containerName).CombinedOutput()
	log.Printf("container logs:\n%s", out)
}
