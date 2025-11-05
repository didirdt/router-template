package employeerepo

import (
	"database/sql"
	"errors"
	"fmt"
	"router-template/entities"
	"router-template/entities/app"
	"router-template/repository/built_in/databasefactory"

	"github.com/nyaruka/phonenumbers"
)

func newEmployeeMysqlImpl() EmployeeRepo {
	conn := databasefactory.AppDb.GetConnection()
	return &employeeMysqlImpl{conn: conn.(*sql.DB)}
}

type employeeMysqlImpl struct {
	conn *sql.DB
}

func (e *employeeMysqlImpl) GetEmployee() (list []entities.Employee, er error) {
	rows, er := e.conn.Query("Select id, name, address, phone_number, balance from employee")
	if er != nil {
		return list, er
	}

	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var item entities.Employee
		if er = rows.Scan(&item.Id, &item.Name, &item.Address, &item.PhoneNumber, &item.Balance); er != nil {
			return list, er
		}

		list = append(list, item)
	}

	if len(list) == 0 {
		return list, app.ErrNoRecord
	} else {
		return
	}
}

func (e *employeeMysqlImpl) GetEmployeeById(id int64) (employee entities.Employee, er error) {
	row := e.conn.QueryRow(`SELECT 
		id,
		name,
		address,
		phone_number,
		balance
		FROM employee WHERE id=?`, id)

	if er = row.Scan(&employee.Id, &employee.Name, &employee.Address, &employee.PhoneNumber, &employee.Balance); er != nil {
		if er == sql.ErrNoRows {
			return
		} else {
			return employee, errors.New(fmt.Sprint("error while get employee : ", er.Error()))
		}
	}

	return
}

func (e *employeeMysqlImpl) CreateEmployee(name, address, phone_number string) (employee entities.Employee, er error) {
	phone_numbers, er := phonenumbers.Parse(phone_number, "ID")
	formatPhoneNumber := phonenumbers.Format(phone_numbers, phonenumbers.NATIONAL)
	if er != nil {
		return employee, errors.New(fmt.Sprint("error while create employee : ", er.Error()))
	}

	result, er := e.conn.Exec(`INSERT INTO employee
		(name, address, phone_number) VALUES (?, ?, ?)`,
		name, address, formatPhoneNumber)
	if er != nil {
		return employee, errors.New(fmt.Sprint("error while create employee : ", er.Error()))
	}

	lastInsertId, er := result.LastInsertId()
	if er != nil {
		return employee, errors.New(fmt.Sprint("error while get last insert id employee : ", er.Error()))
	}

	employee, er = e.GetEmployeeById(lastInsertId)
	return
}

func (e *employeeMysqlImpl) UpdateEmployee(id int64, name, address, phone_number string) (employee entities.Employee, er error) {
	phone_numbers, er := phonenumbers.Parse(phone_number, "ID")
	formatPhoneNumber := phonenumbers.Format(phone_numbers, phonenumbers.NATIONAL)
	if er != nil {
		return employee, errors.New(fmt.Sprint("error while create employee : ", er.Error()))
	}

	result, er := e.conn.Exec(`UPDATE employee SET name=?, address=?, phone_number=? WHERE id=?`,
		name, address, formatPhoneNumber, id)
	if result == nil || er != nil {
		return employee, errors.New(fmt.Sprint("error while update employee : ", er.Error()))
	}

	employee, er = e.GetEmployeeById(id)
	return
}

func (e *employeeMysqlImpl) DeleteEmployee(id int64) (employee entities.Employee, er error) {
	employee, er = e.GetEmployeeById(id)
	if er != nil {
		return employee, errors.New(fmt.Sprint("error while get employee : ", er.Error()))
	}

	result, er := e.conn.Exec(`DELETE FROM employee WHERE id=?`, id)
	if result == nil || er != nil {
		return employee, errors.New(fmt.Sprint("error while delete employee : ", er.Error()))
	}

	return
}
