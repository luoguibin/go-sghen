package controllers

import(
	"SghenApi/models"
	"strconv"
	"fmt"
	"net/http"
	"sync"
	"github.com/gorilla/websocket"
)

var (
	upgrader 		= websocket.Upgrader{
						ReadBufferSize: 	2048,
						WriteBufferSize: 	2048,
						CheckOrigin: func(r *http.Request) bool {
							return true
						},
					}

	GameStatus		=	1

	gLoginService 	GLoginService
	GMapMap 		sync.Map	
)


type WSServerController struct {
	BaseController
}


func init() {
	fmt.Println("WSServerController::init()");

	gMapService0 := &GMapService{}
	gMapService0.Init("HumanWorld ")
	GMapMap.Store(0, gMapService0)

	// gMapService1 := &GMapService{}
	// gMapService1.Init("Hell")
	// GMapMap.Store(1, gMapService1)

	// gMapService2 := &GMapService{}
	// gMapService2.Init("Heaven")
	// GMapMap.Store(2, gMapService2)

	go gLoginService.start()
	go gameCommond()
}

func gameCommond() {
	for {
		var str string
		fmt.Scan(&str)

		if (str == "game:finish") {
			logoutAll()
		} else if (str == "game:start") {
			GameStatus = 1
		} else if (str == "game:getUser") {
			fmt.Print("input id:")
			fmt.Scan(&str)
			id, err := strconv.ParseInt(str, 10, 64)
			if err == nil {
				getUserData(id)
			} else {
				fmt.Println(err)
			}
		}
	}
}


func ResetGameData(gameData *models.GameData) {
	gameData.GBloodAll = gameData.GBloodBase + 30000
}

/**
 * WebSocket连接入口
 * 在BeforeRouter检测jwt中的合法后才给予长连接
 */
func (c *WSServerController) Get() {
	ws, err := upgrader.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil) 
	if err != nil { 
		models.MConfig.MLogger.Error("get ws error:\n%s", err)
	} 

	models.MConfig.MLogger.Debug("get ws: " + ws.RemoteAddr().String())

	uId,_ := strconv.ParseInt(c.Ctx.Input.Query("uId"), 10, 64)
	gameClient := &GameClient {
		ID:				uId,
		Conn:			ws,
		GameStatus:		GStatusLogin,
	}
	GLoginChan <- gameClient
}

func logoutAll() {
	GMapMap.Range(func (key , gMap_ interface{}) bool{
		gMap, ok := (gMap_).(*GMapService)
		if !ok {
			models.MConfig.MLogger.Error("addGameClient() gMap cast error")
		} else {
			gMap.logoutAll()
		}
		return true
	})
	GameStatus = -1
}

func getUserData(id int64) {
	GMapMap.Range(func (key , gMap_ interface{}) bool{
		gMap, ok := (gMap_).(*GMapService)
		if !ok {
			models.MConfig.MLogger.Error("addGameClient() gMap cast error")
		} else {
			data := gMap.getUserData(id)
			fmt.Println(data)
		}
		return true
	})
}