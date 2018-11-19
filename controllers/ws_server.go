package controllers

import(
	"SghenApi/models"
	"strconv"
	"time"
	"fmt"
	"sync"
	"math/rand"
	"net/http"
	"github.com/gorilla/websocket"
)

type WSServerController struct {
	BaseController
}

func init() {
	fmt.Println("WSServerController::init()");
    go dataCenter()
}

var (
	upgrader 	= websocket.Upgrader{
		ReadBufferSize: 	2048,
		WriteBufferSize: 	2048,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	clients  sync.Map
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
	uGame := models.Game0 {
		ID: 		uId,
		GName: 		ws.RemoteAddr().String(),
		GBlood: 	300000,
		GBloodAll: 	353535,
		GLevel:		103,
		GPower:		5000,
	}
	wsUser := models.WsUser {
		ID:		uId,
		Conn:	ws,
		WsData: uGame,	
	}

	clients.Store(uId, wsUser)
	// go func(wsUser models.WsUser) {
		for { 
			var action models.GameAction0
			err := wsUser.Conn.ReadJSON(&action)
			if err != nil {
				models.MConfig.MLogger.Error("ws read msg error: " + err.Error())
				break;
			}
			client_, ok := clients.Load(action.Target)
			if !ok {
				models.MConfig.MLogger.Error("clients.Load error")
				break;
			}
			
			client, ok := (client_).(*models.WsUser)
			if !ok {
				models.MConfig.MLogger.Error("clients cast error")
				break;
			}

			switch action.Action {
				case "fist":
					ran := rand.Intn(200)
					if rand.Intn(10) < 5 {
						ran = wsUser.WsData.GPower + ran
					} else {
						ran = wsUser.WsData.GPower - ran
					}
					client.WsData.GBlood -= ran

					if client.WsData.GBlood < 0 {
						client.WsData.GBlood = 0
					}
				case "skill":
					ran := rand.Intn(200)
					ran = int(float32(wsUser.WsData.GPower) * 1.3) + ran
					client.WsData.GBlood -= ran

					if client.WsData.GBlood < 0 {
						client.WsData.GBlood = 0
					}
				case "skill_big":
					ran := rand.Intn(10000)
					ran = int(float32(wsUser.WsData.GPower) * 3.3) + ran
					client.WsData.GBlood -= ran

					if client.WsData.GBlood < 0 {
						client.WsData.GBlood = 0
					}
				case "msg":
					action.Target = wsUser.ID
					client.Conn.WriteJSON(action)
				default:
			}
		} 
	// }(wsUser)
}


func dataCenter() {	
	for {
		// ①读取ws数据
		

		// ②计算
		wsDatas := make([]interface{}, 0)
		clients.Range(func(key, client_ interface{}) bool {
			client, ok := (client_).(*models.WsUser)
			if !ok {
				models.MConfig.MLogger.Error("dataCenter() clients cast error")
				return true
			}
			wsDatas = append(wsDatas, client.WsData)
			return true
		})

		// ③发送ws数据
		count := 0
		clients.Range(func(key, client_ interface{}) bool {
			client, ok := (client_).(*models.WsUser)
			count++
			if !ok {
				models.MConfig.MLogger.Error("dataCenter() clients cast error")
				return true
			}
			err := client.Conn.WriteJSON(wsDatas) 
			if err != nil { 
				models.MConfig.MLogger.Debug(err.Error())
				client.Conn.Close() 
				clients.Delete(key)
			} 
			return true
		})

		time.Sleep(time.Second * 1)
		// fmt.Println("clients count=", count)
	}
}
