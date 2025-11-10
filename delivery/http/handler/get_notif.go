package handler

import (
	"net/http"
	"router-template/delivery"
	"router-template/entities"
	"router-template/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetNotif(ctx *gin.Context) {
	payloads := entities.SendNotif{}
	id, er := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if er != nil {
		ctx.String(http.StatusBadRequest, "invalid id parameter")
		return
	}

	payloads.Id = id
	ucase := usecase.NewNotificationUsecase()
	notifications, er := ucase.ReceiveNotification(payloads)
	if er != nil {
		delivery.PrintError(er.Error())
		ctx.String(http.StatusInternalServerError, "internal service error")
	} else {
		ctx.JSON(http.StatusOK, notifications)
	}
}
