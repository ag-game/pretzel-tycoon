package game

import (
	"bytes"
	"strings"
)

// Note: Screen fits 40x18 characters.

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

` + string(bytes.TrimRight(centeredText("PRESS ENTER TO CONTINUE..."), "\n"))),

	// viewStartDayProduction1
	[]byte(`                 DAY %d

STARTING SUPPLIES

PRODUCTION AMOUNTS

HOW MANY BATCHES (DOZENS) OF PRETZELS
DO YOU WISH TO MAKE ?`),

	// viewStartDayProduction2
	[]byte(`                 DAY %d

STARTING SUPPLIES

PRODUCTION AMOUNTS

HOW MANY BATCHES (DOZENS) OF PRETZELS
DO YOU WISH TO MAKE ?%d

HOW MANY ADVERTISING SIGNS ($1.25 EACH)
DO YOU WISH TO MAKE ?`),

	// viewStartDayProduction3
	[]byte(`                 DAY %d

STARTING SUPPLIES

PRODUCTION AMOUNTS

HOW MANY BATCHES (DOZENS) OF PRETZELS
DO YOU WISH TO MAKE ?%d

HOW MANY ADVERTISING SIGNS ($1.25 EACH)
DO YOU WISH TO MAKE ?%d

WHAT PRICE (IN CENTS) DO YOU WISH TO
CHARGE FOR PRETZELS ?`),

	// viewDay
	nil,

	// viewFinancialReport
	[]byte(` $$  PRETZELSVILLE FINANCIAL REPORT  $$

                 DAY %d

%7d PRETZELS SOLD     INCOME

%7s PER PRETZEL       %s


%7d PRETZELS MADE     EXPENSES

%7d SIGNS MADE        %s

              PROFIT %s

              ASSETS %s

` + string(bytes.TrimRight(centeredText("PRESS ENTER TO CONTINUE..."), "\n"))),
}

func centeredText(text string) []byte {
	if len(text) > 40 {
		return []byte(text + "\n")
	}
	spaces := (40 - len(text)) / 2
	return []byte(strings.Repeat(" ", spaces) + text + strings.Repeat(" ", 40-len(text)-spaces) + "\n")
}
