package models

import (
	"encoding/json"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestETLResult_JSONMarshalling(t *testing.T) {
	in := ETLResult{
		Date:        "2025-09-20",
		Channel:     "Google",
		CampaignID:  "C1",
		UTMCampaign: "fall2025",
		Clicks:      10,
		Impressions: 100,
		Cost:        50.0,
		Leads:       2,
		Opportunities: 3,
		ClosedWon:   1,
		Revenue:     200.0,
		CPC:         5.0,
		CPA:         25.0,
		CVRLeadToOpp: 1.5,
		CVROppToWon:  0.33,
		ROAS:        4.0,
	}
	b, err := json.Marshal(in)
	assert.NoError(t, err)
	var out ETLResult
	assert.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in, out)
}

func TestCRMAPIResponse_JSONMarshalling(t *testing.T) {
	in := CRMAPIResponse{
		External: struct {
			CRM struct {
				Opportunities []Opportunity `json:"opportunities"`
			} `json:"crm"`
		}{
			CRM: struct {
				Opportunities []Opportunity `json:"opportunities"`
			}{
				Opportunities: []Opportunity{{OpportunityID: "O1", ContactEmail: "a@b.com", Stage: "lead", Amount: 100.0, CreatedAt: "2025-09-20", UTMCampaign: "fall2025", UTMSource: "google", UTMMedium: "cpc"}},
			},
		},
	}
	b, err := json.Marshal(in)
	assert.NoError(t, err)
	var out CRMAPIResponse
	assert.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in.External.CRM.Opportunities[0], out.External.CRM.Opportunities[0])
}

func TestAdsAPIResponse_JSONMarshalling(t *testing.T) {
	in := AdsAPIResponse{
		External: struct {
			Ads struct {
				Performance []AdPerformance `json:"performance"`
			} `json:"ads"`
		}{
			Ads: struct {
				Performance []AdPerformance `json:"performance"`
			}{
				Performance: []AdPerformance{{Date: "2025-09-20", CampaignID: "C1", Channel: "Google", Clicks: 10, Impressions: 100, Cost: 50.0, UTMCampaign: "fall2025", UTMSource: "google", UTMMedium: "cpc"}},
			},
		},
	}
	b, err := json.Marshal(in)
	assert.NoError(t, err)
	var out AdsAPIResponse
	assert.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in.External.Ads.Performance[0], out.External.Ads.Performance[0])
}

func TestOpportunity_Fields(t *testing.T) {
	opp := Opportunity{
		OpportunityID: "O1",
		ContactEmail:  "a@b.com",
		Stage:         "lead",
		Amount:        100.0,
		CreatedAt:     "2025-09-20",
		UTMCampaign:   "fall2025",
		UTMSource:     "google",
		UTMMedium:     "cpc",
	}
	assert.Equal(t, "O1", opp.OpportunityID)
	assert.Equal(t, "a@b.com", opp.ContactEmail)
	assert.Equal(t, "lead", opp.Stage)
	assert.Equal(t, 100.0, opp.Amount)
	assert.Equal(t, "2025-09-20", opp.CreatedAt)
	assert.Equal(t, "fall2025", opp.UTMCampaign)
	assert.Equal(t, "google", opp.UTMSource)
	assert.Equal(t, "cpc", opp.UTMMedium)
}

func TestAdPerformance_Fields(t *testing.T) {
	ad := AdPerformance{
		Date:        "2025-09-20",
		CampaignID:  "C1",
		Channel:     "Google",
		Clicks:      10,
		Impressions: 100,
		Cost:        50.0,
		UTMCampaign: "fall2025",
		UTMSource:   "google",
		UTMMedium:   "cpc",
	}
	assert.Equal(t, "2025-09-20", ad.Date)
	assert.Equal(t, "C1", ad.CampaignID)
	assert.Equal(t, "Google", ad.Channel)
	assert.Equal(t, 10, ad.Clicks)
	assert.Equal(t, 100, ad.Impressions)
	assert.Equal(t, 50.0, ad.Cost)
	assert.Equal(t, "fall2025", ad.UTMCampaign)
	assert.Equal(t, "google", ad.UTMSource)
	assert.Equal(t, "cpc", ad.UTMMedium)
}
