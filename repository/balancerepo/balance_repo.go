package balancerepo

import (
	"router-template/entities"
)

type BalanceRepo interface {
	TopupBalance(id int64, balance float64) (entities.Employee, error)
	SendBalance(balances entities.SendBalance, ch chan entities.EmployeeBalance) (entities.EmployeeBalance, error)
}
