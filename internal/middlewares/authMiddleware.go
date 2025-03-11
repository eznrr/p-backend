package middlewares

import (
    "strings"
    "net/http"
    "gitlab.com/ltp2-b-crepusculo/ps-backend-enzo-rodrigues/internal/utils"
    "github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "401", "message": "Unauthorized", "reason": "Non-existent token"})
            c.Abort()
            return
        }

        token := strings.TrimPrefix(authHeader, "Bearer ")
        if token == authHeader {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "401", "message": "Unauthorized", "reason": "Non-existent token"})
            c.Abort()
            return
        }

        claims, err := utils.ValidateJWT(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "401", "message": "Unauthorized", "reason": "Token Unauthorized"})
            c.Abort()
            return
        }

        c.Set("email", claims.Email)
        c.Next()
    }
}
