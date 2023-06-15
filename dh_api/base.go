/*
Package dh_api provides a client for interacting with the Dinahosting API.

The package includes types and functions for making requests to the Dinahosting API and handling the responses.
*/

package dh_api

import (
	"net/http"
	"net/url"
)

const (
	defaultBaseURL      = "https://dinahosting.com/special/api.php"
	defaultResponseType = "json"
)

// NewClient creates a new DinahostingClient with the specified username and password.
// It returns a pointer to the DinahostingClient.
func NewClient(username, password string) *DinahostingClient {
	return &DinahostingClient{
		httpClient:   http.DefaultClient,
		baseURL:      defaultBaseURL,
		username:     username,
		password:     password,
		responseType: defaultResponseType,
	}
}

/*
DoRequest performs a request to the Dinahosting API with the given parameters.

It takes the API parameters as a url.Values and returns an http.Response and an error if any.

Example usage:

	// Create parameters for the API request
	params := url.Values{}
	params.Add("param1", "value1")
	params.Add("param2", "value2")

	// Perform the API request
	response, err := client.DoRequest(params)
	if err != nil {
		// Handle error
	}

	// Access the response data...
*/
func (c *DinahostingClient) DoRequest(params url.Values) (*http.Response, error) {
	params.Set("AUTH_USER", c.username)
	params.Set("AUTH_PWD", c.password)
	params.Set("responseType", c.responseType)

	req, err := http.NewRequest("POST", c.baseURL, nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = params.Encode()
	return c.httpClient.Do(req)
}
