package main

import (
    "log"
    "gitlab.com/ltp2-b-crepusculo/ps-backend-enzo-rodrigues/internal/routes"
    "gitlab.com/ltp2-b-crepusculo/ps-backend-enzo-rodrigues/internal/database"
    "github.com/gin-gonic/gin"
)

func main() {
    database.ConnectDB()

    router := gin.Default()

    routes.InitializeRoutes(router)

    log.Fatal(router.Run(":8080"))
}
