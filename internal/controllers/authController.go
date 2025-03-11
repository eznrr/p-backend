package controllers

import (
    "context"
    "time"
    "net/http"

    "gitlab.com/ltp2-b-crepusculo/ps-backend-enzo-rodrigues/internal/models/users"
    "gitlab.com/ltp2-b-crepusculo/ps-backend-enzo-rodrigues/internal/utils"
    "gitlab.com/ltp2-b-crepusculo/ps-backend-enzo-rodrigues/internal/database"
    
    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "github.com/go-playground/validator/v10"
    
)

type Credentials struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type AuthController struct {
    userCollection *mongo.Collection
    validate       *validator.Validate
}

func NewAuthController() *AuthController {
    userCollection := database.GetCollection("users")
    return &AuthController{userCollection, validator.New()}
}

func (ac *AuthController) Register(c *gin.Context) {
    var user models.User

    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "400", "message": "Bad Request", "reason": err.Error()})
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var existingUser models.User
	err := ac.userCollection.FindOne(ctx, bson.M{
		"$or": []bson.M{
			{"username": user.Username},
			{"email": user.Email},
		},
	}).Decode(&existingUser)
    
    if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "409", "message": "Conflict", "reason": "User already exists",})
		return
	}

    user.ID = utils.GenerateStringID()
    user.Password = utils.HashPassword(user.Password)
    user.CreatedAt = time.Now()

    _, err = ac.userCollection.InsertOne(context.Background(), user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "500", "message": "Internal Server Error", "reason": "Failed to create user"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Successfully registered user", "user": user})
}

func (ac *AuthController) Login(c *gin.Context) {
    var creds Credentials

    if err := c.ShouldBindJSON(&creds); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "400", "message": "Bad Request", "reason": err.Error()})
        return
    }

    var user models.User

    err := ac.userCollection.FindOne(context.Background(), bson.M{"email": creds.Email}).Decode(&user)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "401", "message": "Unauthorized", "reason": "Invalid email"})
        return
    }

    if !utils.CheckPasswordHash(creds.Password, user.Password) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "401", "message": "Unauthorized", "reason": "Invalid Password"})
        return
    }

    token, err := utils.GenerateJWT(creds.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "500", "message": "Internal Server Error", "reason": "Failed to generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Logged in successfully", "token": token})
}
