package controllers

import (
	"SghenApi/models"
	"encoding/json"
	"strconv"
	"fmt"
)

type GLoginService struct {
	
}

func (service *GLoginService) start() {
	fmt.Println("GLoginService::Start()")
	for {
		v, ok := <- gameManager.loginChan
		if ok {
			switch v.GameStatus {
			case models.StatusLogin:
				checkLogin(v)
			case models.StatusLogout:
				checkLogout(v)
			case models.StatusLogoutAll:
				checkLogout(v)
			}
		} else {
			fmt.Println("loginChan get error")
		}
	}
}

func checkLogin(gameClient *models.GameClient) {
	if (GameStatus != 1) {
		gameClient.Conn.WriteJSON(models.GameOrder{
			OrderType: 	models.OrderSystemInfo,
			FromType:	models.FromSystem,
			FromID:		models.IDSystem,
			Data:		models.GameOrderMsg {
							ToType:		models.FromUser,
							ToID:		gameClient.ID,
							Msg: 		"系统维护中，代码：" + strconv.Itoa(GameStatus),
						},
		})
		gameClient.Conn.Close()
		return
	}
	_, ok := gameManager.gameClientMap.Load(gameClient.ID)
	if !ok {
		addGameUser(gameClient)
	} else {
		gameClient.Conn.WriteJSON(models.GameOrder{
			OrderType: 	models.OrderMsgSystem,
			FromType:	models.FromSystem,
			FromID:		models.IDSystem,
			Data:		models.GameOrderMsg {
							ToType:		models.FromUser,
							ToID:		gameClient.ID,
							Msg: 		"重复登录",
						},
		})
	}
}

func addGameUser(gameClient *models.GameClient) {
	gameData, err := models.QueryGameData(gameClient.ID)
	if err != nil {
		gameClient.Conn.WriteJSON(models.GameOrder{
			OrderType: 	models.OrderMsgSystem,
			FromType:	models.FromSystem,
			FromID:		models.IDSystem,
			Data:		models.GameOrderMsg {
							ToType:		models.FromUser,
							ToID:		gameClient.ID,
							Msg: 		"该账号下未查询到游戏数据",
						},
		})
		gameClient.Conn.Close()
		return
	}

	gameClient.GameData = gameData
	d, err := json.Marshal(gameClient.GameData)
	if err != nil {
		gameClient.Conn.WriteJSON(models.GameOrder{
			OrderType: 	models.OrderMsgSystem,
			FromType:	models.FromSystem,
			FromID:		models.IDSystem,
			Data:		models.GameOrderMsg {
							ToType:		models.FromUser,
							ToID:		gameClient.ID,
							Msg: 		"该账号下游戏数据解析出错",
						},
		})
		gameClient.Conn.Close()
		return
	}
	gameClient.Conn.WriteJSON(models.GameOrder{
		OrderType: 	models.OrderUserData,
		FromType:	models.FromSystem,
		FromID:		models.IDSystem,
		Data:		string(d),
	})
	
	gameManager.gameClientMap.Store(gameClient.ID, gameClient)
	go gameManager.gameClientHandle(gameClient)
}

func checkLogout(gameClient *models.GameClient) {
	gameClient.Conn.Close()
	gameManager.gameClientMap.Delete(gameClient.ID)

	err := models.UpdateGameData(gameClient.GameData)
	if err != nil {
		models.MConfig.MLogger.Error(err.Error())
	}
}