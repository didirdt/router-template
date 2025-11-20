package balancerepo

import (
	"database/sql"
	"errors"
	"fmt"
	"router-template/entities"
	"router-template/repository/built_in/databasefactory"
	"router-template/repository/employeerepo"
	"sync"
)

func newBalanceMysqlImpl() BalanceRepo {
	conn := databasefactory.AppDb.GetConnection()
	return &balanceMysqlImpl{conn: conn.(*sql.DB)}
}

type balanceMysqlImpl struct {
	conn *sql.DB
	mu   sync.RWMutex
}

func (e *balanceMysqlImpl) TopupBalance(id int64, balance float64) (employee entities.Employee, er error) {
	e.mu.Lock()
	defer e.mu.Unlock()

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

// func (b *balanceMysqlImpl) SendBalance(balances entities.SendBalance, employeeBalance *entities.EmployeeBalance, ch chan *entities.EmployeeBalance, wg *sync.WaitGroup) {
// 	emprepo, _ := employeerepo.NewEmployeeRepo()
// 	employee, er := emprepo.GetEmployeeById(balances.Id)
// 	if er != nil {
// 		message := fmt.Sprintf("error while Process payment send Balance from ID %d : %s", balances.Id, er.Error())
// 		employeeBalance.Id = balances.Id
// 		employeeBalance.Message = message
// 		ch <- employeeBalance
// 		return
// 	}

// 	toEmployee, er := emprepo.GetEmployeeById(balances.ToId)
// 	if er != nil {
// 		message := fmt.Sprintf("error while Process payment send Balance Balance To ID %d : %s", balances.ToId, er.Error())
// 		employeeBalance.Id = balances.ToId
// 		employeeBalance.Message = message
// 		ch <- employeeBalance
// 		return
// 	}

// 	employeeBalance.Mutex.Lock()
// 	defer wg.Done()
// 	defer employeeBalance.Mutex.Unlock()

// 	employeeBalance.Id = employee.Id
// 	employeeBalance.Name = employee.Name
// 	employeeBalance.Balance = employee.Balance - balances.Balance

// 	if employeeBalance.Balance < 0 {
// 		employeeBalance.Message = fmt.Sprintf("error while Process payment send Balance : %v -  tidak cukup", employeeBalance.Balance)
// 		ch <- employeeBalance
// 		return
// 	}

// 	result, er := b.conn.Exec(`UPDATE employee SET balance=? WHERE id=?`, employeeBalance.Balance, employee.Id)
// 	if result == nil || er != nil {
// 		employeeBalance.Message = fmt.Sprint("error while Process payment send Balance : %s - %s"+er.Error(), result)
// 		ch <- employeeBalance
// 		return
// 	}

// 	balance, er := b.TopupBalance(toEmployee.Id, balances.Balance)
// 	if er != nil {
// 		employeeBalance.Message = fmt.Sprint("error while Process payment send Balance : %s - %s"+er.Error(), balance)
// 		ch <- employeeBalance
// 		return
// 	}

// 	employeeBalance.Message = fmt.Sprint(
// 		"success send to : ", toEmployee.Id,
// 		" sebesar : ", balances.Balance,
// 		" nilai awal balance : ", employee.Balance,
// 		" nilai akhir balance : ", employeeBalance.Balance)

// 	ch <- employeeBalance
// }

func (b *balanceMysqlImpl) SendBalance(balances entities.SendBalance, ch chan *entities.EmployeeBalance, wg *sync.WaitGroup) {
	defer wg.Done()
	b.mu.Lock()
	defer b.mu.Unlock()

	emrepo, _ := employeerepo.NewEmployeeRepo()
	senderEmployee, _ := emrepo.GetEmployeeById(balances.Id)

	employeeBalance := &entities.EmployeeBalance{}
	employeeBalance.Id = senderEmployee.Id
	employeeBalance.Name = senderEmployee.Name

	tx, err := b.conn.Begin()
	if err != nil {
		message := fmt.Sprintf("error starting transaction: %s", err.Error())
		employeeBalance.Message = ErrorSendBalance(balances.Id, balances.ToId, balances.Balance, message)
		ch <- employeeBalance
		return
	}
	defer tx.Rollback()

	var senderBalance float64
	err = tx.QueryRow("SELECT balance FROM employee WHERE id = ? FOR UPDATE", balances.Id).Scan(&senderBalance)
	if err != nil {
		message := fmt.Sprintf("error sender: %s", err.Error())
		employeeBalance.Message = ErrorSendBalance(balances.Id, balances.ToId, balances.Balance, message)
		ch <- employeeBalance
		return
	}

	var receiverBalance float64
	err = tx.QueryRow("SELECT balance FROM employee WHERE id = ? FOR UPDATE", balances.ToId).Scan(&receiverBalance)
	if err != nil {
		message := fmt.Sprintf("error receiver: %s", err.Error())
		employeeBalance.Message = ErrorSendBalance(balances.Id, balances.ToId, balances.Balance, message)
		ch <- employeeBalance
		return
	}

	if senderBalance < balances.Balance {
		message := fmt.Sprintf("Saldo kurang: %.2f, butuh: %.2f", senderBalance, balances.Balance)
		employeeBalance.Message = ErrorSendBalance(balances.Id, balances.ToId, balances.Balance, message)
		ch <- employeeBalance
		return
	}

	newSenderBalance := senderBalance - balances.Balance
	_, err = tx.Exec("UPDATE employee SET balance = ? WHERE id = ?", newSenderBalance, balances.Id)
	if err != nil {
		message := fmt.Sprintf("error updating sender: %s", err.Error())
		employeeBalance.Message = ErrorSendBalance(balances.Id, balances.ToId, balances.Balance, message)
		ch <- employeeBalance
		return
	}

	newReceiverBalance := receiverBalance + balances.Balance
	_, err = tx.Exec("UPDATE employee SET balance = ? WHERE id = ?", newReceiverBalance, balances.ToId)
	if err != nil {
		message := fmt.Sprintf("error updating receiver: %s", err.Error())
		employeeBalance.Message = ErrorSendBalance(balances.Id, balances.ToId, balances.Balance, message)
		ch <- employeeBalance
		return
	}

	if err := tx.Commit(); err != nil {
		message := fmt.Sprintf("error committing: %s", err.Error())
		employeeBalance.Message = ErrorSendBalance(balances.Id, balances.ToId, balances.Balance, message)
		ch <- employeeBalance
		return
	}

	employeeBalance.Balance = newSenderBalance
	employeeBalance.Message = fmt.Sprint("Berhasil Kirim Balance ke ID : ", balances.ToId,
		" sebesar : ", balances.Balance,
		" nilai awal balance : ", senderBalance,
		" nilai akhir balance : ", newSenderBalance)

	ch <- employeeBalance
}

func ErrorSendBalance(id int64, toId int64, balance float64, message string) (serr string) {
	serr = fmt.Sprintf("%s dari ID : %d, ke ID : %d, sebesar : %.2f", message, id, toId, balance)
	return
}
