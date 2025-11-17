package bni

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type params struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

func GetTokenBni(ctx *gin.Context) {
	payload := params{}
	// er := ctx.BindJSON(&payload)
	payload.ClientId = os.Getenv("bni.client_id")
	payload.ClientSecret = os.Getenv("bni.secret_key")
	payload.GrantType = "client_credentials"

	body := map[string]string{
		"grant_type": payload.GrantType,
	}
	RequestBody, _ := json.Marshal(body)

	req, err := http.NewRequest(
		"POST",
		"https://sandbox.bni.co.id/api/oauth/token",
		bytes.NewBuffer(RequestBody),
	)

	if err != nil {
		ctx.String(http.StatusBadRequest, "Failed Parse Params: "+err.Error())
		return
	}

	req.SetBasicAuth(payload.ClientId, payload.ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{
		Timeout: time.Second * 30,
	}
	resp, err := client.Do(req)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Failed Parse Params: "+err.Error())
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Failed request Token: "+err.Error())
		return
	}

	fmt.Println("RESPONSE: ", string(respBody))

	ctx.JSON(http.StatusOK, gin.H{"response": string(respBody)})
	// if er != nil {
	// 	ctx.String(http.StatusBadRequest, "Failed Parse Params: "+er.Error())
	// 	return
	// } else {
	// }
}
