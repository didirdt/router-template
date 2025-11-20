package usecase

import (
	"router-template/entities"
	"router-template/repository/balancerepo"
	"sync"
)

type BalancesUsecase interface {
	TopupBalance(Id int64, Balance float64) (employee entities.Employee, er error)
	SendBalance(balances []entities.SendBalance) ([]*entities.EmployeeBalance, error)
}

func NewBalancesUsecase() BalancesUsecase {
	return &balancesUsecase{}
}

type balancesUsecase struct{}

func (b *balancesUsecase) TopupBalance(id int64, balance float64) (employee entities.Employee, er error) {
	repoBalance, _ := balancerepo.NewBalanceRepo()
	employee, er = repoBalance.TopupBalance(id, balance)

	return
}

func (b *balancesUsecase) SendBalance(balances []entities.SendBalance) (employees []*entities.EmployeeBalance, er error) {
	ch := make(chan *entities.EmployeeBalance, len(balances))
	emprepo, _ := balancerepo.NewBalanceRepo()
	var wg sync.WaitGroup
	// var mu sync.Mutex

	for _, balance := range balances {
		wg.Add(1)
		go emprepo.SendBalance(balance, ch, &wg)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for result := range ch {
		if result.Id != 0 {
			// mu.Lock()
			employee := &entities.EmployeeBalance{
				Id:      result.Id,
				Name:    result.Name,
				Balance: result.Balance,
				Message: result.Message,
			}
			employees = append(employees, employee)
			// mu.Unlock()
		}

		if len(employees) == len(balances) {
			break
		}
	}

	return employees, er
}
