package handler

import (
	"net/http"
	"router-template/delivery"
	"router-template/entities"
	"router-template/usecase"

	"github.com/gin-gonic/gin"
)

func UpdateEmployeeHandler(ctx *gin.Context) {
	payload := entities.Employee{}
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

	ucase := usecase.NewEmployeeUsecase()
	employee, er := ucase.UpdateEmployee(payload.Id, payload.Name, payload.Address, payload.PhoneNumber)
	if er != nil {
		delivery.PrintError(er.Error())
		ctx.String(http.StatusInternalServerError, er.Error())
	} else {
		ctx.JSON(http.StatusOK, employee)
	}
}
