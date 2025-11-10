package handler

import (
	"net/http"
	"router-template/delivery"
	"router-template/entities"
	"router-template/entities/app"
	"router-template/entities/statuscode"
	"router-template/usecase"

	"github.com/gin-gonic/gin"
)

func CreateEmployeeHandler(ctx *gin.Context) {
	payload := entities.CreateEmployee{}
	var er error
	if ctx.ContentType() == "application/json" {
		er = ctx.BindJSON(&payload)
	} else {
		er = ctx.Bind(&payload)
	}
	if er != nil {
		ctx.String(http.StatusBadRequest, "Failed Parse Params: "+er.Error())
	}

	ucase := usecase.NewEmployeeUsecase()
	employee, er := ucase.CreateEmployee(payload.Name, payload.Address, payload.PhoneNumber)
	if er != nil {
		if er == app.ErrDuplicateEntry {
			ctx.String(statuscode.StatusDuplicate, "Data karyawan sudah tersedia!")
		} else {
			delivery.PrintError(er.Error())
			ctx.String(http.StatusInternalServerError, er.Error())
		}
	} else {
		ctx.JSON(http.StatusOK, employee)
	}
}
