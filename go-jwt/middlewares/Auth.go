package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/tuanpnt17/eleven-golang-projects/go-jwt/helpers"
	"net/http"
)

func Authenticate() gin.HandlerFunc {
	return func(context *gin.Context) {
		clientToken := context.GetHeader("Authorization")
		if clientToken == "" {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			context.Abort()
			return
		}

		claims, errMsg := helpers.ValidateToken(clientToken)
		if errMsg != "" {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errMsg})
			context.Abort()
			return
		}

		context.Set("email", claims.Email)
		context.Set("first_name", claims.FirstName)
		context.Set("last_name", claims.LastName)
		context.Set("uid", claims.Uid)
		context.Set("user_type", claims.UserType)
		context.Next()
	}

}
