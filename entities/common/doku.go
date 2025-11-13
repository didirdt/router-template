package common

import (
	"os"

	"github.com/PTNUSASATUINTIARTHA-DOKU/doku-golang-library/controllers"
	"github.com/PTNUSASATUINTIARTHA-DOKU/doku-golang-library/doku"
	"github.com/randyardiansyah25/libpkg/util/env"
)

func GetDoku() (snap doku.Snap, err error) {

	privateKey := env.GetString("doku.private_key")
	publicKey := env.GetString("doku.public_key")
	clientId := env.GetString("doku.client_id")
	secretKey := env.GetString("doku.secret_key")
	issuer := env.GetString("doku.issuer")
	isProduction := false

	doku.TokenController = controllers.TokenController{}
	doku.VaController = controllers.VaController{}
	doku.NotificationController = controllers.NotificationController{}
	doku.DirectDebitController = &controllers.DirectDebitController{}

	dokuPrivKey, err := os.ReadFile(privateKey)
	if err != nil {
		return snap, err
	}

	dokuPublicKey, err := os.ReadFile(publicKey)
	if err != nil {
		return snap, err
	}

	snap = doku.Snap{
		PrivateKey:   string(dokuPrivKey),
		PublicKey:    string(dokuPublicKey),
		ClientId:     clientId,
		IsProduction: isProduction,
		SecretKey:    secretKey,
		Issuer:       issuer,
	}

	return snap, err
}
