package middlewares

import (
	"github.com/Safwanseban/Project-Ecommerce/auth"
	"github.com/gin-gonic/gin"
)
func AdminAuth() gin.HandlerFunc{
	return func(context *gin.Context) {
		// tokenString := context.GetHeader("Authorization")
		tokenString,err:=context.Cookie("Adminjwt")
		if tokenString == "" {
			context.JSON(401, gin.H{"error": "request does not contain an access token"})
			context.Abort()
			return
		}
		err= auth.ValidateToken(tokenString)
		if err != nil {
			context.JSON(401, gin.H{"error": err.Error()})
			context.Abort()
			return
		}
		context.Next()
	}
}

func UserAuth() gin.HandlerFunc{
	return func(context *gin.Context) {
		// tokenString := context.GetHeader("Authorization")
		tokenString,err:=context.Cookie("UserAuth")
		if tokenString == "" {
			context.JSON(401, gin.H{"error": "request does not contain an access token"})
			context.Abort()
			return
		}
		err= auth.ValidateToken(tokenString)
		context.Set("user",auth.P)
		if err != nil {
			context.JSON(401, gin.H{"error": err.Error()})
			context.Abort()
			return
		}
	
		context.Next()
	}
}