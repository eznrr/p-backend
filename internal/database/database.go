package database

import (
    "context"
    "log"
    "time"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func ConnectDB() {
    var err error

    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

    defer cancel()
    Client, err = mongo.Connect(ctx, clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    err = Client.Ping(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Connected to MongoDB!")
}

func GetCollection(collectionName string) *mongo.Collection {
    return Client.Database("Culinario").Collection(collectionName)
}
