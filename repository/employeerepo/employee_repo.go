package employeerepo

import (
	"router-template/entities"
)

type EmployeeRepo interface {
	GetEmployee() ([]entities.Employee, error)
	GetEmployeeById(id int64) (entities.Employee, error)
	CreateEmployee(name, address, phone_number string) (entities.Employee, error)
	UpdateEmployee(id int64, name, address, phone_number string) (entities.Employee, error)
	DeleteEmployee(id int64) (entities.Employee, error)
}
