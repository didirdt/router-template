package excel

import (
	"net/http"
	"router-template/delivery"
	"router-template/usecase"

	"github.com/gin-gonic/gin"
)

func SendBalance(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.String(http.StatusBadRequest, "Failed Parse Params: "+err.Error())
		return
	}

	ucase := usecase.NewBalancesUsecase()
	_, report, er := ucase.SendBalanceExcel(file)
	if er != nil {
		delivery.PrintError(er.Error())
		ctx.String(http.StatusInternalServerError, "internal service error")
	} else {
		ctx.JSON(http.StatusOK, gin.H{"report": report, "file": file.Filename})
	}
}
