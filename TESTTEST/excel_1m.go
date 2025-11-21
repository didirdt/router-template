//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"math/rand"

	"github.com/xuri/excelize/v2"
)

func main() {
	loop := 10000
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	f.SetCellValue("Sheet1", "A1", "No")
	f.SetCellValue("Sheet1", "B1", "ID sender")
	f.SetCellValue("Sheet1", "C1", "ID Receiver")
	f.SetCellValue("Sheet1", "D1", "Balance")

	noNumb := 1
	for index := range make([]struct{}, loop) {
		index++
		if index == 1 {
			continue
		}

		cellA := fmt.Sprintf("A%d", index)
		f.SetCellValue("Sheet1", cellA, noNumb)

		userId := []int{1, 2, 9, 12, 14, 15}
		randNumber := rand.Intn(5)
		userIdReceiver := append(userId[:randNumber], userId[randNumber+1:]...)
		randNumberTO := rand.Intn(5)
		randomFrom := userId[randNumber]
		randomTO := userIdReceiver[randNumberTO]

		cellB := fmt.Sprintf("B%d", index)
		f.SetCellValue("Sheet1", cellB, randomFrom)

		cellC := fmt.Sprintf("C%d", index)
		f.SetCellValue("Sheet1", cellC, randomTO)

		cellD := fmt.Sprintf("D%d", index)
		f.SetCellValue("Sheet1", cellD, "$10")
		noNumb++
	}

	if err := f.SaveAs("10kdata.xlsx"); err != nil {
		fmt.Println(err)
	}
}
