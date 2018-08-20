package models

import (
	"github.com/astaxie/beego/logs"
)

type Config struct {
	CodeMsgMap map[int]string
	MLogger *logs.BeeLogger

}

var (
	MConfig Config
	JWT_SECRET_KEY = "SghenMorge"

	IMAGE_SAVE_PATH = "./upload/peotry/images/"
)

var (
	RESP_CODE      = "code"
	RESP_MSG	   = "msg"
	RESP_DATA	   = "data"
	RESP_TOKEN	   = "token"

	RESP_OK        = 1000
	RESP_ERR       = 1001
)

func init() {
	initCodeMsgMap()
	initLog()
}

func initCodeMsgMap() {
	MConfig.CodeMsgMap = make(map[int]string)
	MConfig.CodeMsgMap[RESP_OK] = "请求成功"
	MConfig.CodeMsgMap[RESP_ERR] = "失败,参数错误"
}

func initLog() {
	MLogger := logs.NewLogger(10000) 
	MLogger.SetLogger(logs.AdapterFile,`{"filename":"./logs/temp.log","maxlines":1000000,"maxsize":256}`)
	MLogger.EnableFuncCallDepth(true)
	MLogger.SetLevel(logs.LevelDebug)     // 设置日志写入缓冲区的等级
	
	// log test
	// log.Emergency("Emergency")
    // log.Alert("Alert")
    // log.Critical("Critical")
    // log.Error("Error")
    // log.Warning("Warning")
    // log.Notice("Notice")
    // log.Informational("Informational")
    // log.Debug("Debug")

	// log.Flush() // 将日志从缓冲区读出，写入到文件
    // log.Close()
}
