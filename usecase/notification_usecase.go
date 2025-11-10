package usecase

import (
	"errors"
	"fmt"
	"router-template/entities"
	"router-template/repository/notificationrepo"
)

type NotificationUsecase interface {
	SendNotification(notifications []entities.SendNotif) (notif_responses []entities.SendNotifResponse, er error)
	ReceiveNotification(notification entities.SendNotif) (notif_responses []entities.SendNotifResponse, er error)
}

func NewNotificationUsecase() NotificationUsecase {
	return &notificationUsecase{}
}

type notificationUsecase struct{}

func (b *notificationUsecase) SendNotification(notifications []entities.SendNotif) (notif_responses []entities.SendNotifResponse, er error) {
	for _, notification := range notifications {
		notifrepo, _ := notificationrepo.NewNotificationRepo()
		notif_response, er := notifrepo.SendNotification(notification)
		if er != nil {
			notif_response.Message = "error while Send Notification : " + er.Error()
		}
		notif_responses = append(notif_responses, notif_response)
	}

	return notif_responses, er
}

func (b *notificationUsecase) ReceiveNotification(notification entities.SendNotif) (notif_responses []entities.SendNotifResponse, er error) {
	notifrepo, _ := notificationrepo.NewNotificationRepo()
	notif_responses, er = notifrepo.ReceiveNotification(notification)
	if er != nil {
		return notif_responses, errors.New(fmt.Sprint("error while Receive Notification : ", er.Error()))
	}

	return notif_responses, er
}
