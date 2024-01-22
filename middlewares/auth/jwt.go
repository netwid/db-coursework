package auth

import (
	"github.com/golang-jwt/jwt/v5"
	util "github.com/netwid/db-coursework/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware JWT is jwt middleware
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		claims := &util.Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return util.JwtKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.Set("id", claims.ID)
		c.Next()
	}
}
