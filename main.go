package main

import (
	"BrowserHacker/core"
	"BrowserHacker/data"
	"flag"
)

var(
	FirefoxProfile string
	Firefoxdb string
	ChromeProfile string
	ChromeState string
)
func init(){
	flag.BoolVar(&data.FileName,"o",false,"return result file")
	flag.StringVar(&Firefoxdb,"k","","firefox key4.db profile")
	flag.StringVar(&FirefoxProfile,"f","","firefox login.json profile")
	flag.StringVar(&ChromeProfile,"c","","Chrome login data profile")
	flag.StringVar(&ChromeState,"s","","Chrome login state profile")
	flag.Parse()
	if Firefoxdb != ""{
		core.FireFoxMainDB = Firefoxdb
	}
	if FirefoxProfile !=""{
		core.FireFoxMainProfile = FirefoxProfile
	}
	if ChromeProfile != ""{
		core.ChromeMainLogin = ChromeProfile
	}
	if ChromeState != ""{
		core.ChromeMainState = ChromeState
	}
}


func main() {
	if data.FileName == true {
		data.WriteFile()
	}
	_ = core.FireFoxExec()
	core.ChromeSQLiteV80()
	core.ChromeSQLite()
}
