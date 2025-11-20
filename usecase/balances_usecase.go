package usecase

import (
	"fmt"
	"mime/multipart"
	"router-template/entities"
	"router-template/entities/common"
	"router-template/repository/balancerepo"
	"strconv"
	"strings"
	"sync"

	"github.com/xuri/excelize/v2"
)

type BalancesUsecase interface {
	TopupBalance(Id int64, Balance float64) (employee entities.Employee, er error)
	SendBalance(balances []entities.SendBalance) ([]*entities.EmployeeBalance, error)
	SendBalanceExcel(file *multipart.FileHeader) (employees []*entities.EmployeeBalance, report entities.ReportSendBalance, er error)
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
	var mu sync.Mutex

	for _, balance := range balances {
		wg.Add(1)
		go emprepo.SendBalance(balance, ch, &wg)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for result := range ch {
		mu.Lock()
		// if result.Id != 0 {
		employee := &entities.EmployeeBalance{
			Id:      result.Id,
			Name:    result.Name,
			Balance: result.Balance,
			Message: result.Message,
		}
		employees = append(employees, employee)
		// }

		// glg.Debug("ini ", len(employees), ": dari", len(balances), "==>", result.Id)
		if len(employees) == len(balances) {
			break
		}
		mu.Unlock()
	}

	return employees, er
}

func (b *balancesUsecase) SendBalanceExcel(file *multipart.FileHeader) (employees []*entities.EmployeeBalance, report entities.ReportSendBalance, err error) {
	balances, err := parseExcelFile(file)
	if (err != nil) || (len(balances) == 0) {
		return employees, report, fmt.Errorf("failed parse excel file : %s", err.Error())
	}

	employees, err = b.SendBalance(balances)
	report.TotalData = len(employees)
	report.TotalNilaiTransaksi = nilaiTransaksi(employees)
	report.DataSukses, report.DataGagal = filterEmployeeBalance(employees)
	report.TotalSukses = len(report.DataSukses)
	report.TotalGagal = len(report.DataGagal)
	return employees, report, err
}

func filterEmployeeBalance(ebs []*entities.EmployeeBalance) (success []entities.EmployeeBalance, failed []entities.EmployeeBalance) {
	for _, eb := range ebs {
		if strings.Contains(eb.Message, "Berhasil") {
			success = append(success, entities.EmployeeBalance{
				Id:      eb.Id,
				Name:    eb.Name,
				Balance: eb.Balance,
				Message: eb.Message,
			})
		} else if !strings.Contains(eb.Message, "Berhasil") {
			failed = append(failed, entities.EmployeeBalance{
				Id:      eb.Id,
				Name:    eb.Name,
				Balance: eb.Balance,
				Message: eb.Message,
			})
		}
	}
	return
}

func nilaiTransaksi(ebs []*entities.EmployeeBalance) float64 {
	var total float64 = 0
	for _, eb := range ebs {
		total += eb.Balance
	}
	return total
}

func parseExcelFile(file *multipart.FileHeader) (balances []entities.SendBalance, err error) {
	fopen, err := file.Open()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fopen.Close()

	f, err := excelize.OpenReader(fopen)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}

	for index, row := range rows {
		if index <= 0 {
			continue
		}
		var balance entities.SendBalance

		balance.Id, err = strconv.ParseInt(row[1], 10, 64)
		balance.ToId, err = strconv.ParseInt(row[2], 10, 64)
		balance.Balance, err = common.ParseCurrency(row[3])
		balances = append(balances, balance)
	}
	return balances, err
}
