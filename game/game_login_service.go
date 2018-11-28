package game

import (
	"SghenApi/models"
	"strconv"
	"fmt"
)

var (
	GLoginChan	 = 	make(chan *GameClient)
)

type GLoginService struct {
	
}

func (service *GLoginService) start() {
	fmt.Println("GLoginService::Start()")
	for {
		v, ok := <- GLoginChan
		if ok {
			switch v.GameStatus {
				case GStatusLogin:
					checkLogin(v)
				case GStatusLogout:
					checkLogout(v)
				case GStatusErrorLogout:
					checkLogout(v)
				case GGStatusLogoutAll:
					checkLogout(v)
			}
		} else {
			fmt.Println("loginChan get error")
		}
	}
}

func checkLogin(gameClient *GameClient) {
	if (gServerStatus != 1) {
		gameClient.Conn.WriteJSON(GameOrder{
			OrderType: 	OT_MsgSystemInner,
			FromType:	ITSystem,
			FromID:		IDSYSTEM,
			Data:		GameOrderMsg {
							ToType:		ITPerson,
							ToID:		gameClient.ID,
							Msg: 		"系统维护中，代码：" + strconv.Itoa(gServerStatus),
						},
		})
		gameClient.Conn.Close()
		return
	}

	gameData := getUserData(gameClient.ID)
	
	if gameData != nil {
		gameClient.Conn.WriteJSON(GameOrder{
			OrderType: 	OT_MsgSystem,
			FromType:	ITSystem,
			FromID:		IDSYSTEM,
			Data:		GameOrderMsg {
							ToType:		ITPerson,
							ToID:		gameClient.ID,
							Msg: 		"重复登录",
						},
		})
	} else {
		addGameClient(gameClient)
	}
}

func addGameClient(gameClient *GameClient) {
	gameData, err := models.QueryGameData(gameClient.ID)
	if err != nil {
		gameClient.Conn.WriteJSON(GameOrder{
			OrderType: 	OT_MsgSystem,
			FromType:	ITSystem,
			FromID:		IDSYSTEM,
			Data:		GameOrderMsg {
							ToType:		ITPerson,
							ToID:		gameClient.ID,
							Msg: 		"该账号下未查询到游戏数据",
						},
		})
		gameClient.Conn.Close()
		return
	}
	resetGameData(gameData)

	gameClient.GameData = gameData
	gameClient.Conn.WriteJSON(GameOrder{
		OrderType: 	OT_DataPerson,
		FromType:	ITSystem,
		FromID:		IDSYSTEM,
		Data:		gameData,
	})

	gMap_, ok := gMapMap.Load(gameData.GMapId)
	if !ok {
		models.MConfig.MLogger.Error("addGameClient() gMap load error %s", gameData)
	} else {
		gMap, ok := (gMap_).(*GMapService)
		if !ok {
			models.MConfig.MLogger.Error("addGameClient() gMap cast error")
		} else {
			gameClient.GMap = gMap
			gMap.gameClientMap.Store(gameClient.ID, gameClient)
			go goGameClientHandle(gameClient)
		}
	}
}

func checkLogout(gameClient *GameClient) {
	gameClient.Conn.Close()
	gameClient.GMap.gameClientMap.Delete(gameClient.ID)
	
	err := models.UpdateGameData(gameClient.GameData)
	if err != nil {
		models.MConfig.MLogger.Error(err.Error())
	}
}