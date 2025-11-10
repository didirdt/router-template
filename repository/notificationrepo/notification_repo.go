package notificationrepo

import (
	"router-template/entities"
)

type NotificationRepo interface {
	SendNotification(notification entities.SendNotif) (entities.SendNotifResponse, error)
	ReceiveNotification(notification entities.SendNotif) ([]entities.SendNotifResponse, error)
}
