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

	for i, row := range result.Rows {
		var colors []tablewriter.Colors
		for idx, rowItem := range row {
			if rowItem[:1] == "!"{
				result.Rows[i][idx] = rowItem[1:]
				colors = append(colors, tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor})
			} else {
				colors = append(colors, tablewriter.Colors{tablewriter.Normal, tablewriter.FgWhiteColor})
			}
		}
		table.Rich(row, colors)
	}
	//table.AppendBulk(result.Rows)
	table.Render()
}

func PrintRuler() {
	fmt.Println("==============================================================")
}


