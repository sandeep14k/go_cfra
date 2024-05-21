package store

import (
	"context"
	"fmt"
	"log"

	"gin/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const uri = "mongodb://localhost:27017"

type MongoStore struct {
	Collection *mongo.Collection
}

func (m *MongoStore) OpenConnectionWithMongoDB() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Printf("Error occurred while establishing connection with mongodb: %v", err)
	}

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Printf("Error occurred while pinging the server: %v", err)
	}

	m.Collection = client.Database("cf_recent_Actions").Collection("recent_Actions")
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

}

func (m *MongoStore) StoreRecentActionsInTheDatabase(actions []models.RecentAction) error {
	var toInsertInterface []interface{}
	for _, action := range actions {
		toInsertInterface = append(toInsertInterface, action)
	}

	log.Printf("Trying to insert document to mongodb")

	_, err1 := m.Collection.InsertMany(context.TODO(), toInsertInterface)
	if err1 != nil {
		log.Printf("Error while inserting documents: %v", err1)
		return nil
	}
	log.Printf("Insertion successful")
	//return nil
	return nil
}

func (m *MongoStore) QueryRecentActions() ([]models.RecentAction, error) {
	cursor, err := m.Collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Printf("Error occurred while querying recentActions collection: %v", err)
	}

	var results []models.RecentAction

	err = cursor.All(context.TODO(), &results)
	if err != nil {
		log.Printf("Error occurred while iterating the cursor: %v", err)
	}

	return results, err

}

func (m *MongoStore) GetMaxTimeStamp() (int64, error) {
	groupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "MaxTimeStamp", Value: bson.D{{Key: "$max", Value: "$timeSeconds"}}},
		}},
	}
	//
	pipeline := mongo.Pipeline{groupStage}
	cursor, err := m.Collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		panic(err)
	}
	wrapper := struct {
		MaxTimeStamp int64
	}{}

	for cursor.Next(context.TODO()) {
		if err = cursor.Decode(&wrapper); err != nil {
			panic(err)
		}
	}
	return wrapper.MaxTimeStamp, err
}
