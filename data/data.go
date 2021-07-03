package data

import (
	"fmt"
	"os"
)

var FileName bool


func WriteFile(){
	err := os.Mkdir("result", os.ModePerm)
	if err != nil {
		fmt.Println("Create result floor failed, because result is exist")
	}
}