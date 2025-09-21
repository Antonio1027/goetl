package clients

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestFetchCRMData_Success(t *testing.T) {
	body := json.RawMessage(`{"external":{"crm":{"opportunities":[{"opportunity_id":"O-9001","contact_email":"ana@example.com","stage":"closed_won","amount":5000.0,"created_at":"2025-08-05T10:22:00Z","utm_campaign":"back_to_school","utm_source":"google","utm_medium":"cpc"}]}}}`)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(string(body))),
	}
	origTransport := http.DefaultTransport
	http.DefaultTransport = &mockRoundTripper{response: resp, err: nil}
	defer func() { http.DefaultTransport = origTransport }()

	os.Setenv("CRM_API_URL", "http://mock-crm-url")
	defer os.Unsetenv("CRM_API_URL")

	result, err := FetchCRMData()
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "O-9001", result[0].OpportunityID)
	assert.Equal(t, "closed_won", result[0].Stage)
	assert.Equal(t, 5000.0, result[0].Amount)
}

func TestFetchCRMData_EnvNotSet(t *testing.T) {
	os.Unsetenv("CRM_API_URL")
	result, err := FetchCRMData()
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "CRM_API_URL not set")
}

func TestFetchCRMData_HTTPError(t *testing.T) {
	os.Setenv("CRM_API_URL", "http://mock-crm-url")
	defer os.Unsetenv("CRM_API_URL")
	origTransport := http.DefaultTransport
	http.DefaultTransport = &mockRoundTripper{response: nil, err: errors.New("network error")}
	defer func() { http.DefaultTransport = origTransport }()

	result, err := FetchCRMData()
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "network error")
}

func TestFetchCRMData_BadStatus(t *testing.T) {
	os.Setenv("CRM_API_URL", "http://mock-crm-url")
	defer os.Unsetenv("CRM_API_URL")
	resp := &http.Response{
		StatusCode: http.StatusInternalServerError,
		Status:     "500 Internal Server Error",
		Body:       io.NopCloser(strings.NewReader("")),
	}
	origTransport := http.DefaultTransport
	http.DefaultTransport = &mockRoundTripper{response: resp, err: nil}
	defer func() { http.DefaultTransport = origTransport }()

	result, err := FetchCRMData()
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to fetch CRM data")
}

func TestFetchCRMData_BadJSON(t *testing.T) {
	os.Setenv("CRM_API_URL", "http://mock-crm-url")
	defer os.Unsetenv("CRM_API_URL")
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("not-json")),
	}
	origTransport := http.DefaultTransport
	http.DefaultTransport = &mockRoundTripper{response: resp, err: nil}
	defer func() { http.DefaultTransport = origTransport }()

	result, err := FetchCRMData()
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid character")
}
