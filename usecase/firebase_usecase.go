package usecase

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

type FirebaseUsecase interface {
	FirebaseTest(token string, title string, body string, data string) (string, error)
}

func NewFirebaseUsecase() FirebaseUsecase {
	return &firebaseUsecase{}
}

type firebaseUsecase struct{}

func (e *firebaseUsecase) FirebaseTest(token string, title string, body string, data string) (message string, er error) {

	opt := option.WithCredentialsFile("serviceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return "Got error initializing app", fmt.Errorf("error initializing app: %v", err)
	}

	fcmClient, err := app.Messaging(context.TODO())
	if err != nil {
		log.Fatalf("messaging: %s", err)
	}

	ctx := context.Background()
	response, er := fcmClient.Send(ctx, &messaging.Message{
		Data: map[string]string{
			"data": data,
		},
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Token: token, // a token that you received from a client
	})

	return response, er
}
