package utils

import (
    "context"

    "gitlab.com/ltp2-b-crepusculo/ps-backend-enzo-rodrigues/internal/database"
    "gitlab.com/ltp2-b-crepusculo/ps-backend-enzo-rodrigues/internal/models/users"

    "go.mongodb.org/mongo-driver/bson"
)

func GetUsernameByEmail(email string) (string, error) {
    userCollection := database.GetCollection("users")

    var user models.User
    err := userCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
    if err != nil {
        return "", err
    }

    return user.Username, nil
}
