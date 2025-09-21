package etl

import (
	"log"
	"goetl/internal/models"
	"goetl/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const etlCollection = "etl_results"


// SaveResult persists an ETLResult in MongoDB
func SaveResult(res models.ETLResult) {
	collection, ctx, cancel := db.GetCollection(etlCollection)
	if collection == nil || ctx == nil || cancel == nil {
		return
	}
	if cancel != nil {
		defer cancel()
	}
	filter := bson.M{"date": res.Date, "channel": res.Channel, "campaignid": res.CampaignID}
	update := bson.M{"$set": res}
	_, err := collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		log.Printf("Failed to upsert ETLResult: %v", err)
	}
}


// GetResultsByCampaign returns ETL results filtered by utm_campaign, date range, and paginated
func GetResultsByCampaign(dateStart, dateEnd, utmCampaign string, limit int, offset int) []models.ETLResult {
       return getResults(dateStart, dateEnd, map[string]string{"utmcampaign": utmCampaign}, limit, offset)
}


// GetResultsByChannel returns ETL results filtered by channel, date range, and paginated
func GetResultsByChannel(dateStart, dateEnd, channel string, limit, offset int) []models.ETLResult {
	return getResults(dateStart, dateEnd, map[string]string{"channel": channel}, limit, offset)
}


// getResults is a shared helper for channel/campaign queries with pagination by date range and limit/offset
func getResults(dateStart, dateEnd string, filterFields map[string]string, limit, offset int) []models.ETLResult {
       collection, ctx, cancel := db.GetCollection(etlCollection)
       if collection == nil || ctx == nil || cancel == nil {
	       return nil
       }
       defer cancel()
       filter := bson.M{}
       for k, v := range filterFields {
	       if v != "" {
		       filter[k] = v
	       }
       }
       // filter collection by date range
       if dateStart != "" && dateEnd != "" {
	       filter["date"] = bson.M{"$gte": dateStart, "$lte": dateEnd}
       } else if dateStart != "" {
	       filter["date"] = bson.M{"$gte": dateStart}
       } else if dateEnd != "" {
	       filter["date"] = bson.M{"$lte": dateEnd}
       }

       opts := options.Find()
       if limit > 0 {
	       opts.SetLimit(int64(limit))
       }
       if offset > 0 {
	       opts.SetSkip(int64(offset))
       }

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		log.Printf("Failed to fetch filtered ETLResults: %v", err)
		return nil
       }
       defer cursor.Close(ctx)

       var results []models.ETLResult
       for cursor.Next(ctx) {
	       var res models.ETLResult
	       if err := cursor.Decode(&res); err != nil {
			log.Printf("Failed to decode ETLResult: %v", err)
		       continue
	       }
	       results = append(results, res)
       }
       if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
       }
       return results
}
