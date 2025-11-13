package common

import (
	"github.com/PTNUSASATUINTIARTHA-DOKU/doku-golang-library/controllers"
	"github.com/PTNUSASATUINTIARTHA-DOKU/doku-golang-library/doku"
	"github.com/randyardiansyah25/libpkg/util/env"
)

func GetDoku() (snap doku.Snap, err error) {

	privateKey := env.GetString("doku.privateKey")
	publicKey := env.GetString("doku.publicKey")
	clientId := env.GetString("doku.clientId")
	secretKey := env.GetString("doku.secretKey")
	issuer := env.GetString("doku.issuer")
	isProduction := false

	doku.TokenController = controllers.TokenController{}
	doku.VaController = controllers.VaController{}
	doku.NotificationController = controllers.NotificationController{}
	doku.DirectDebitController = &controllers.DirectDebitController{}

	snap = doku.Snap{
		PrivateKey:   privateKey,
		ClientId:     clientId,
		IsProduction: isProduction,
		SecretKey:    secretKey,
		Issuer:       issuer,
		PublicKey:    publicKey,
	}

	return snap, err
}
