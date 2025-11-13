package handler

import (
	"net/http"
	"router-template/delivery"
	"router-template/entities"
	"router-template/entities/common"
	"router-template/usecase"

	"github.com/gin-gonic/gin"
)

func TopupBalance(ctx *gin.Context) {
	payload := entities.TopupBalance{}
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

	resultToken, invalidToken := common.CheckToken(payload.Token)
	if invalidToken != nil {
		token, er := common.GetToken(payload.Balance)
		if er != nil {
			delivery.PrintError(er.Error())
			ctx.String(http.StatusInternalServerError, "Token gagal digenerate")
			return
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message":   resultToken,
				"new_token": token,
			})
			return
		}
	}

	ucase := usecase.NewBalancesUsecase()
	balance, er := ucase.TopupBalance(payload.Id, payload.Balance)
	if er != nil {
		delivery.PrintError(er.Error())
		ctx.String(http.StatusInternalServerError, "internal service error")
	} else {
		ctx.JSON(http.StatusOK, balance)
	}
}
