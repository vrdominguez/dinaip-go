package dh_api

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateInstance(t *testing.T) {
	t.Run("Can create an instance with new", func(t *testing.T) {
		client := NewClient("myUser", "myPassword")

		assert.IsType(t, &DinahostingClient{}, client)

		assert.Equal(t, "myUser", client.username)
		assert.Equal(t, "myPassword", client.password)
		assert.Equal(t, defaultBaseURL, client.baseURL)
		assert.Equal(t, defaultResponseType, client.responseType)
	})
}

func TestDoRequest(t *testing.T) {
	// Create a test server to mock the Dinahosting API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"response": "ok"}`))
	}))
	defer server.Close()

	t.Run("Can do a httop request to dinahosting", func(t *testing.T) {
		// Create a new client with the test server URL
		client := &DinahostingClient{
			httpClient: http.DefaultClient,
			baseURL:    server.URL,
			username:   "testuser",
			password:   "testpass",
		}

		// Test the doRequest method
		params := url.Values{}
		params.Set("command", "test")
		resp, err := client.DoRequest(params)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
