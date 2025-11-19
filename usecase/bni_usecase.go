package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"router-template/entities"
	"router-template/entities/common"
	"strings"
	"time"

	"github.com/kpango/glg"
)

type BniUsecase interface {
	GetTokenBni(entities.BniToken) (message map[string]interface{}, er error)
	GetBalance(entities.BniToken) (message map[string]interface{}, er error)
}

func NewBniUsecase() BniUsecase {
	return &bniUsecase{}
}

type bniUsecase struct{}
type bodyType struct {
	AccountNo string `json:"accountNo"`
	ClientId  string `json:"clientId"`
}
type bodyReq struct {
	AccountNo string `json:"accountNo"`
	ClientId  string `json:"clientId"`
	Signature string `json:"signature"`
}

func (e *bniUsecase) GetTokenBni(payload entities.BniToken) (message map[string]interface{}, er error) {
	payload.ClientId = os.Getenv("bni.client_id")
	payload.ClientSecret = os.Getenv("bni.client_secret")
	payload.GrantType = "client_credentials"

	glg.Info("Payload:", payload)
	data := url.Values{}
	data.Set("grant_type", payload.GrantType)

	req, err := http.NewRequest(
		"POST",
		"https://sandbox.bni.co.id/api/oauth/token",
		strings.NewReader(data.Encode()),
	)

	if err != nil {
		return message, err
	}

	glg.Info(req)
	req.SetBasicAuth(payload.ClientId, payload.ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	glg.Info(req.Header)
	client := &http.Client{
		Timeout: time.Second * 30,
	}

	resp, err := client.Do(req)
	if err != nil {
		return message, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return message, err
	}
	err = json.Unmarshal(respBody, &message)

	return message, err
}

func (e *bniUsecase) GetBalance(payload entities.BniToken) (message map[string]interface{}, er error) {
	payload.ApiSecret = os.Getenv("bni.api_secret")
	payload.AppName = os.Getenv("bni.app_name")

	body := map[string]any{}
	body["accountNo"] = payload.AccountNo
	body["clientId"] = common.GenerateClientId(payload.AppName)
	signature := common.GenerateSignature(body, payload.ApiSecret)
	reqBody := bodyReq{
		AccountNo: payload.AccountNo,
		ClientId:  common.GenerateClientId(payload.AppName),
		Signature: signature,
	}
	dataBody, err := json.Marshal(reqBody)
	if err != nil {
		return message, err
	}

	glg.Info("Payload : ", payload)
	glg.Info("Body : ", string(dataBody))
	dataUrl := url.Values{}
	dataUrl.Set("access_token", payload.AccessToken)
	fullURL := fmt.Sprintf("%s?%s", "https://sandbox.bni.co.id/H2H/v2/getbalance", dataUrl.Encode())

	glg.Info("Url : ", fullURL)
	req, err := http.NewRequest(
		"POST",
		fullURL,
		bytes.NewBuffer(dataBody),
	)

	if err != nil {
		return message, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", payload.ApiKey)

	glg.Info("header : ", req.Header)
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	err = json.Unmarshal(respBody, &message)
	glg.Info("Message : ", message)
	glg.Info("Status : ", resp.StatusCode)
	return message, err
}
