package notificationrepo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"router-template/delivery/broker"
	"router-template/entities"
	"router-template/repository/built_in/databasefactory"
	"router-template/repository/built_in/keyvaluefactory"
	"router-template/repository/employeerepo"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/randyardiansyah25/libpkg/util/env"
	"github.com/redis/go-redis/v9"
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

	ctx := context.Background()
	rdb := keyvaluefactory.AppStore.GetStore().(*redis.Client)
	key := fmt.Sprint("notif_go_", notification.Id, "_", time.Now().Unix())
	err = rdb.Set(ctx, key, json_notification, 0).Err()
	if err != nil {
		notif.Status = "failed"
		return notif, err
	}

	notif.Status = "success"
	return notif, nil
}

func (e *notificationMysqlImpl) ReceiveNotification(notification entities.SendNotif) (notifications []entities.SendNotifResponse, err error) {
	emprepo, _ := employeerepo.NewEmployeeRepo()
	employee, err := emprepo.GetEmployeeById(notification.Id)
	if err != nil {
		return []entities.SendNotifResponse{}, errors.New(fmt.Sprint("error while Send notification : ", err.Error()))
	}

	ctx := context.Background()
	rdb := keyvaluefactory.AppStore.GetStore().(*redis.Client)

	get_key := fmt.Sprint("notif_go_", employee.Id, "_*")
	get_notifications := rdb.Keys(ctx, get_key)

	for _, notif := range get_notifications.Val() {
		get_value := rdb.Get(ctx, notif).Val()
		err = json.Unmarshal([]byte(get_value), &notification)

		notifications = append(notifications, entities.SendNotifResponse{
			Id:      notification.Id,
			Name:    employee.Name,
			Message: notification.Message,
			Status:  "success",
		})
	}
	return notifications, err
}
