package notificationrepo

import (
	"errors"
	"os"
	"router-template/repository/built_in/databasefactory"
)

func NewNotificationRepo() (NotificationRepo, error) {
	driverName := os.Getenv("app.database_driver")
	if driverName == databasefactory.DRIVER_MYSQL {
		return newNotificationMysqlImpl(), nil
	} else {
		return nil, errors.New("unimplemented database driver")
	}

}
