package balancerepo

import (
	"router-template/entities"
	"sync"
)

type BalanceRepo interface {
	TopupBalance(id int64, balance float64) (entities.Employee, error)
	SendBalance(balances entities.SendBalance, ch chan *entities.EmployeeBalance, wg *sync.WaitGroup)
}
