package controllers

import (
    "context"
    "net/http"
    "io/ioutil"
    "encoding/base64"
    
    "gitlab.com/ltp2-b-crepusculo/ps-backend-enzo-rodrigues/internal/models/users"
    "gitlab.com/ltp2-b-crepusculo/ps-backend-enzo-rodrigues/internal/database"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "github.com/gin-gonic/gin"
)

type UserController struct {
    userCollection *mongo.Collection
}

func NewUserController() *UserController {
    userCollection := database.GetCollection("users")
    return &UserController{userCollection}
}

func (uc *UserController) GetProfile(c *gin.Context) {
    email := c.GetString("email")
	if email == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "401", "message": "Unauthorized", "reason": "User not logged in"})
		return
	}

	var user models.User
	err := uc.userCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "500", "message": "Internal Server Error", "reason": "Failed to find user"})
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, user)
}

func (uc *UserController) GetProfileByUsername(c *gin.Context) {
	username := c.Param("username")

    var user models.User
    err := uc.userCollection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "404", "message": "Not Found", "reason": "Profile not found"})
        return
    }

    c.JSON(http.StatusOK, user)
}

func (uc *UserController) UpdateProfileImage(c *gin.Context) {
    email := c.GetString("email")
    if email == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "401", "message": "Unauthorized", "reason": "User not logged in"})
        return
    }

    file, _, err := c.Request.FormFile("profile_image")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "400", "message": "Bad Request", "reason": "Invalid image data"})
        return
    }
    defer file.Close()

    imageData, err := ioutil.ReadAll(file)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "500", "message": "Internal Servel Error", "reason": "Error reading image data"})
        return
    }

    encodedImage := base64.StdEncoding.EncodeToString(imageData)

    filter := bson.M{"email": email}
    update := bson.M{"$set": bson.M{"profile_image": encodedImage}}

    _, err = uc.userCollection.UpdateOne(context.Background(), filter, update)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "500", "message": "Internal Server Error", "reason": "Error updating profile image"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Profile image updated successfully"})
}