package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenValue, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message:": "no token in cookie found"})
			c.Abort()
			return
		}

		claims, err := VerifyToken(tokenValue)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message:": "unauthorized"})
			c.Abort()
			return
		}

		idStr, ok := claims["id"].(string)
		if !ok {
			c.JSON(401, gin.H{"message": "unauthorized: invalid claims"})
			c.Abort()
			return
		}

		objectID, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			c.JSON(401, gin.H{"message": "unauthorized: invalid ObjectID"})
			c.Abort()
			return
		}

		c.Set("id", objectID)

		c.Next()
	}
}
