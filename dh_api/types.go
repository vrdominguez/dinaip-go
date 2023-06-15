/*
Package dh_api provides a client for interacting with the Dinahosting API.

The package includes types and functions for making requests to the Dinahosting API and handling the responses.
*/
package dh_api

import "net/http"

const (
	TypeA    RecordType = "A"
	TypeAAAA RecordType = "AAAA"
)

// RecordType represents the type of DNS record.
type RecordType string

// DinahostingClient represents a client for the Dinahosting API.
type DinahostingClient struct {
	httpClient   *http.Client
	baseURL      string
	username     string
	password     string
	responseType string
}

// DomainGetZonesResponse represents the response structure for the "Domain_Zone_GetTypeA" and "Domain_Zone_GetTypeAAAA" API commands.
type DomainGetZonesResponse struct {
	TrId         string        `json:"trId"`
	ResponseCode int           `json:"responseCode"`
	Errors       []interface{} `json:"errors"`
	Data         []struct {
		Ip       string `json:"ip"`
		Hostname string `json:"hostname"`
	} `json:"data"`
	Command string `json:"command"`
}

// GenericResponse represents a generic response structure for Dinahosting API commands.
type GenericResponse struct {
	TrId         string        `json:"trId"`
	ResponseCode int           `json:"responseCode"`
	Errors       []interface{} `json:"errors"`
	Message      string        `json:"message"`
}
