package entities

import (
	"router-template/flatbuffer/Notification"

	flatbuffers "github.com/google/flatbuffers/go"
)

type SendNotif struct {
	Id      int64  `json:"id"`
	Message string `json:"message"`
}

type SendNotifResponse struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

func (s *SendNotifResponse) SerializeNotification() (data []byte) {
	builder := flatbuffers.NewBuilder(0)

	nameOffset := builder.CreateString(s.Name)
	messageOffset := builder.CreateString(s.Message)
	statusOffset := builder.CreateString(s.Status)

	Notification.NotificationStart(builder)
	Notification.NotificationAddId(builder, s.Id)
	Notification.NotificationAddName(builder, nameOffset)
	Notification.NotificationAddMessage(builder, messageOffset)
	Notification.NotificationAddStatus(builder, statusOffset)
	NotificationOffset := Notification.NotificationEnd(builder)
	builder.Finish(NotificationOffset)

	return builder.FinishedBytes()
}

func (n *SendNotifResponse) DeserializeNotification(data []byte) (SendNotifResponse, error) {
	notification := Notification.GetRootAsNotification(data, 0)
	notifResponse := SendNotifResponse{}

	notifResponse.Id = notification.Id()
	notifResponse.Name = string(notification.Name())
	notifResponse.Message = string(notification.Message())
	notifResponse.Status = string(notification.Status())

	return notifResponse, nil
}
