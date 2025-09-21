package db

import (
	"context"
	"fmt"
	"sync"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var (
	clientInstance *mongo.Client
	clientInstanceError error
	mongoOnce sync.Once
)

// GetMongoClient returns a singleton MongoDB client instance
func GetMongoClient() (*mongo.Client, error) {
	mongoOnce.Do(func() {
		uri := getMongoURI()
		log.Printf("Connecting to MongoDB at %s", uri)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		clientOptions := options.Client().ApplyURI(uri)
		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			clientInstanceError = err
			return
		}
		// Ping to verify connection
		if err = client.Ping(ctx, nil); err != nil {
			clientInstanceError = err
			return
		}
		clientInstance = client
	})
	return clientInstance, clientInstanceError
}

func getMongoURI() string {
	host := getEnv("MONGO_HOST")
	if host == "" {
		host = "mongodb"
	}
	port := getEnv("MONGO_PORT")
	if port == "" {
		port = "27017"
	}
	user := getEnv("MONGO_USERNAME")
	pass := getEnv("MONGO_PASSWORD")
	db := getEnv("MONGO_DATABASE")
	if user != "" && pass != "" {
		if db != "" {
			return fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", user, pass, host, port, db)
		}
		return fmt.Sprintf("mongodb://%s:%s@%s:%s", user, pass, host, port)
	}
	return fmt.Sprintf("mongodb://%s:%s", host, port)
}
