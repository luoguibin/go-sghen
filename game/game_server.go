package game

import (
	"fmt"
	"go-sghen/helper"
	"go-sghen/models"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"

	"github.com/astaxie/beego/context"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  2048,
		WriteBufferSize: 2048,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	MGameServer = &GameServer{}

	GameServerStatus = 1
)

type GameServer struct {
	GameMapService       *GameMapService
	GameAuthorityService *GameAuthorityService
}

/*
 * init MGameServer
 */
func init() {
	MGameServer.Start()

	// listen the program died
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		fmt.Println("programer exit")

		MGameServer.GameMapService.LogoutAll()
		GameServerStatus = -1

		os.Exit(0)
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			var str string
			fmt.Scan(&str)

			if str == "game:finish" {
				MGameServer.GameMapService.LogoutAll()
				GameServerStatus = -1
			} else if str == "game:start" {
				GameServerStatus = 1
			} else if str == "game:getUser" {
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
			} else if str == "game:threadcount" {
				fmt.Printf("	threadcount=%d\n", runtime.NumGoroutine())
			} else if str == "exit" {
				break
			}
		}
	}()
}

/*
 * game server start
 */
func (gameServer *GameServer) Start() {
	fmt.Println("GameServer::Start()")

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

	gameClient := &GameClient{
		ID:         uId,
		Conn:       ws,
		GameStatus: GStatusLogin,
	}
	MGameServer.GameAuthorityService.LoginChan <- gameClient
}

/*
 * the client order handler
 */
func GoGameClientHandle(gameClient *GameClient) {
	preTime := helper.GetMillisecond()
	for {
		// 获取指令
		var order models.GameOrder
		err := gameClient.Conn.ReadJSON(&order)
		if err != nil {
			models.MConfig.MLogger.Error("ws read msg error: " + err.Error())
			if gameClient.GameStatus != GStatusLogout {
				gameClient.GameStatus = GStatusLogout
				MGameServer.GameAuthorityService.LoginChan <- gameClient
			}
			return
		}

		curTime := helper.GetMillisecond()
		// fmt.Printf("%v  %d\n", order, curTime - preTime)
		if curTime-preTime < 300 {
			continue
		}
		preTime = curTime
		if order.OrderType < 100 {
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
	gameData.Speed = gameData.SpeedBase
}

/*
 * reset the client move data
 */
func ResetGameDataMove(gameData *models.GameData, orderAction []models.GameOrderAction) {
	gameData.MoveOrder = orderAction
}
