package controllers

import(
	"SghenApi/models"
	"strconv"
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
)

type WSServerController struct {
	BaseController
}

func init() {
	fmt.Println("WSServerController::init()");
	gameManager.Init()
}

var (
	upgrader 	= websocket.Upgrader{
		ReadBufferSize: 	2048,
		WriteBufferSize: 	2048,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	gameManager	GameManager
)

/**
 * WebSocket连接入口
 * 在BeforeRouter检测jwt中的合法后才给予长连接
 */
func (c *WSServerController) Get() {
	ws, err := upgrader.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil) 
	if err != nil { 
		models.MConfig.MLogger.Error("get ws error: " + err.Error())
	} 
	models.MConfig.MLogger.Debug("get ws: " + ws.RemoteAddr().String())

	uId,_ := strconv.ParseInt(c.Ctx.Input.Query("uId"), 10, 64)
	gameClient := &models.GameClient {
		ID:				uId,
		Conn:			ws,
		GameStatus:		models.StatusLogin,
	}
	gameManager.loginChan <- gameClient
}