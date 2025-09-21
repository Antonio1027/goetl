package clients

import (
	"goetl/internal/models"
)

// GetOpportunities fetches CRM opportunities data using the clients package and transforms them
func GetOpportunities() ([]models.Opportunity, error) {
	opportunities, err := FetchCRMData()
	if err != nil {
		return nil, err
	}
	return opportunities, nil
}

// FetchAdsData fetches ads data using the clients package and transforms them
func GetAdsData() ([]models.AdPerformance, error) {
	ads, err := FetchAdsData()
	if err != nil {
		return nil, err
	}
	return ads, nil
}
