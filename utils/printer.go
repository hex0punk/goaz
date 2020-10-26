package printer

import (
	"github.com/hex0punk/color"
	"github.com/olekukonko/tablewriter"
	"os"
	"fmt"
)

type ResultTable struct {
	Rows			[][]string
	Columns			[]string
}

var	(
	Danger = color.New(color.FgRed).PrintfFunc()
	Warning = color.New(color.FgYellow).PrintfFunc()
	Info = color.New(color.FgCyan).PrintfFunc()
	Data = color.New(color.FgWhite).PrintfFunc()
	InfoHeading = color.New(color.FgGreen).Add(color.Underline).PrintfFunc()
)

func PrintTable(result *ResultTable){
	fmt.Println()
	PrintRuler()
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(result.Columns)
	table.SetBorder(false)

	var headerColors []tablewriter.Colors
	for i, _ := range result.Columns{
		if i == 0{
			headerColors = append(headerColors, tablewriter.Colors{tablewriter.Bold, tablewriter.BgHiBlueColor, tablewriter.FgHiWhiteColor})
		} else {
			headerColors = append(headerColors, tablewriter.Colors{tablewriter.Bold, tablewriter.BgHiBlueColor, tablewriter.FgHiWhiteColor})
		}
	}
	table.SetHeaderColor(headerColors...)

	var rowColors []tablewriter.Colors
	for i, _ := range result.Columns{
		if i == 0{
			rowColors = append(rowColors, tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlueColor})
		} else {
			rowColors = append(rowColors, tablewriter.Colors{tablewriter.Bold, tablewriter.FgGreenColor})
		}
	}
	table.SetColumnColor(rowColors...)

	//table.SetColumnColor(tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlackColor},
	//	tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor})


	table.AppendBulk(result.Rows)
	table.Render()
}

func PrintRuler() {
	fmt.Println("==============================================================")
}


