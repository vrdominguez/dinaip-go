package dh_api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTypeIpRecords(t *testing.T) {
	username := "testuser"
	password := "testpassword"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"trId":"","responseCode":1000,"errors":[],"data":[{"ip":"192.0.2.1","hostname":"example"}],"command":""}`))
	}))
	defer ts.Close()

	client := &DinahostingClient{
		httpClient:   http.DefaultClient,
		baseURL:      ts.URL,
		username:     username,
		password:     password,
		responseType: "json",
	}

	t.Run("Can get type A or AAAA records", func(t *testing.T) {
		domain := "example.com"
		response, err := client.GetTypeIpRecords(domain, TypeA)

		assert.NoError(t, err)
		assert.Equal(t, 1000, response.ResponseCode)
		assert.Len(t, response.Data, 1)

		record := response.Data[0]
		assert.Equal(t, "example", record.Hostname)
		assert.Equal(t, "192.0.2.1", record.Ip)
	})

	t.Run("Fals with invalid record", func(t *testing.T) {
		domain := "example.com"
		_, err := client.GetTypeIpRecords(domain, "CNAME")
		assert.Error(t, err)

	})
}

func TestSetTypeIpRecord(t *testing.T) {
	username := "testuser"
	password := "testpassword"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"trId":"","responseCode":1000,"errors":[],"data":null,"command":""}`))
	}))
	defer ts.Close()

	client := &DinahostingClient{
		httpClient:   http.DefaultClient,
		baseURL:      ts.URL,
		username:     username,
		password:     password,
		responseType: "json",
	}

	t.Run("Can set type A record", func(t *testing.T) {
		domain := "example.com"
		hostname := "example"
		ip := "192.0.2.1"
		oldIp := ""
		err := client.SetTypeIpRecord(domain, hostname, ip, oldIp, TypeA)

		assert.NoError(t, err)
	})

	t.Run("Can update type A record", func(t *testing.T) {
		domain := "example.com"
		hostname := "example"
		ip := "192.0.2.1"
		oldIp := "192.0.2.2"
		err := client.SetTypeIpRecord(domain, hostname, ip, oldIp, TypeA)

		assert.NoError(t, err)
	})

	t.Run("Can set type AAAA record", func(t *testing.T) {
		domain := "example.com"
		hostname := "example"
		ip := "2001:db8::1"
		oldIp := ""
		err := client.SetTypeIpRecord(domain, hostname, ip, oldIp, TypeAAAA)

		assert.NoError(t, err)
	})

	t.Run("Can update type AAAA record", func(t *testing.T) {
		domain := "example.com"
		hostname := "example"
		ip := "2001:db8::1"
		oldIp := "2001:db8::2"
		err := client.SetTypeIpRecord(domain, hostname, ip, oldIp, TypeAAAA)

		assert.NoError(t, err)
	})

	t.Run("Fails with invalid record type", func(t *testing.T) {
		domain := "example.com"
		hostname := "example"
		ip := "192.0.2.1"
		oldIp := ""
		err := client.SetTypeIpRecord(domain, hostname, ip, oldIp, "CNAME")

		assert.Error(t, err)
	})
}
