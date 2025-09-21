package clients

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"time"
	"goetl/internal/models"
)

func FetchAdsData() ([]models.AdPerformance, error) {
	adsURL := os.Getenv("ADS_API_URL")
	if adsURL == "" {
		return nil, errors.New("ADS_API_URL not set")
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(adsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch ads data: status " + resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var adsResp models.AdsAPIResponse
	if err := json.Unmarshal(body, &adsResp); err != nil {
		return nil, err
	}
	return adsResp.External.Ads.Performance, nil
}
