package balancerepo

import (
	"errors"
	"os"
	"router-template/repository/built_in/databasefactory"
)

func NewBalanceRepo() (BalanceRepo, error) {
	driverName := os.Getenv("app.database_driver")
	if driverName == databasefactory.DRIVER_MYSQL {
		return newBalanceMysqlImpl(), nil
	} else {
		return nil, errors.New("unimplemented database driver")
	}

}
