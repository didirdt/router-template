package notificationrepo

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"router-template/delivery/broker"
	"router-template/entities"
	"router-template/repository/built_in/databasefactory"
	"router-template/repository/employeerepo"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/randyardiansyah25/libpkg/util/env"
)

func newNotificationMysqlImpl() NotificationRepo {
	conn := databasefactory.AppDb.GetConnection()
	return &notificationMysqlImpl{conn: conn.(*sql.DB)}
}

type notificationMysqlImpl struct {
	conn *sql.DB
}

func (e *notificationMysqlImpl) SendNotification(notification entities.SendNotif) (entities.SendNotifResponse, error) {
	emprepo, _ := employeerepo.NewEmployeeRepo()
	employee, err := emprepo.GetEmployeeById(notification.Id)
	if err != nil {
		return entities.SendNotifResponse{}, errors.New(fmt.Sprint("error while Send notification : ", err.Error()))
	}

	json_notification, _ := json.Marshal(notification)
	xchange := env.GetString("rabbit.exchange")
	notif := entities.SendNotifResponse{
		Id:      notification.Id,
		Name:    employee.Name,
		Message: notification.Message,
	}

	err = broker.BrokerChannel.Publish(
		xchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        json_notification,
			Timestamp:   time.Now(),
		},
	)

	if err != nil {
		notif.Status = "failed"
		return notif, err
	}

	notif.Status = "success"
	return notif, nil
}

func (e *notificationMysqlImpl) ReceiveNotification(notification entities.SendNotif) (notifications []entities.SendNotifResponse, err error) {
	// emprepo, _ := employeerepo.NewEmployeeRepo()
	// employee, err := emprepo.GetEmployeeById(notification.Id)
	// if err != nil {
	// 	return []entities.SendNotifResponse{}, errors.New(fmt.Sprint("error while Send notification : ", err.Error()))
	// }

	// queueName := env.GetString("rabbit.queue_name")
	// message, _ := broker.BrokerChannel.Consume(
	// 	queueName,
	// 	"",
	// 	false,
	// 	false,
	// 	false,
	// 	false,
	// 	nil,
	// )

	// for msg := range message {
	// 	var prettyJson bytes.Buffer
	// 	if er := json.Indent(&prettyJson, msg.Body, "", "    "); er != nil {
	// 		delivery.PrintErrorf("We got message, it seems the message is not in json format : %s, error : %v\n", string(msg.Body), er)
	// 		msg.Reject(false)
	// 	} else {
	// 		delivery.PrintLogf("A notification message was received :\n%s\n", prettyJson.String())
	// 		notif := entities.SendNotifResponse{
	// 			Id:      notification.Id,
	// 			Name:    employee.Name,
	// 			Message: string(msg.Body),
	// 			Status:  "success",
	// 		}
	// 		notifications = append(notifications, notif)
	// 		msg.Ack(false)
	// 	}
	// }
	return notifications, err
}
