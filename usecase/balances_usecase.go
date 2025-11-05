package usecase

import (
	"router-template/entities"
	"router-template/repository/balancerepo"
)

type BalancesUsecase interface {
	TopupBalance(Id int64, Balance float64) (employee entities.Employee, er error)
	SendBalance(balances []entities.SendBalance) ([]entities.EmployeeBalance, error)
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

func (b *balancesUsecase) SendBalance(balances []entities.SendBalance) (employees []entities.EmployeeBalance, er error) {
	for _, balance := range balances {
		emprepo, _ := balancerepo.NewBalanceRepo()
		employee, er := emprepo.SendBalance(balance)
		if er != nil {
			employee.Message = "error while Process payment send Balance : " + er.Error()
		}
		employees = append(employees, employee)
	}
	return employees, er
}
