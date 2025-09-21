package db

import (
	"context"
	"log"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"goetl/internal/utils"
)


func getDBName() string {
	dbName := "goetldb"
	if env := getEnv("MONGO_DATABASE"); env != "" {
		dbName = env
	}
	return dbName
}


func getEnv(key string) string {
	if v := utils.Getenv(key); v != "" {
		return v
	}
	return ""
}


// GetCollection returns the MongoDB collection for a given name or nil on error
func GetCollection(collectionName string) (*mongo.Collection, context.Context, context.CancelFunc) {
	client, err := GetMongoClient()
	if err != nil {
		log.Printf("MongoDB connection error: %v", err)
		return nil, nil, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	collection := client.Database(getDBName()).Collection(collectionName)
	return collection, ctx, cancel
}
