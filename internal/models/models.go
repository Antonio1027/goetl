package models
// ETLResult represents the consolidated data to persist after ETL processing
type ETLResult struct {
	Date           string  `json:"date"`
	Channel        string  `json:"channel"`
	CampaignID     string  `json:"campaign_id"`
	UTMCampaign    string  `json:"utm_campaign"`
	Clicks         int     `json:"clicks"`
	Impressions    int     `json:"impressions"`
	Cost           float64 `json:"cost"`
	Leads          int     `json:"leads"`
	Opportunities  int     `json:"opportunities"`
	ClosedWon      int     `json:"closed_won"`
	Revenue        float64 `json:"revenue"`
	CPC            float64 `json:"cpc"`
	CPA            float64 `json:"cpa"`
	CVRLeadToOpp   float64 `json:"cvr_lead_to_opp"`
	CVROppToWon    float64 `json:"cvr_opp_to_won"`
	ROAS           float64 `json:"roas"`
}

type CRMAPIResponse struct {
	External struct {
		CRM struct {
			Opportunities []Opportunity `json:"opportunities"`
		} `json:"crm"`
	} `json:"external"`
}

type Opportunity struct {
	OpportunityID string  `json:"opportunity_id"`
	ContactEmail  string  `json:"contact_email"`
	Stage         string  `json:"stage"`
	Amount        float64 `json:"amount"`
	CreatedAt     string  `json:"created_at"`
	UTMCampaign   string  `json:"utm_campaign"`
	UTMSource     string  `json:"utm_source"`
	UTMMedium     string  `json:"utm_medium"`
}

type AdsAPIResponse struct {
	External struct {
		Ads struct {
			Performance []AdPerformance `json:"performance"`
		} `json:"ads"`
	} `json:"external"`
}

type AdPerformance struct {
	Date        string  `json:"date"`
	CampaignID  string  `json:"campaign_id"`
	Channel     string  `json:"channel"`
	Clicks      int     `json:"clicks"`
	Impressions int     `json:"impressions"`
	Cost        float64 `json:"cost"`
	UTMCampaign string  `json:"utm_campaign"`
	UTMSource   string  `json:"utm_source"`
	UTMMedium   string  `json:"utm_medium"`
}
