package handler

import (
	"net/http"
	"router-template/delivery"
	"router-template/entities"
	"router-template/usecase"

	"github.com/gin-gonic/gin"
)

func DeleteEmployeeHandler(ctx *gin.Context) {
	payload := entities.GetId{}
	var er error
	if ctx.ContentType() == "application/json" {
		er = ctx.BindJSON(&payload)
	} else {
		er = ctx.Bind(&payload)
	}
	if er != nil {
		ctx.String(http.StatusBadRequest, er.Error())
	}

	ucase := usecase.NewEmployeeUsecase()
	employee, er := ucase.DeleteEmployee(payload.Id)
	if er != nil {
		delivery.PrintError(er.Error())
		ctx.String(http.StatusInternalServerError, "internal service error")
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "Employee deleted successfully", "data": employee})
	}
}
