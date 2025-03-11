package routes

import (
    "github.com/gin-gonic/gin"
    "gitlab.com/ltp2-b-crepusculo/ps-backend-enzo-rodrigues/internal/controllers"
    "gitlab.com/ltp2-b-crepusculo/ps-backend-enzo-rodrigues/internal/middlewares"
    "github.com/gin-contrib/cors"
)

func InitializeRoutes(router *gin.Engine) {
    recipeController := controllers.NewRecipeController()
    authController := controllers.NewAuthController()
    userController := controllers.NewUserController()

    config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Authorization", "Content-Type"}

    router.Use(cors.New(config))

    router.POST("/register", authController.Register)
    router.POST("/login", authController.Login)

    router.GET("/recipes", recipeController.GetAllRecipes)
    router.GET("/recipes/:id", recipeController.GetRecipeByID)
    router.GET("/recipes/search/:title", recipeController.SearchRecipes)

    router.GET("/:username", userController.GetProfileByUsername)

    auth := router.Group("/auth")
    auth.Use(middlewares.AuthMiddleware())
    {
        auth.GET("/profile", userController.GetProfile)
        auth.POST("/profile/image", userController.UpdateProfileImage)
        auth.POST("/create", recipeController.CreateRecipe)
        auth.PUT("/edit/:id", recipeController.UpdateRecipe)
        auth.DELETE("/delete/:id", recipeController.DeleteRecipe)
    }

}
