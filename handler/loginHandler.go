package handler

import (
	"fmt"
	models2 "github.com/Assyl00/goProject/internal/models"
	"github.com/Assyl00/goProject/token"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	UserName string `json:"UserName" form:"UserName" binding:"required"`
	Password string `json:"Password" form:"Password" binding:"required"`
}

func LoginHandler(context *gin.Context) {
	var loginObj LoginRequest
	if err := context.ShouldBindJSON(&loginObj); err != nil {
		var errors []models2.ErrorDetail = make([]models2.ErrorDetail, 0, 1)
		errors = append(errors, models2.ErrorDetail{
			ErrorType:    models2.ErrorTypeValidation,
			ErrorMessage: fmt.Sprintf("%v", err),
		})
		badRequest(context, http.StatusBadRequest, "invalid request", errors)
	}

	var claims = &models2.JwtClaims{}
	claims.Username = loginObj.UserName
	claims.Roles = []int{1, 2, 3}
	claims.Audience = context.Request.Header.Get("Referer")

	var tokenCreationTime = time.Now().UTC()
	var expirationTime = tokenCreationTime.Add(time.Duration(2) * time.Hour)
	tokeString, err := token.GenrateToken(claims, expirationTime)

	if err != nil {
		badRequest(context, http.StatusBadRequest, "error in generating token", []models2.ErrorDetail{
			{
				ErrorType:    models2.ErrorTypeError,
				ErrorMessage: err.Error(),
			},
		})
	}

	ok(context, http.StatusOK, "token created", tokeString)
}
