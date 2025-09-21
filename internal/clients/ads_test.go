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

// mockRoundTripper implements http.RoundTripper for mocking HTTP requests
type mockRoundTripper struct {
	response *http.Response
	err      error
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.response, m.err
}

func TestFetchAdsData_Success(t *testing.T) {
	// Prepare mock response
	body := json.RawMessage(`{"external": {"ads": {"performance": [{"campaign_id": "cmp1", "channel": "Google", "clicks": 100, "impressions": 1000, "cost": 123.45}]}}}`)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(string(body))),
	}

	// Patch http.DefaultClient's Transport
	origTransport := http.DefaultTransport
	http.DefaultTransport = &mockRoundTripper{response: resp, err: nil}
	defer func() { http.DefaultTransport = origTransport }()

	// Set env var
	os.Setenv("ADS_API_URL", "http://mock-ads-url")
	defer os.Unsetenv("ADS_API_URL")

	// Call function
	result, err := FetchAdsData()
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "cmp1", result[0].CampaignID)
	assert.Equal(t, "Google", result[0].Channel)
	assert.Equal(t, 100, result[0].Clicks)
	assert.Equal(t, 1000, result[0].Impressions)
	assert.Equal(t, 123.45, result[0].Cost)
}

func TestFetchAdsData_EnvNotSet(t *testing.T) {
	os.Unsetenv("ADS_API_URL")
	result, err := FetchAdsData()
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "ADS_API_URL not set")
}

func TestFetchAdsData_HTTPError(t *testing.T) {
	os.Setenv("ADS_API_URL", "http://mock-ads-url")
	defer os.Unsetenv("ADS_API_URL")
	origTransport := http.DefaultTransport
	http.DefaultTransport = &mockRoundTripper{response: nil, err: errors.New("network error")}
	defer func() { http.DefaultTransport = origTransport }()

	result, err := FetchAdsData()
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "network error")
}

func TestFetchAdsData_BadStatus(t *testing.T) {
	os.Setenv("ADS_API_URL", "http://mock-ads-url")
	defer os.Unsetenv("ADS_API_URL")
	resp := &http.Response{
		StatusCode: http.StatusInternalServerError,
		Status:     "500 Internal Server Error",
		Body:       io.NopCloser(strings.NewReader("")),
	}
	origTransport := http.DefaultTransport
	http.DefaultTransport = &mockRoundTripper{response: resp, err: nil}
	defer func() { http.DefaultTransport = origTransport }()

	result, err := FetchAdsData()
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to fetch ads data")
}

func TestFetchAdsData_BadJSON(t *testing.T) {
	os.Setenv("ADS_API_URL", "http://mock-ads-url")
	defer os.Unsetenv("ADS_API_URL")
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("not-json")),
	}
	origTransport := http.DefaultTransport
	http.DefaultTransport = &mockRoundTripper{response: resp, err: nil}
	defer func() { http.DefaultTransport = origTransport }()

	result, err := FetchAdsData()
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid character")
}
