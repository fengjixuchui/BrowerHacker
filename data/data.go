package data

import "os"

const SQLiteName = `sqlite3`  // 定义数据库连接类型

var Env = os.Getenv("USERPROFILE") // 定义用户的环境变量

var LocalDB = `LocalDB`
