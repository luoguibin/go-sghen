package models

import (
	"os"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/config"
	"github.com/jinzhu/gorm"
)

type Config struct {
	SGHENENV	string

	JwtSecretKey string

	PathTypeMap	map[string]string
	MaxUploadSize	int
	MaxUploadCount  int


	dBHost     	string
	dBName     	string
	dBUsername 	string
	dBPassword 	string
	dBMaxIdle  int
	dBMaxConn  int

	CodeMsgMap 	map[int]string

	MLogger *logs.BeeLogger
}

var (
	MConfig Config

	dbOrmDefault   *gorm.DB


)

var (
	STR_CODE    	= 	"code"
	STR_MSG	    	= 	"msg"
	STR_DATA		= 	"data"
	STR_TOKEN		= 	"token"

	CODE_OK  		= 	1000
	CODE_ERR    	= 	1001
	CODE_ERR_TOKEN	=	1002
)

func init() {
	initConfParams()
	initCodeMsgMap()
	initPathTypeMap()
	initLog()
}

func initConfParams() {
	SGHENENV := os.Getenv("SGHENENV")
	if len(SGHENENV) <= 0 {
		SGHENENV = "dev"
	}
	appConf, err := config.NewConfig("ini", "conf/app.conf")
	if err != nil {
		return
	}

	MConfig = Config{}
	if appConf != nil {
		MConfig.SGHENENV = SGHENENV
		MConfig.JwtSecretKey = appConf.String(SGHENENV + "::jwtSecretKey")
		MConfig.MaxUploadSize, _ = appConf.Int(SGHENENV + "::maxUploadSize")
		MConfig.MaxUploadCount, _ = appConf.Int(SGHENENV + "::maxUploadCount")

		MConfig.dBHost = appConf.String(SGHENENV + "::dbHost")
		MConfig.dBName = appConf.String(SGHENENV + "::dbName")
		MConfig.dBUsername = appConf.String(SGHENENV + "::dbUsername")
		MConfig.dBPassword = appConf.String(SGHENENV + "::dbPassword")
		MConfig.dBMaxIdle, _ = appConf.Int(SGHENENV + "::dbMaxIdle")
		MConfig.dBMaxConn, _ = appConf.Int(SGHENENV + "::dbMaxConn")
	}
	
}

func initCodeMsgMap() {
	MConfig.CodeMsgMap = make(map[int]string)
	MConfig.CodeMsgMap[CODE_OK] = "请求成功"
	MConfig.CodeMsgMap[CODE_ERR] = "失败,参数错误"
}

func initPathTypeMap() {
	MConfig.PathTypeMap = make(map[string]string, 0)
	MConfig.PathTypeMap["peotry"] = "./file/peotry/img/"
	MConfig.PathTypeMap["normal"] = "./file/normal/"
}

func initLog() {
	MLogger := logs.NewLogger(10000) 
	MLogger.SetLogger(logs.AdapterFile,`{"filename":"./logs/temp.log","maxlines":1000000,"maxsize":256,"perm": "0644"}`)
	MLogger.EnableFuncCallDepth(true)
	MLogger.Async()
	MLogger.SetLevel(logs.LevelDebug)
}

func NewLog() *logs.BeeLogger{
	logger := logs.NewLogger(10000) 
	logger.SetLogger(logs.AdapterFile,`{"filename":"./logs/temp0.log","maxlines":1000000,"maxsize":256,"perm": "0644"}`)
	logger.EnableFuncCallDepth(true)
	logger.Async()
	logger.SetLevel(logs.LevelDebug)

	return logger
}


func GetDb() *gorm.DB {
	return dbOrmDefault
}