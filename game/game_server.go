package game

import(
	"SghenApi/models"
	"net/http"
	"strconv"
	"sync"
	"time"
	"fmt"
	"runtime"
	"github.com/gorilla/websocket"
	"github.com/astaxie/beego/context"
	"github.com/goinggo/mapstructure"
)

var (
	upgrader 		= websocket.Upgrader{
						ReadBufferSize: 	2048,
						WriteBufferSize: 	2048,
						CheckOrigin: func(r *http.Request) bool {
							return true
						},
					}

	gServerStatus	=	1

	gLoginService 	GLoginService
	gMapMap 		sync.Map	
)


type GameServer struct {
	
}


func init() {
	fmt.Println("GameServer::init()");

	gMapService0 := &GMapService{}
	gMapService0.Init("HumanWorld")
	gMapMap.Store(0, gMapService0)

	// gMapService1 := &GMapService{}
	// gMapService1.Init("Hell")
	// gMapMap.Store(1, gMapService1)

	// gMapService2 := &GMapService{}
	// gMapService2.Init("Heaven")
	// gMapMap.Store(2, gMapService2)

	go gLoginService.start()
	go goDataCenter()
	go goGameCommond()
}

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
	GLoginChan <- gameClient
}

func goGameCommond() {
	for {
		var str string
		fmt.Scan(&str)

		if (str == "game:finish") {
			logoutAll()
		} else if (str == "game:start") {
			gServerStatus = 1
		} else if (str == "game:getUser") {
			fmt.Print("input id:")
			fmt.Scan(&str)
			id, err := strconv.ParseInt(str, 10, 64)
			if err == nil {
				fmt.Println(getUserData(id))
			} else {
				fmt.Println(err)
			}
		} else if (str == "game:threadcount") {
			fmt.Printf("	threadcount=%d\n", runtime.NumGoroutine())
		}
	}
}

func goGameClientHandle(gameClient *GameClient) {
	preTime := time.Now().UnixNano() / 1e6
	for {
		// 获取指令 
		var order GameOrder
		err := gameClient.Conn.ReadJSON(&order)
		if err != nil {
			models.MConfig.MLogger.Error("ws read msg error: " + err.Error())
			gameClient.GameStatus = GStatusErrorLogout
			GLoginChan <- gameClient
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
		v := order.OrderType >> TYPE_TRANS << TYPE_TRANS
		switch v {
			case OT_Msg:
				dealOrderMsg(gameClient, &order)
			case OT_Skill:
				dealOrderSkill(gameClient, &order)
			case OT_Action:
				dealOrderAction(gameClient, &order)
			default:
				models.MConfig.MLogger.Error(string(gameClient.ID) + " order invalid: " + string(order.OrderType))
		}
	} 
}

/**
 * 处理消息指令：
 *		个体对个体的消息指令，则直接执行
 *		个体对个体自建群组的消息指令，则直接执行
 * 		个体对大众的消息指令，加入中心指令队列
 */
func dealOrderMsg(gameClient *GameClient, order *GameOrder) {
	switch order.OrderType {
		case OT_MsgPerson:
			var orderMsg GameOrderMsg
			err := mapstructure.Decode(order.Data, &orderMsg)
			if err != nil {
				models.MConfig.MLogger.Error("mapstructure.Decode error %s", err.Error())
				return
			}
			
			client_, ok := gameClient.GMap.gameClientMap.Load(orderMsg.ToID)
			if !ok {
				// fmt.Println(orderMsg)
				models.MConfig.MLogger.Error("gameClientMap.Load error")
				return
			}
			
			client, ok := (client_).(*GameClient)
			if !ok {
				models.MConfig.MLogger.Error("gameClientMap cast error")
				return
			}
			client.Conn.WriteJSON(order)
		case OT_MsgAll:
			var orderMsg GameOrderMsg
			err := mapstructure.Decode(order.Data, &orderMsg)
			if err != nil {
				models.MConfig.MLogger.Error("mapstructure.Decode error %s", err.Error())
				return
			}
			order.Data = orderMsg
			pushCenterOrder(order)
		default:
	}
}

func dealOrderSkill(gameClient *GameClient, order *GameOrder) {
	gameClient.GMap.dealOrderSkill(gameClient, order)
}

func dealOrderAction(gameClient *GameClient, order *GameOrder) {
	gameClient.GMap.dealOrderAction(gameClient, order)
}

func pushCenterOrder(order *GameOrder) {
	gMapMap.Range(func (key, gMap_ interface{}) bool {
		gMap, ok := gMap_.(*GMapService)
		if !ok {
			models.MConfig.MLogger.Error("dataCenter() gameClientMap cast error")
		} else {
			gMap.orderList.PushBack(&order)
		}
		return true
	})
}

func goDataCenter() {
	for {
		gMapMap.Range(func (key, gMap_ interface{}) bool {
			gMap, ok := gMap_.(*GMapService)
			if !ok {
				models.MConfig.MLogger.Error("dataCenter() gameClientMap cast error")
			} else {
				gMap.dataCenter()
			}
			return true
		})
		time.Sleep(time.Second * 1)
	}
}

func resetGameData(gameData *models.GameData) {
	gameData.GBloodAll = gameData.GBloodBase + 30000
}

func logoutAll() {
	gMapMap.Range(func (key , gMap_ interface{}) bool{
		gMap, ok := (gMap_).(*GMapService)
		if !ok {
			models.MConfig.MLogger.Error("addGameClient() gMap cast error")
		} else {
			gMap.logoutAll()
		}
		return true
	})
	gServerStatus = -1
}

func getUserData(id int64) *models.GameData {
	var gameData *models.GameData
	gMapMap.Range(func (key , gMap_ interface{}) bool{
		gMap, ok := (gMap_).(*GMapService)
		if !ok {
			models.MConfig.MLogger.Error("addGameClient() gMap cast error")
		} else {
			gameData = gMap.getUserData(id)
			return false
		}
		return true
	})
	return gameData
}