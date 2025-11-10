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
	type resp struct {
		index int
		r     entities.SendNotifResponse
	}

	resultsCh := make(chan resp, len(notifications))
	notif_responses = make([]entities.SendNotifResponse, len(notifications))
	for i, notification := range notifications {
		i2, n := i, notification
		go func(idx int, notif entities.SendNotif) {
			notifrepo, _ := notificationrepo.NewNotificationRepo()
			notif_response, err := notifrepo.SendNotification(notif)
			if err != nil {
				notif_response.Message = "error while Send Notification : " + err.Error()
			}
			resultsCh <- resp{index: idx, r: notif_response}
		}(i2, n)
	}

	for i := 0; i < len(notifications); i++ {
		rr := <-resultsCh
		notif_responses[rr.index] = rr.r
	}
	close(resultsCh)

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
