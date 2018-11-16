package controllers

import(
	"SghenApi/models"
	"strconv"
	"time"
	"fmt"
	"math/rand"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/astaxie/beego/logs"
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

	clients   	= make(map[*models.WsUser]bool)

	wsLogger *logs.BeeLogger
)

/**
 * WebSocket连接入口
 */
func (c *WSServerController) Get() {
	ws, err := upgrader.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil) 
	if err != nil { 
		wsLogger.Error("get ws error: " + err.Error())
	} 
	wsLogger.Debug("get ws: " + ws.RemoteAddr().String())

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

	clients[&wsUser] = true
	go func(wsuser models.WsUser) {
		for { 
			var action models.GameAction0
			err := wsuser.Conn.ReadJSON(&action)
			if err == nil {
				for client := range clients {
					if client.ID == action.Target {
						// fmt.Println(action)

						switch action.Action {
							case "fist":
								ran := rand.Intn(200)
								if rand.Intn(10) < 5 {
									ran = wsuser.WsData.GPower + ran
								} else {
									ran = wsuser.WsData.GPower - ran
								}
								client.WsData.GBlood -= ran

								if client.WsData.GBlood < 0 {
									client.WsData.GBlood = 0
								}
								break;
							case "skill":
								ran := rand.Intn(200)
								ran = int(float32(wsuser.WsData.GPower) * 1.3) + ran
								client.WsData.GBlood -= ran

								if client.WsData.GBlood < 0 {
									client.WsData.GBlood = 0
								}
								break;
							case "skill_big":
								ran := rand.Intn(10000)
								ran = int(float32(wsuser.WsData.GPower) * 3.3) + ran
								client.WsData.GBlood -= ran

								if client.WsData.GBlood < 0 {
									client.WsData.GBlood = 0
								}
								break;
							case "msg":
								action.Target = wsuser.ID
								client.Conn.WriteJSON(action)
								break
						}
						break;
					}
				}
			} 
			
			time.Sleep(time.Second * 2)
		} 
	}(wsUser)
}


func dataCenter() {
	wsLogger = models.NewLog()
	
	for {
		// ①读取ws数据
		

		// ②计算
		wsDatas := make([]interface{}, 0)
		for client := range clients { 
			wsDatas = append(wsDatas, client.WsData)
		} 

		// ③发送ws数据
		for client := range clients { 
			err := client.Conn.WriteJSON(wsDatas) 
			if err != nil { 
				wsLogger.Debug(err.Error())
				client.Conn.Close() 
				delete(clients, client) 
			} 
		} 

		time.Sleep(time.Second * 1)
		// fmt.Println("clients count=", len(clients))
	}
}
