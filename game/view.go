package game

import (
	"bytes"
	"strings"
)

// Note: The screen size is 40x22.

var viewFinancialReport = ` $$  PRETZELSVILLE FINANCIAL REPORT  $$


` + string(centeredText("DAY 1")) + `

%7d PRETZELS SOLD     INCOME

%7s PER PRETZEL       %s


%7d PRETZELS MADE     EXPENSES

%7d SIGNS MADE        %s


              PROFIT %s

              ASSETS %s


` + string(bytes.TrimRight(centeredText("PRESS SPACE TO CONTINUE..."), "\n"))

func centeredText(text string) []byte {
	if len(text) > 40 {
		return []byte(text + "\n")
	}
	spaces := (40 - len(text)) / 2
	return []byte(strings.Repeat(" ", spaces) + text + strings.Repeat(" ", 40-len(text)-spaces) + "\n")
}
