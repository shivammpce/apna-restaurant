package middleware

import (
	"net/http"
	"os"
	"strings"
	"time"

	"apna-restaurant-2.0/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.GetHeader("Authorization")
		authToken = strings.TrimSpace(authToken)
		if len(authToken) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization"})
			c.Abort()
			return
		}
		splittedToken := strings.Split(authToken, " ")
		if len(splittedToken) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Malformed Authorization token"})
			c.Abort()
			return
		}
		tokenObj := &utils.Token{}
		parsedToken, err := jwt.ParseWithClaims(splittedToken[1], tokenObj, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("TOKEN_PASSWORD")), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Malformed Authorization"})
			c.Abort()
			return
		}
		if !parsedToken.Valid {
			c.JSON(http.StatusForbidden, gin.H{"error": "Resource Forbidden"})
			c.Abort()
			return
		}
		if tokenObj.ExpiresAt < time.Now().Local().Unix() {
			c.JSON(http.StatusForbidden, gin.H{"error": "Token Expired"})
			c.Abort()
			return
		}
		c.Next()
	}
}
