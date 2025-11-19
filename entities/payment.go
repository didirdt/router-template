package entities

import (
	"time"
)

type VirtualAccount struct {
	PartnerServiceId      string `json:"partner_service_id"`
	CustomerNo            string `json:"customer_no"`
	VirtualAccountNo      string `json:"virtual_account_no"`
	VirtualAccountName    string `json:"virtual_account_name"`
	VirtualAccountEmail   string `json:"virtual_account_email"`
	VirtualAccountPhone   string `json:"virtual_account_phone"`
	TrxId                 string `json:"trx_id"`
	TotalAmount           string `json:"total_amount"`
	Currency              string `json:"currency"`
	BankName              string `json:"bank_name"`
	VirtualAccountTrxType string `json:"virtual_account_trx_type"`
	ExpiredDate           string `json:"expired_date"`
	Description           string `json:"description"`
}

func (va *VirtualAccount) GenerateVirtualAccountNo() string {
	va.VirtualAccountNo = va.PartnerServiceId + va.CustomerNo
	return va.VirtualAccountNo
}

func (va *VirtualAccount) GetExpiredDate() string {
	return time.Now().Add(time.Hour * 24).Format("2006-01-02T15:04:05+07:00")
}

type Qris struct {
	PartnerReferenceNo string             `json:"partnerReferenceNo"`
	Amount             QrisAmount         `json:"amount"`
	MerchantId         string             `json:"merchantId"`
	TerminalId         string             `json:"terminalId"`
	ValidityPeriod     string             `json:"validityPeriod"`
	AdditionalInfo     additionalInfoQris `json:"additionalInfo"`
}
type QrisAmount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type additionalInfoQris struct {
	PostalCode     string `json:"postalCode"`
	FeeType        string `json:"feeType"`
	ValidityPeriod string `json:"validityPeriod"`
}

type CreateVaResponseQris struct {
	ResponseCode       string             `json:"responseCode"`
	ResponseMessage    string             `json:"responseMessage"`
	ReferenceNo        string             `json:"referenceNo"`
	PartnerReferenceNo string             `json:"partnerReferenceNo"`
	QrContent          string             `json:"qrContent"`
	TerminalId         string             `json:"terminalId"`
	AdditionalInfo     additionalInfoQris `json:"additionalInfo"`
}

func (va *Qris) GetValidityPeriod() string {
	return time.Now().Add(time.Hour * 24).Format("2006-01-02T15:04:05+07:00")
}

type BniToken struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	GrantType    string `json:"grantType"`
	ApiKey       string `json:"apiKey"`
	ApiSecret    string `json:"apiSecret"`
	AppName      string `json:"appName"`
	ClientId64   string `json:"clientId64"`
	AccessToken  string `json:"access_token"`
	AccountNo    string `json:"account_no"`
}
