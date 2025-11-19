package usecase

import (
	"router-template/entities"
	"router-template/repository/balancerepo"
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
	ch := make(chan entities.EmployeeBalance, len(balances))
	for _, balance := range balances {
		emprepo, _ := balancerepo.NewBalanceRepo()
		go emprepo.SendBalance(balance, ch)

		// if er != nil {
		// 	employee.Message = "error while Process payment send Balance : " + er.Error()
		// }

		// employees = append(employees, employee)
	}

	// for em := range ch {
	// 	employees = append(employees, em)
	// }

	for {
		em := <-ch
		if em.Id != 0 {
			employees = append(employees, &em)
		}

		if len(employees) == len(balances) {
			break
		}
	}

	close(ch)
	return employees, er
}
