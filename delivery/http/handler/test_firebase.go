package handler

import (
	"encoding/json"
	"net/http"
	"router-template/usecase"

	"github.com/gin-gonic/gin"
)

func TestFirebase(ctx *gin.Context) {
	type DataNotif struct {
		Name string `json:"name"`
	}

	type params struct {
		Token string    `json:"token_device"`
		Title string    `json:"title"`
		Body  string    `json:"body"`
		Data  DataNotif `json:"data"`
	}

	payload := params{}
	err := ctx.BindJSON(&payload)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Failed Parse Params: "+err.Error())
	}

	jsonBytes, err := json.Marshal(payload.Data)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Failed Parse Params: "+err.Error())
	}

	usecase := usecase.NewFirebaseUsecase()
	firebase, _ := usecase.FirebaseTest(payload.Token, payload.Title, payload.Body, string(jsonBytes))

	ctx.JSON(http.StatusOK, gin.H{"Firebase": firebase, "payload": payload})
}
