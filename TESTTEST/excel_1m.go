//go:build ignore
// +build ignore

package main

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

func main() {
	loop := 1000000
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

	for index := range make([]struct{}, loop) {
		index++
		if index == 1 {
			continue
		}

		cell := fmt.Sprintf("A%d", index)
		f.SetCellValue("Sheet1", cell, index+1)

		// cell := fmt.Sprintf("B%d", index)
		// f.SetCellValue("Sheet1", cell, randInt(1, 10))
	}

	// f.SetActiveSheet(index)

	if err := f.SaveAs("Book1.xlsx"); err != nil {
		fmt.Println(err)
	}
}
