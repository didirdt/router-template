package handler

import (
	"net/http"
	"router-template/delivery"
	"router-template/entities"
	"router-template/usecase"

	"github.com/gin-gonic/gin"
)

func PayWithVa(ctx *gin.Context) {
	var payload entities.VirtualAccount
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

	usecase := usecase.NewDokuUsecase()
	response, er := usecase.PayWithVA(payload)
	if er != nil {
		delivery.PrintError(er.Error())
		ctx.String(http.StatusInternalServerError, "internal service error Doku: "+er.Error())
	} else {
		ctx.JSON(http.StatusOK, gin.H{"payload": payload, "response": response})
	}
}

func PayWithQris(ctx *gin.Context) {
	var payload entities.Qris
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

	usecase := usecase.NewDokuUsecase()
	response, er := usecase.PayWithQris(payload)
	if er != nil {
		delivery.PrintError(er.Error())
		ctx.String(http.StatusInternalServerError, "internal service error: "+er.Error())
	} else {
		ctx.JSON(http.StatusOK, gin.H{"payload": payload, "response": response})
	}
}
