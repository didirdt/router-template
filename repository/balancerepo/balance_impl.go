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
	mu   sync.Mutex
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

	employeeBalance := &entities.EmployeeBalance{}

	// Start transaction
	tx, err := b.conn.Begin()
	if err != nil {
		employeeBalance.Message = fmt.Sprintf("error starting transaction: %s", err.Error())
		ch <- employeeBalance
		return
	}
	defer tx.Rollback()

	// Lock sender row
	var senderBalance float64
	err = tx.QueryRow("SELECT balance FROM employee WHERE id = ? FOR UPDATE", balances.Id).Scan(&senderBalance)
	if err != nil {
		employeeBalance.Message = fmt.Sprintf("error locking sender: %s", err.Error())
		ch <- employeeBalance
		return
	}

	// Lock receiver row
	var receiverBalance float64
	err = tx.QueryRow("SELECT balance FROM employee WHERE id = ? FOR UPDATE", balances.ToId).Scan(&receiverBalance)
	if err != nil {
		employeeBalance.Message = fmt.Sprintf("error locking receiver: %s", err.Error())
		ch <- employeeBalance
		return
	}

	// Check balance
	if senderBalance < balances.Balance {
		employeeBalance.Id = balances.Id
		employeeBalance.Balance = senderBalance
		employeeBalance.Message = fmt.Sprintf("insufficient balance: %.2f, needed: %.2f", senderBalance, balances.Balance)
		ch <- employeeBalance
		return
	}

	// Update sender
	newSenderBalance := senderBalance - balances.Balance
	_, err = tx.Exec("UPDATE employee SET balance = ? WHERE id = ?", newSenderBalance, balances.Id)
	if err != nil {
		employeeBalance.Message = fmt.Sprintf("error updating sender: %s", err.Error())
		ch <- employeeBalance
		return
	}

	// Update receiver
	newReceiverBalance := receiverBalance + balances.Balance
	_, err = tx.Exec("UPDATE employee SET balance = ? WHERE id = ?", newReceiverBalance, balances.ToId)
	if err != nil {
		employeeBalance.Message = fmt.Sprintf("error updating receiver: %s", err.Error())
		ch <- employeeBalance
		return
	}

	// Commit
	if err := tx.Commit(); err != nil {
		employeeBalance.Message = fmt.Sprintf("error committing: %s", err.Error())
		ch <- employeeBalance
		return
	}

	emrepo, _ := employeerepo.NewEmployeeRepo()
	senderEmployee, _ := emrepo.GetEmployeeById(balances.Id)

	employeeBalance.Id = senderEmployee.Id
	employeeBalance.Name = senderEmployee.Name
	employeeBalance.Balance = newSenderBalance

	employeeBalance.Message = fmt.Sprint("Berhasil Kirim Balance ke ID : ", balances.ToId,
		" sebesar : ", balances.Balance,
		" nilai awal balance : ", senderBalance,
		" nilai akhir balance : ", newSenderBalance)

	ch <- employeeBalance
}
