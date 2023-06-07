package dh_api

import (
	"net/http"
	"net/url"
)

const (
	defaultBaseURL      = "https://dinahosting.com/special/api.php"
	defaultResponseType = "json"
)

func NewClient(username, password string) *DinahostingClient {
	return &DinahostingClient{
		httpClient:   http.DefaultClient,
		baseURL:      defaultBaseURL,
		username:     username,
		password:     password,
		responseType: defaultResponseType,
	}
}

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
