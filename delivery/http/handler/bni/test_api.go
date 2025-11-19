package bni

import (
	"net/http"
	"router-template/delivery"
	"router-template/entities"
	"router-template/usecase"

	"github.com/gin-gonic/gin"
)

func GetTokenBni(ctx *gin.Context) {

	payload := entities.BniToken{}
	var er error
	if ctx.ContentType() == "application/json" {
		er = ctx.BindJSON(&payload)
	} else {
		er = ctx.Bind(&payload)
	}
	if er != nil {
		ctx.String(http.StatusBadRequest, "Failed Parse Params: "+er.Error())
		return
	}

	usecase := usecase.NewBniUsecase()
	rawjson, er := usecase.GetTokenBni(payload)
	if er != nil {
		delivery.PrintError(er.Error())
		ctx.String(http.StatusInternalServerError, "internal service error: "+er.Error())
	} else {
		ctx.JSON(http.StatusOK, gin.H{"response": rawjson})
	}
}

func GetBalance(ctx *gin.Context) {
	payload := entities.BniToken{}
	var er error
	if ctx.ContentType() == "application/json" {
		er = ctx.BindJSON(&payload)
	} else {
		er = ctx.Bind(&payload)
	}
	if er != nil {
		ctx.String(http.StatusBadRequest, "Failed Parse Params: "+er.Error())
		return
	}

	usecase := usecase.NewBniUsecase()
	rawjson, er := usecase.GetBalance(payload)
	if er != nil {
		delivery.PrintError(er.Error())
		ctx.String(http.StatusInternalServerError, "internal service error: "+er.Error())
	} else {
		ctx.JSON(http.StatusOK, gin.H{"response": rawjson})
	}
}
