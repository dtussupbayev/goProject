package middleware

import (
	"github.com/Assyl00/goProject/internal/models"
	"github.com/Assyl00/goProject/token"
	"github.com/gin-gonic/gin"
	//"fmt"
	"net/http"
)

func ReturnUnauthorized(context *gin.Context) {
	context.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{
		Error: []models.ErrorDetail{
			{
				ErrorType:    models.ErrorTypeUnauthorized,
				ErrorMessage: "You are not authorized to access this path",
			},
		},
		Status:  http.StatusUnauthorized,
		Message: "Unauthorized access",
	})
}

func ValidateToken() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString := context.Request.Header.Get("apikey")
		referer := context.Request.Header.Get("Referer")
		valid, claims := token.VerifyToken(tokenString, referer)
		if !valid {
			ReturnUnauthorized(context)
		}
		if len(context.Keys) == 0 {
			context.Keys = make(map[string]interface{})
		}
		context.Keys["Username"] = claims.Username
		//context.Keys["Roles"] = claims.Roles

	}
}

func Authorization(validRoles []int) gin.HandlerFunc {
	return func(context *gin.Context) {

		if len(context.Keys) == 0 {
			ReturnUnauthorized(context)
		}

	}
}
