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

func (b *balanceMysqlImpl) SendBalance(balances entities.SendBalance, ch chan entities.EmployeeBalance) (employeeBalance entities.EmployeeBalance, er error) {
	emprepo, _ := employeerepo.NewEmployeeRepo()
	employee, er := emprepo.GetEmployeeById(balances.Id)
	if er != nil {
		message := fmt.Sprintf("error while Process payment send Balance : %s", er.Error())
		ch <- entities.EmployeeBalance{Id: balances.ToId, Message: message}
		return
		// return entities.EmployeeBalance{Id: balances.Id}, errors.New(fmt.Sprint("error while Process payment from employee : ", er.Error()))
	}

	toEmployee, er := emprepo.GetEmployeeById(balances.ToId)
	if er != nil {
		message := fmt.Sprintf("error while Process payment send Balance : %s", er.Error())
		ch <- entities.EmployeeBalance{Id: balances.ToId, Message: message}
		return
		// return entities.EmployeeBalance{Id: balances.ToId}, errors.New(fmt.Sprint("error while Process payment to employee : ", er.Error()))
	}

	employeeBalance = entities.EmployeeBalance{
		Id:      employee.Id,
		Name:    employee.Name,
		Balance: employee.Balance,
	}

	employeeBalance.Mutex.Lock()
	employeeBalance.Balance = employee.Balance - balances.Balance
	if employeeBalance.Balance < 0 {
		employeeBalance.Message = fmt.Sprintf("error while Process payment send Balance : %v -  tidak cukup", employeeBalance.Balance)
		ch <- employeeBalance
		return
		// return employeeBalance, errors.New("error while Process payment send Balance : balance not enough")
	}

	result, er := b.conn.Exec(`UPDATE employee SET balance=? WHERE id=?`, employeeBalance.Balance, employee.Id)
	if result == nil || er != nil {
		employeeBalance.Message = fmt.Sprint("error while Process payment send Balance : %s - %s"+er.Error(), result)
		ch <- employeeBalance
		return
		// return employeeBalance, errors.New(fmt.Sprint("error while Process payment send Balance : ", er.Error()))
	}

	balance, er := b.TopupBalance(toEmployee.Id, balances.Balance)
	if er != nil {
		employeeBalance.Message = fmt.Sprint("error while Process payment send Balance : %s - %s"+er.Error(), balance)
		ch <- employeeBalance
		return
		// return employeeBalance, errors.New(fmt.Sprint("error while Process payment send Balance :", balance, er.Error()))
	}
	// fmt.Printf("Reload balances: %+v\n", employeeBalance)
	employeeBalance.Mutex.Unlock()

	employeeBalance.Message = fmt.Sprint(
		"success send to : ", toEmployee.Id,
		" sebesar : ", balances.Balance,
		" nilai awal balance : ", employee.Balance,
		" nilai akhir balance : ", employeeBalance.Balance)

	ch <- employeeBalance
	return
}
