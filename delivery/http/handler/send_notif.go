package handler

import (
	"net/http"
	"router-template/delivery"
	"router-template/entities"
	"router-template/usecase"

	"github.com/gin-gonic/gin"
)

func SendNotif(ctx *gin.Context) {
	var payloads []entities.SendNotif

	var er error
	if ctx.ContentType() == "application/json" {
		er = ctx.BindJSON(&payloads)
	} else {
		er = ctx.Bind(&payloads)
	}
	if er != nil {
		ctx.String(http.StatusBadRequest, er.Error())
	}

	ucase := usecase.NewNotificationUsecase()
	notifications, er := ucase.SendNotification(payloads)
	if er != nil {
		delivery.PrintError(er.Error())
		ctx.String(http.StatusInternalServerError, "internal service error")
	} else {
		ctx.JSON(http.StatusOK, notifications)
	}
}
