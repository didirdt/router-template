package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func TestDoku(ctx *gin.Context) {

	payload, err := ctx.GetQuery("id")
	if err {
		ctx.String(http.StatusBadRequest, "Failed Parse Params")
	}

	ctx.JSON(http.StatusOK, gin.H{"payload": payload})
}
