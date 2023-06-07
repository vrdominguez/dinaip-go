package dh_api

import "net/http"

const (
	TypeA    RecordType = "A"
	TypeAAAA RecordType = "AAAA"
)

type RecordType string

type DinahostingClient struct {
	httpClient   *http.Client
	baseURL      string
	username     string
	password     string
	responseType string
}

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

// {"trId":"dh647f659234c059.72531182","responseCode":1000,"message":"Success.","data":null,"command":"Domain_Zone_AddTypeA"}
type GenericResponse struct {
	TrId         string        `json:"trId"`
	ResponseCode int           `json:"responseCode"`
	Errors       []interface{} `json:"errors"`
	Message      string        `json:"message"`
}
