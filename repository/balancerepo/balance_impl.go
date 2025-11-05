package balancerepo

import (
	"database/sql"
	"errors"
	"fmt"
	"router-template/entities"
	"router-template/repository/built_in/databasefactory"
	"router-template/repository/employeerepo"
)

func newBalanceMysqlImpl() BalanceRepo {
	conn := databasefactory.AppDb.GetConnection()
	return &balanceMysqlImpl{conn: conn.(*sql.DB)}
}

type balanceMysqlImpl struct {
	conn *sql.DB
}

func (e *balanceMysqlImpl) TopupBalance(id int64, balance float64) (employee entities.Employee, er error) {
	emprepo, _ := employeerepo.NewEmployeeRepo()
	employee, er = emprepo.GetEmployeeById(id)
	if er != nil {
		return employee, errors.New(fmt.Sprint("error while get employee : ", er.Error()))
	}

	endBalance := employee.Balance + balance
	result, er := e.conn.Exec(`UPDATE employee SET balance=? WHERE id=?`, endBalance, id)
	if result == nil || er != nil {
		return employee, errors.New(fmt.Sprint("error while update employee : ", er.Error()))
	}

	employee, er = emprepo.GetEmployeeById(id)
	return
}

func (b *balanceMysqlImpl) SendBalance(balances entities.SendBalance) (employeeBalance entities.EmployeeBalance, er error) {
	emprepo, _ := employeerepo.NewEmployeeRepo()
	employee, er := emprepo.GetEmployeeById(balances.Id)
	if er != nil {
		return entities.EmployeeBalance{Id: balances.Id}, errors.New(fmt.Sprint("error while Process payment from employee : ", er.Error()))
	}

	toEmployee, er := emprepo.GetEmployeeById(balances.ToId)
	if er != nil {
		return entities.EmployeeBalance{Id: balances.ToId}, errors.New(fmt.Sprint("error while Process payment to employee : ", er.Error()))
	}

	employeeBalance = entities.EmployeeBalance{
		Id:      employee.Id,
		Name:    employee.Name,
		Balance: employee.Balance,
	}
	endBalance := employee.Balance - balances.Balance
	if endBalance < 0 {
		return employeeBalance, errors.New("error while Process payment send Balance : balance not enough")
	}

	result, er := b.conn.Exec(`UPDATE employee SET balance=? WHERE id=?`, endBalance, employee.Id)
	if result == nil || er != nil {
		return employeeBalance, errors.New(fmt.Sprint("error while Process payment send Balance : ", er.Error()))
	}

	balance, er := b.TopupBalance(toEmployee.Id, balances.Balance)
	if er != nil {
		return employeeBalance, errors.New(fmt.Sprint("error while Process payment send Balance :", balance, er.Error()))
	}

	return
}
