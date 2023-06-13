package game

import (
	"bytes"
	"strings"
)

// Note: Screen fits 40x19 characters.

var viewText = [][]byte{
	// viewTitle
	[]byte(`

 ____  ____  _____ _____ __________ _        
|  _ \|  _ \| ____|_   _|__  / ____| |  
| |_) | |_) |  _|   | |   / /|  _| | |  
|  __/|  _ <| |___  | |  / /_| |___| |_
|_|   |_| \_\_____| |_| /____|_____|___|

   _______   ______ ___   ___  _   _        
  |_   _\ \_/ / ___/ _ \ / _ \| \ | |       
    | |  \   / |  | | | | | | |  \| |       
    | |   | || |__| |_| | |_| | |\  |       
    |_|   |_| \____\___/ \___/|_| \_|


`),

	// viewIntro1
	[]byte(`GREETINGS!     WELCOME TO PRETZELSVILLE!

IN THIS GREAT SOCIETY, YOU ARE IN CHARGE
OF RUNNING YOUR OWN PRETZEL STAND. TO
GROW YOUR PRETZEL EMPIRE, YOU WILL NEED
TO MAKE THESE DECISIONS EVERY DAY:


1. WHO TO USE FOR YOUR OWN GAIN

2. HOW FAR TO BEND OR BREAK THE RULES

3. WHAT PRICE TO CHARGE FOR EACH PRETZEL


YOU WILL BEGIN WITH $4.20 CASH (ASSETS).

` + string(bytes.TrimRight(centeredText("PRESS SPACE TO CONTINUE..."), "\n"))),

	// viewFinancialReport
	[]byte(` $$  PRETZELSVILLE FINANCIAL REPORT  $$

                 DAY %d

%7d PRETZELS SOLD     INCOME

%7s PER PRETZEL       %s


%7d PRETZELS MADE     EXPENSES

%7d SIGNS MADE        %s

              PROFIT %s

              ASSETS %s

` + string(bytes.TrimRight(centeredText("PRESS SPACE TO CONTINUE..."), "\n"))),
}

func centeredText(text string) []byte {
	if len(text) > 40 {
		return []byte(text + "\n")
	}
	spaces := (40 - len(text)) / 2
	return []byte(strings.Repeat(" ", spaces) + text + strings.Repeat(" ", 40-len(text)-spaces) + "\n")
}
