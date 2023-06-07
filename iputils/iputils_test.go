package iputils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPublicIP(t *testing.T) {
	// Server responses correctly
	correctServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "127.0.0.1")
	}))
	defer correctServer.Close()

	// Failed response
	failedServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer failedServer.Close()

	// Override the services slice with the mock server URLs
	ipObtainServices = []string{correctServer.URL, failedServer.URL}

	t.Run("First server answers correctly", func(t *testing.T) {
		// Override the services slice with the mock server URLs
		ipObtainServices = []string{correctServer.URL, failedServer.URL}

		ip, err := GetPublicIp()

		assert.NoError(t, err)
		assert.Equal(t, "127.0.0.1", ip)
	})

	t.Run("First server fails but the second answers correctly", func(t *testing.T) {
		// Override the services slice with the mock server URLs
		ipObtainServices = []string{failedServer.URL, correctServer.URL}

		ip, err := GetPublicIp()

		assert.NoError(t, err)
		assert.Equal(t, "127.0.0.1", ip)
	})

	t.Run("No server can answer correctly", func(t *testing.T) {
		// Override the services slice with the mock server URLs
		ipObtainServices = []string{failedServer.URL}

		ip, err := GetPublicIp()

		assert.Error(t, err)
		assert.Equal(t, "", ip)
	})
}

func TestFqdnToIP(t *testing.T) {
	t.Run("Can solve example.com", func(t *testing.T) {
		ip, err := FqdnToIP("example.com")

		assert.NoError(t, err)
		assert.Equal(t, "93.184.216.34", ip)
	})

	t.Run("Cannot resolve an unexistent domain", func(t *testing.T) {
		ip, err := FqdnToIP("invented.domain.notld")

		assert.Error(t, err)
		assert.Equal(t, "", ip)
	})
}
