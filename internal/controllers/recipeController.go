package controllers

import (
    "time"
    "context"
    "encoding/json"
    "net/http"
    "encoding/base64"

    "gitlab.com/ltp2-b-crepusculo/ps-backend-enzo-rodrigues/internal/models/recipes"
    "gitlab.com/ltp2-b-crepusculo/ps-backend-enzo-rodrigues/internal/database"
    "gitlab.com/ltp2-b-crepusculo/ps-backend-enzo-rodrigues/internal/utils"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
)

type RecipeController struct {
    recipeCollection *mongo.Collection
    validate         *validator.Validate
}

func NewRecipeController() *RecipeController {
    recipeCollection := database.GetCollection("recipes")
    return &RecipeController{recipeCollection, validator.New()}
}

func (rc *RecipeController) GetAllRecipes(c *gin.Context) {
    cursor, err := rc.recipeCollection.Find(context.TODO(), bson.D{})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "500", "message": "Internal Server Error", "reason": "Failed to get recipes"})
        return
    }
    defer cursor.Close(context.TODO())

    var recipes []models.Recipe
    for cursor.Next(context.TODO()) {
        var recipe models.Recipe
        cursor.Decode(&recipe)

        recipes = append(recipes, recipe)
    }

    c.JSON(http.StatusOK, recipes)
}

func (rc *RecipeController) GetRecipeByID(c *gin.Context) {
    id := c.Param("id")
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "400", "message": "Bad Request", "reason": "Invalid recipe ID"})
        return
    }

    var recipe models.Recipe
    err = rc.recipeCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&recipe)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "404", "message": "Not Found", "reason": "Recipe not found"})
        return
    }

    c.JSON(http.StatusOK, recipe)
}

func (rc *RecipeController) SearchRecipes(c *gin.Context) {
    title := c.Param("title")
    if title == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "400", "message": "Bad Request", "reason": "Title parameter is required"})
        return
    }

    filter := bson.M{"title": bson.M{"$regex": title, "$options": "i"}}
    cursor, err := rc.recipeCollection.Find(context.TODO(), filter)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "500", "message": "Internal Server Error", "reason": "Failed to search recipes"})
        return
    }
    defer cursor.Close(context.TODO())

    var recipes []models.Recipe
    for cursor.Next(context.TODO()) {
        var recipe models.Recipe
        cursor.Decode(&recipe)
        recipes = append(recipes, recipe)
    }

    c.JSON(http.StatusOK, recipes)
}

func (rc *RecipeController) CreateRecipe(c *gin.Context) {
	var recipe models.Recipe

	err := c.Request.ParseMultipartForm(10 << 20) 
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "400", "message": "Bad Request", "reason": "Failed to parse form"})
		return
	}

	recipe.Title = c.PostForm("title")
	recipe.PrepTime = c.PostForm("prep_time")
	recipe.Servings = c.PostForm("servings")
	recipe.Difficulty = c.PostForm("difficulty")

	ingredients := c.PostForm("ingredients")
	instructions := c.PostForm("instructions")
	if err := json.Unmarshal([]byte(ingredients), &recipe.Ingredients); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "400", "message": "Bad Request", "reason": "Invalid ingredients format"})
		return
	}
	if err := json.Unmarshal([]byte(instructions), &recipe.Instructions); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "400", "message": "Bad Request", "reason": "Invalid instructions format"})
		return
	}

	imageData := c.PostForm("image")
	if imageData == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "400", "message": "Bad Request", "reason": "Image is required"})
		return
	}

	imageBytes, err := base64.StdEncoding.DecodeString(imageData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "400", "message": "Bad Request", "reason": "Invalid image data"})
		return
	}
	recipe.Image = imageBytes

	email := c.GetString("email")

	username, err := utils.GetUsernameByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "500", "message": "Internal Server Error", "reason": "Could not retrieve username"})
		return
	}

	recipe.ID = primitive.NewObjectID()
	recipe.CreatedAt = time.Now()
	recipe.CreatedBy = username

	if err := rc.validate.Struct(recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "400", "message": "Bad Request", "reason": err.Error()})
		return
	}

	_, err = rc.recipeCollection.InsertOne(context.TODO(), recipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "500", "message": "Internal Server Error", "reason": "Failed to create recipe"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Recipe created successfully"})
}

func (rc *RecipeController) UpdateRecipe(c *gin.Context) {
    id := c.Param("id")
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "400", "message": "Bad Request", "reason": "Invalid recipe ID"})
        return
    }

    var recipe models.Recipe
    if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "400", "message": "Bad Request", "reason": "Failed to parse form"})
        return
    }

    recipe.Title = c.PostForm("title")
    recipe.PrepTime = c.PostForm("prep_time")
    recipe.Servings = c.PostForm("servings")
    recipe.Difficulty = c.PostForm("difficulty")

    ingredients := c.PostForm("ingredients")
    instructions := c.PostForm("instructions")
    if err := json.Unmarshal([]byte(ingredients), &recipe.Ingredients); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "400", "message": "Bad Request", "reason": "Invalid ingredients format"})
        return
    }
    if err := json.Unmarshal([]byte(instructions), &recipe.Instructions); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "400", "message": "Bad Request", "reason": "Invalid instructions format"})
        return
    }

    imageData := c.PostForm("image")
    if imageData != "" {
        imageBytes, err := base64.StdEncoding.DecodeString(imageData)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "400", "message": "Bad Request", "reason": "Invalid image data"})
            return
        }
        recipe.Image = imageBytes
    }

    email := c.GetString("email")
    username, err := utils.GetUsernameByEmail(email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "500", "message": "Internal Server Error", "reason": "Could not retrieve username"})
        return
    }

    recipe.CreatedBy = username 

    if err := rc.validate.Struct(recipe); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "400", "message": "Bad Request", "reason": err.Error()})
        return
    }

    _, err = rc.recipeCollection.UpdateOne(context.TODO(), bson.M{"_id": objID}, bson.D{{"$set", recipe}})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "500", "message": "Internal Server Error", "reason": "Failed to update recipe"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Recipe updated successfully"})
}

func (rc *RecipeController) DeleteRecipe(c *gin.Context) {
    id := c.Param("id")
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "400", "message": "Bad Request", "reason": "Invalid recipe ID"})
        return
    }

    _, err = rc.recipeCollection.DeleteOne(context.TODO(), bson.M{"_id": objID})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "500", "message": "Internal Server Error", "reason": "Failed to delete recipe"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Recipe deleted successfully"})
}
