package printer

import "github.com/DharmaOfCode/color"


var	(
	Danger = color.New(color.FgRed).PrintfFunc()
	Warning = color.New(color.FgYellow).PrintfFunc()
	Info = color.New(color.FgCyan).PrintfFunc()
	Data = color.New(color.FgWhite).PrintfFunc()

	InfoHeading = color.New(color.FgGreen).Add(color.Underline).PrintfFunc()
)


