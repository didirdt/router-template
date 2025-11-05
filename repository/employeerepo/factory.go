package employeerepo

import (
	"errors"
	"os"
	"router-template/repository/built_in/databasefactory"
)

func NewEmployeeRepo() (EmployeeRepo, error) {
	driverName := os.Getenv("app.database_driver")

	switch i := driverName; {
	case i == databasefactory.DRIVER_MOCK:
		return newEmployeeMockImpl(), nil
	case i == databasefactory.DRIVER_MYSQL:
		return newEmployeeMysqlImpl(), nil
	default:
		return nil, errors.New("unimplemented database driver")
	}
}
