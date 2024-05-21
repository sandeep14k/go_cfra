package middleware

import (
	"context"
	"fmt"
	"gin/database"
	helper "gin/helper"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users_data")

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorization header provided")})
			c.Abort()
			return
		}
		claims, err := helper.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("user_name", claims.User_name)
		c.Set("codeforce_handle", claims.Codeforce_handle)
		c.Set("uid", claims.Uid)
		c.Next()
	}
}
func CheckTokenValid() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid, ok := ctx.Get("uid")
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "the token is invalid"})
			return
		}
		c, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		count, err := userCollection.CountDocuments(c, bson.M{"user_id": uid})

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}
		if count < 1 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "user with this token not exsit"})
			ctx.Abort()
			return
		}
	}
}
