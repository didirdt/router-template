package usecase

import (
	"router-template/entities"
	"router-template/entities/app"
	"router-template/repository/employeerepo"
)

type EmployeeUsecase interface {
	GetEmployeeList() ([]entities.Employee, error)
	GetEmployee(id int64) (entities.Employee, error)
	CreateEmployee(name, address, phone_number string) (entities.Employee, error)
	UpdateEmployee(id int64, name, address, phone_number string) (entities.Employee, error)
	DeleteEmployee(id int64) (entities.Employee, error)
}

func NewEmployeeUsecase() EmployeeUsecase {
	return &employeeUsecase{}
}

type employeeUsecase struct{}

func (e *employeeUsecase) GetEmployeeList() (detail []entities.Employee, er error) {
	repo, _ := employeerepo.NewEmployeeRepo()
	detail, er = repo.GetEmployee()

	if er != nil {
		return detail, er
	}

	if len(detail) == 0 {
		return detail, app.ErrNoRecord
	}

	return
}

func (e *employeeUsecase) GetEmployee(id int64) (employee entities.Employee, er error) {
	repo, _ := employeerepo.NewEmployeeRepo()
	employee, er = repo.GetEmployeeById(id)

	return
}

func (e *employeeUsecase) CreateEmployee(name, address, phone_number string) (employee entities.Employee, er error) {
	repo, _ := employeerepo.NewEmployeeRepo()
	employee, er = repo.CreateEmployee(name, address, phone_number)

	return
}

func (e *employeeUsecase) UpdateEmployee(id int64, name, address, phone_number string) (employee entities.Employee, er error) {
	repo, _ := employeerepo.NewEmployeeRepo()
	employee, er = repo.UpdateEmployee(id, name, address, phone_number)
	if er != nil {
		return employee, er
	}

	return
}

func (e *employeeUsecase) DeleteEmployee(id int64) (employee entities.Employee, er error) {
	repo, _ := employeerepo.NewEmployeeRepo()
	employee, er = repo.DeleteEmployee(id)
	if er != nil {
		return employee, er
	}

	return
}
