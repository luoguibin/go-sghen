package game

import (
	"fmt"
	"go-sghen/models"
	"strconv"
)

type GameAuthorityService struct {
	LoginChan chan *GameClient
}

func (service *GameAuthorityService) Start() {
	fmt.Println("GameAuthorityService::Start()")

	service.LoginChan = make(chan *GameClient)
	for {
		v, ok := <-service.LoginChan
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
	if GameServerStatus != 1 {
		gameClient.Conn.WriteJSON(GameOrder{
			OrderType: OT_MsgSystemInner,
			FromType:  ITSystem,
			FromID:    IDSYSTEM,
			Data: GameOrderMsg{
				ToType: ITPerson,
				ToID:   gameClient.ID,
				Msg:    "系统维护中，代码：" + strconv.Itoa(GameServerStatus),
			},
		})
		gameClient.Conn.Close()
		return
	}

	client := MGameServer.GameMapService.GetUserData(gameClient.ID)

	if client != nil && client.GameData != nil {
		gameClient.Conn.WriteJSON(GameOrder{
			OrderType: OT_MsgSystemInner,
			FromType:  ITSystem,
			FromID:    IDSYSTEM,
			Data: GameOrderMsg{
				ToType: ITPerson,
				ToID:   gameClient.ID,
				Msg:    "重复登录",
			},
		})
		gameClient.Conn.Close()
	} else {
		addGameClient(gameClient)
	}
}

func addGameClient(gameClient *GameClient) {
	gameData, err := models.QueryGameData(gameClient.ID)
	if err != nil {
		gameClient.Conn.WriteJSON(GameOrder{
			OrderType: OT_MsgSystem,
			FromType:  ITSystem,
			FromID:    IDSYSTEM,
			Data: GameOrderMsg{
				ToType: ITPerson,
				ToID:   gameClient.ID,
				Msg:    "该账号下未查询到游戏数据",
			},
		})
		gameClient.Conn.Close()
		return
	}
	ResetGameData(gameData)

	gameClient.GameData = gameData
	gameClient.Conn.WriteJSON(GameOrder{
		OrderType: OT_DataPerson,
		FromType:  ITSystem,
		FromID:    IDSYSTEM,
		Data:      gameData,
	})

	MGameServer.GameMapService.AddGameClient(gameClient)
	go GoGameClientHandle(gameClient)
}

func checkLogout(gameClient *GameClient) {
	gameClient.Conn.Close()
	MGameServer.GameMapService.RemoveGameClient(gameClient)
	err := models.UpdateGameData(gameClient.GameData)
	if err != nil {
		models.MConfig.MLogger.Error(err.Error())
	}
}
