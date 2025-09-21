package api

import (
	"github.com/gin-gonic/gin"
	"goetl/internal/etl"
	"goetl/internal/models"
	"time"
	"net/http"
	"fmt"
	"goetl/internal/utils"
)

func RegisterRoutes(r *gin.Engine) {
	r.POST("/ingest/run", ingestRunHandler)
	r.GET("/metrics/channel", metricsByChannelHandler)
	r.GET("/metrics/campaign", metricsByCampaignHandler)
}

// metricsByCampaignHandler handles GET /metrics/campaign?from=YYYY-MM-DD&to=YYYY-MM-DD&utm_campaign=google_ads&limit=10&offset=0
func metricsByCampaignHandler(c *gin.Context) {
       from := c.Query("from")
       to := c.Query("to")
       utmCampaign := c.Query("utm_campaign")
       limit := utils.ParseQueryInt(c, "limit", 10)
       offset := utils.ParseQueryInt(c, "offset", 0)

       // Fetch results filtered by utm_campaign and date range
       results := etl.GetResultsByCampaign(from, to, utmCampaign, limit, offset)

       c.JSON(http.StatusOK, gin.H{
	       "total": len(results),
	       "limit": limit,
	       "offset": offset,
	       "results": results,
       })
}

// metricsByChannelHandler handles GET /metrics/channel?from=YYYY-MM-DD&to=YYYY-MM-DD&channel=google_ads&pageSize=10&page=0
func metricsByChannelHandler(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	channel := c.Query("channel")
	limit := utils.ParseQueryInt(c, "limit", 10)
	offset := utils.ParseQueryInt(c, "offset", 0)

	results := etl.GetResultsByChannel(from, to, channel, limit, offset)

	c.JSON(http.StatusOK, gin.H{
		"total":  len(results),
		"limit":  limit,
		"offset": offset,
		"results": results,
	})
}


// ingestRunHandler handles POST /ingest/run?since=YYYY-MM-DD
func ingestRunHandler(c *gin.Context) {
	since := c.Query("since")
	if since == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing 'since' query parameter"})
		return
	}

	// Retry logic with backoff (simple, 3 attempts)
	var lastErr error
	var results []models.ETLResult
	var err error
	for attempt := 1; attempt <= 3; attempt++ {
		results, err = etl.RunETL(since)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"message": fmt.Sprintf("ETL process completed successfully. Processed %d records.", len(results)),
				"results": results,
			})
			return
		}
		lastErr = err
		// If the client canceled the request, stop retrying
		if c.Request.Context().Err() != nil {
			break
		}
		time.Sleep(time.Duration(attempt) * 500 * time.Millisecond)
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"error": lastErr.Error(),
	})
}
