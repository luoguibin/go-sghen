package models

type Config struct {
	CodeMsgMap map[int]string
}

var (
	MConfig Config
)

var (
	RESP_CODE      = "code"
	RESP_MSG	   = "msg"
	RESP_DATA	   = "data"

	RESP_OK        = 1000
	RESP_ERR       = 1001
)

func init() {
	initCodeMsgMap()
}

func initCodeMsgMap() {
	MConfig.CodeMsgMap = make(map[int]string)
	MConfig.CodeMsgMap[RESP_OK] = "请求成功"
	MConfig.CodeMsgMap[RESP_ERR] = "失败,参数错误"
}