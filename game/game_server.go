package game

import(
	"SghenApi/models"
	"SghenApi/helper"
	"net/http"
	"strconv"
	"time"
	"fmt"
	"runtime"
	"github.com/gorilla/websocket"
	"github.com/astaxie/beego/context"
)

var (
	upgrader 		= websocket.Upgrader{
						ReadBufferSize: 	2048,
						WriteBufferSize: 	2048,
						CheckOrigin: func(r *http.Request) bool {
							return true
						},
					}
	MGameServer		= &GameServer{}

	GameServerStatus	=	1
)


type GameServer struct {
	GameMapService			*GameMapService
	GameAuthorityService	*GameAuthorityService
}

/*
 * init MGameServer
 */
func init() {
	MGameServer.Start()

	go func () {
		for {
			var str string
			fmt.Scan(&str)
	
			if (str == "game:finish") {
				MGameServer.GameMapService.LogoutAll()
				GameServerStatus = -1
			} else if (str == "game:start") {
				GameServerStatus = 1
			} else if (str == "game:getUser") {
				fmt.Print("input id:")
				fmt.Scan(&str)
				id, err := strconv.ParseInt(str, 10, 64)
				if err == nil {
					client := MGameServer.GameMapService.GetUserData(id)
					if client != nil {
						fmt.Println(client.GameData)
					} else {
						fmt.Println("query result is nil")
					}
				} else {
					fmt.Println(err)
				}
			} else if (str == "game:threadcount") {
				fmt.Printf("	threadcount=%d\n", runtime.NumGoroutine())
			} else if (str == "exit") {
				break
			}
		}
	} ()
}

/*
 * game server start
 */
func (gameServer *GameServer) Start() {
	fmt.Println("GameServer::Start()");

	MGameServer.GameAuthorityService = &GameAuthorityService{}
	go MGameServer.GameAuthorityService.Start()

	MGameServer.GameMapService = &GameMapService{}
	MGameServer.GameMapService.Start()
}

/*
 * add the context to the server
 */
func AddToServer(Ctx *context.Context, uId int64) {
	ws, err := upgrader.Upgrade(Ctx.ResponseWriter, Ctx.Request, nil) 
	if err != nil { 
		models.MConfig.MLogger.Error("get ws error:\n%s", err)
	} 
	models.MConfig.MLogger.Debug("get ws: " + ws.RemoteAddr().String())

	gameClient := &GameClient {
		ID:				uId,
		Conn:			ws,
		GameStatus:		GStatusLogin,
	}
	MGameServer.GameAuthorityService.LoginChan <- gameClient
}

/*
 * the client order handler
 */
func GoGameClientHandle(gameClient *GameClient) {
	preTime := time.Now().UnixNano() / 1e6
	for {
		// 获取指令 
		var order GameOrder
		err := gameClient.Conn.ReadJSON(&order)
		if err != nil {
			models.MConfig.MLogger.Error("ws read msg error: " + err.Error())
			if gameClient.GameStatus != GStatusLogout {
				gameClient.GameStatus = GStatusLogout
				MGameServer.GameAuthorityService.LoginChan <- gameClient
			}
			return
		}
		
		curTime := time.Now().UnixNano() / 1e6
		// fmt.Printf("%v  %d\n", order, curTime - preTime)
		if (curTime - preTime < 300) {
			continue
		}
		preTime = curTime
		if (order.OrderType < 100) {
			models.MConfig.MLogger.Error("ws read msg error: order.OrderType < 1000")
			continue
		}
		v := order.OrderType / 10000 * 10000
		switch v {
			case OT_Msg:
				MGameServer.GameMapService.DealOrderMsg(gameClient, &order)
			case OT_Skill:
				MGameServer.GameMapService.DealOrderSkill(gameClient, &order)
			case OT_Action:
				MGameServer.GameMapService.DealOrderAction(gameClient, &order)
			default:
				models.MConfig.MLogger.Error(string(gameClient.ID) + " order invalid: " + string(order.OrderType))
		}
	} 
}

/*
 * reset the game data when login successfully
 */
func ResetGameData(gameData *models.GameData) {
	gameData.BloodAll = gameData.BloodBase + 300000
	gameData.X0 = gameData.X
	gameData.Y0 = gameData.Y
	gameData.X1 = gameData.X
	gameData.Y1 = gameData.Y
	gameData.Speed = gameData.SpeedBase
	gameData.MoveTime = 0
	gameData.EndTime = 0
	gameData.Move = 0
}

/*
 * reset the client move data
 */
func ResetGameDataMove(gameData *models.GameData, orderAction *GameOrderAction) {
	curTime := helper.GetMillisecond()
	if gameData.Move == 1 {
		if gameData.EndTime < curTime {
			gameData.X = gameData.X1
			gameData.Y = gameData.Y1
			gameData.Move = 0
		} else {
			stayTime := gameData.EndTime - curTime
			moveRatio := float64(1) - float64(stayTime) / float64(gameData.MoveTime)
			gameData.X = int(float64(gameData.X1 - gameData.X0) * moveRatio) + gameData.X0
			gameData.Y = int(float64(gameData.Y1 - gameData.Y0) * moveRatio) + gameData.Y0
		}
	}
	if orderAction != nil {
		gameData.X1 = orderAction.X
		gameData.Y1 = orderAction.Y
		gameData.X0 = gameData.X
		gameData.Y0 = gameData.Y
		distance := helper.GClientDistance(gameData.X, gameData.Y, orderAction.X, orderAction.Y)
		moveTime := distance / float64(gameData.Speed)
		gameData.MoveTime = int64(moveTime * 1000)
		gameData.EndTime = curTime + gameData.MoveTime
		gameData.Move = 1
		// fmt.Printf("distance=%v  moveTime=%v\n", distance, gameData.GMoveTime)
	}
	
}