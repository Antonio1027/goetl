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

func FetchCRMData() ([]models.Opportunity, error) {
	crmURL := os.Getenv("CRM_API_URL")
	if crmURL == "" {
		return nil, errors.New("CRM_API_URL not set")
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(crmURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch CRM data: status " + resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var crmResp models.CRMAPIResponse
	if err := json.Unmarshal(body, &crmResp); err != nil {
		return nil, err
	}
	return crmResp.External.CRM.Opportunities, nil
}
