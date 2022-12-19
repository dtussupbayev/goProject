package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
)


type Product struct {
	Id int
	Name string
}

func GetAllProducts(context *gin.Context){
	context.AbortWithStatusJSON(http.StatusOK,[]Product{
		{
			Id: 1,
			Name: "Mascara L'Oreal Telescopic",
		},
		{
			Id: 2,
			Name: "Essence clear brow gel",
		},
	})
}
