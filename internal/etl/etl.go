package etl

import (
	"goetl/internal/models"
	"goetl/internal/utils"
	"goetl/internal/clients"
	"strings"
	"log"
	"time"
)


// RunETL orchestrates the ETL process: Extract, Transform, Load. It accepts an optional 'since' parameter to filter data.
func RunETL(since string) ([]models.ETLResult, error) {
       ads, crm, err := Extract()
       if err != nil {
	       return nil, err
       }
       results, err := Transform(ads, crm, since)

       if err != nil {
            return nil, err
       }
	   if len(results) == 0 {
		   log.Println("No ETL results to load")
		   return results, nil
	   }
       Load(results)
       // code to wait the persistence of results
       time.Sleep(1 * time.Second)
       return results, nil
}


// Extract fetches data from third-party services and returns raw ads and crm data.
func Extract() ([]models.AdPerformance, []models.Opportunity, error) {
	ads, err := clients.GetAdsData()
	if err != nil {
		return nil, nil, err
	}
	crm, err := clients.GetOpportunities()
	if err != nil {
		return nil, nil, err
	}
	return ads, crm, nil
}


// Transformation of data, filters by 'since', deduplicates, normalizes, crosses, calculates metrics, and persists results.
func Transform(ads []models.AdPerformance, opportunities []models.Opportunity, since string) ([]models.ETLResult, error) {
	opportunities = TransformOpportunitiesData(opportunities)
	ads = TransformPerformanceData(ads)

	// Filter by 'since' date if provided
	adsFiltered := make([]models.AdPerformance, 0)
	for _, ad := range ads {
		if since == "" || ad.Date >= since {
			adsFiltered = append(adsFiltered, ad)
		}
	}
	opportunitiesFiltered := make([]models.Opportunity, 0)
	for _, opp := range opportunities {
		if since == "" || strings.HasPrefix(opp.CreatedAt, since) || opp.CreatedAt >= since {
			opportunitiesFiltered = append(opportunitiesFiltered, opp)
		}
	}

	// Deduplicate ads by (date, channel, campaign_id)
	adsMap := make(map[string]models.AdPerformance)
	for _, ad := range adsFiltered {
		key := ad.Date + ":" + ad.Channel + ":" + ad.CampaignID
		adsMap[key] = ad
	}

	// Deduplicate CRM by (created_at, utm_campaign, utm_source, utm_medium)
	opportunitiesMap := make(map[string]models.Opportunity)
	for _, opp := range opportunitiesFiltered {
		key := opp.CreatedAt + ":" + opp.UTMCampaign + ":" + opp.UTMSource + ":" + opp.UTMMedium
		opportunitiesMap[key] = opp
	}

	// Cross ads and CRM by utm_campaign, utm_source, utm_medium
	results := make([]models.ETLResult, 0)
	for _, ad := range adsMap {
		var leads, opportunities, closedWon int
		var revenue float64
		for _, opp := range opportunitiesMap {
			if ad.Date == opp.CreatedAt && ad.UTMCampaign == opp.UTMCampaign && ad.UTMSource == opp.UTMSource && ad.UTMMedium == opp.UTMMedium {
				opportunities++
				if opp.Stage == "lead" {
					leads++
				}
				if opp.Stage == "closed_won" {
					closedWon++
					revenue += opp.Amount
				}
			}
		}

		// Calculate metrics
		cpc := 0.0
		if ad.Clicks > 0 {
			cpc = ad.Cost / float64(ad.Clicks)
		}
		cpa := 0.0
		if leads > 0 {
			cpa = ad.Cost / float64(leads)
		}
		cvrLeadToOpp := 0.0
		if leads > 0 {
			cvrLeadToOpp = float64(opportunities) / float64(leads)
		}
		cvrOppToWon := 0.0
		if opportunities > 0 {
			cvrOppToWon = float64(closedWon) / float64(opportunities)
		}
		roas := 0.0
		if ad.Cost > 0 {
			roas = revenue / ad.Cost
		}

		res := models.ETLResult{
			Date:          ad.Date,
			Channel:       ad.Channel,
			CampaignID:    ad.CampaignID,
			UTMCampaign:   ad.UTMCampaign,
			Clicks:        ad.Clicks,
			Impressions:   ad.Impressions,
			Cost:          ad.Cost,
			Leads:         leads,
			Opportunities: opportunities,
			ClosedWon:     closedWon,
			Revenue:       revenue,
			CPC:           utils.RoundFloat(cpc, 2),
			CPA:           utils.RoundFloat(cpa, 2),
			CVRLeadToOpp:  utils.RoundFloat(cvrLeadToOpp, 2),
			CVROppToWon:   utils.RoundFloat(cvrOppToWon, 2),
			ROAS:          utils.RoundFloat(roas, 2),
		}
		results = append(results, res)
	}

	return results, nil
}


// Load persists the ETL results into MongoDB.
func Load(data []models.ETLResult) {
       for _, result := range data {
	       SaveResult(result)
       }
}


// transforms and normalizes ads performance data
func TransformPerformanceData(data []models.AdPerformance) []models.AdPerformance {
	var err error
	transformed := make([]models.AdPerformance, 0, len(data))
	for _, ad := range data {
		if ad.Date != "" {
			ad.Date, err = utils.NormalizeDate(ad.Date)
			if err != nil {
				ad.Date = ""
			}
		}
		if ad.Date != "" && ad.Channel != "" && ad.CampaignID != "" {
			ad.Channel = utils.SanitizeString(ad.Channel)
			ad.CampaignID = utils.SanitizeString(ad.CampaignID)
			ad.Clicks = utils.SanitizeInt(ad.Clicks)
			ad.Impressions = utils.SanitizeInt(ad.Impressions)
			ad.Cost = utils.SanitizeFloat(ad.Cost)
			transformed = append(transformed, ad)
		}
	}
	return transformed
}


// transforms and normalizes CRM opportunities data
func TransformOpportunitiesData(data []models.Opportunity) []models.Opportunity {
	var err error
	transformed := make([]models.Opportunity, 0, len(data))
	for _, opp := range data {
		if opp.CreatedAt != "" {
			opp.CreatedAt, err = utils.NormalizeDate(opp.CreatedAt)
			if err != nil {
				opp.CreatedAt = ""
			}
		}
		if opp.CreatedAt != "" && opp.UTMCampaign != "" && opp.UTMSource != "" && opp.UTMMedium != "" {
			opp.ContactEmail = utils.SanitizeString(opp.ContactEmail)
			opp.UTMCampaign = utils.SanitizeString(opp.UTMCampaign)
			opp.UTMSource = utils.SanitizeString(opp.UTMSource)
			opp.UTMMedium = utils.SanitizeString(opp.UTMMedium)
			opp.Stage = utils.SanitizeString(opp.Stage)
			opp.Amount = utils.SanitizeFloat(opp.Amount)
			opp.OpportunityID = utils.SanitizeString(opp.OpportunityID)
			transformed = append(transformed, opp)
		}
	}
	return transformed
}