package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"router-template/entities"
	"router-template/entities/common"
	"time"

	dokuCommon "github.com/PTNUSASATUINTIARTHA-DOKU/doku-golang-library/commons"
	"github.com/PTNUSASATUINTIARTHA-DOKU/doku-golang-library/commons/utils"
	createVaModels "github.com/PTNUSASATUINTIARTHA-DOKU/doku-golang-library/models/va/createVa"
	"github.com/PTNUSASATUINTIARTHA-DOKU/doku-golang-library/services"
)

type DokuUsecase interface {
	PayWithVA(va entities.VirtualAccount) (createVaModels.CreateVaResponseDto, error)
	PayWithQris(va entities.Qris) (entities.CreateVaResponseQris, error)
}

func NewDokuUsecase() DokuUsecase {
	return &dokuUsecase{}
}

type dokuUsecase struct{}

func (d *dokuUsecase) PayWithVA(va entities.VirtualAccount) (response createVaModels.CreateVaResponseDto, er error) {
	resultVa := va.GenerateVirtualAccountNo()
	createVaRequestDTO := createVaModels.CreateVaRequestDto{
		PartnerServiceId:    va.PartnerServiceId,
		CustomerNo:          va.CustomerNo,
		VirtualAccountNo:    resultVa,
		VirtualAccountName:  va.VirtualAccountName,
		VirtualAccountEmail: va.VirtualAccountEmail,
		VirtualAccountPhone: va.VirtualAccountPhone,
		TrxId:               va.TrxId,
		TotalAmount: createVaModels.TotalAmount{
			Value:    va.TotalAmount,
			Currency: va.Currency,
		},
		AdditionalInfo: createVaModels.AdditionalInfo{
			Channel: va.BankName,
			VirtualAccountConfig: createVaModels.VirtualAccountConfig{
				ReusableStatus: true,
			},
		},
		VirtualAccountTrxType: va.VirtualAccountTrxType,
		ExpiredDate:           va.GetExpiredDate(),
		FreeTexts: []createVaModels.FreeTexts{
			{
				Indonesia: va.Description,
			},
		},
	}

	snap, _ := common.GetDoku()
	response, err := snap.CreateVa(createVaRequestDTO)

	return response, err
}

func (d *dokuUsecase) PayWithQris(qris entities.Qris) (responBody entities.CreateVaResponseQris, er error) {
	snap, _ := common.GetDoku()
	var vaServices services.VaServices
	var tokenServices services.TokenServices
	var snapUtils utils.SnapUtils

	timestamp := tokenServices.GenerateTimestamp()
	externalId := snapUtils.GenerateExternalId()

	qris.ValidityPeriod = qris.GetValidityPeriod()
	RequestBody, err := json.Marshal(qris)
	if err != nil {
		return responBody, fmt.Errorf("error body response: %w", err)
	}

	endPointUrl := "/snap-adapter/b2b/v1.0/qr/qr-mpm-generate"
	httpMethod := "POST"
	tokenB2B := snap.GetTokenB2B().AccessToken
	secretKey := snap.SecretKey
	clientId := snap.ClientId
	signature := tokenServices.GenerateSymetricSignature(httpMethod, endPointUrl, tokenB2B, RequestBody, timestamp, secretKey)
	requestHeader := vaServices.GenerateRequestHeaderDto("H2H", signature, timestamp, clientId, externalId, tokenB2B)

	url := dokuCommon.Config.GetBaseUrl(dokuCommon.Config{}, snap.IsProduction) + endPointUrl
	fmt.Println("URL: ", string(url))
	header := map[string]string{
		"X-PARTNER-ID":  requestHeader.XPartnerId,
		"X-TIMESTAMP":   requestHeader.XTimestamp,
		"X-SIGNATURE":   requestHeader.XSignature,
		"Authorization": "Bearer " + requestHeader.Authorization,
		"X-EXTERNAL-ID": requestHeader.XExternalId,
		"CHANNEL-ID":    requestHeader.ChannelId,
		"Content-Type":  "application/json",
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(RequestBody))
	if err != nil {
		return responBody, fmt.Errorf("error body response: %w", err)
	}

	for key, value := range header {
		req.Header.Set(key, value)
	}

	fmt.Println("URL: ", req.Header)
	client := &http.Client{
		Timeout: time.Second * 30,
	}

	resp, err := client.Do(req)
	if err != nil {
		return responBody, fmt.Errorf("error body response: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	fmt.Println("RESPONSE: ", string(respBody))
	if err := json.Unmarshal(respBody, &responBody); err != nil {
		return responBody, fmt.Errorf("error body response: %w", err)
	}

	return responBody, err
}
