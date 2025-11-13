package entities

import "time"

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
	PartnerReferenceNo string             `json:"partner_reference_no"`
	Amount             QrisAmount         `json:"amount"`
	MerchantId         string             `json:"merchant_id"`
	TerminalId         string             `json:"terminal_id"`
	ValidityPeriod     string             `json:"validity_period"`
	AdditionalInfo     additionalInfoQris `json:"additional_info"`
}
type QrisAmount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type additionalInfoQris struct {
	PostalCode     string `json:"postal_code"`
	FeeType        string `json:"fee_type"`
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
