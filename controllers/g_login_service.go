package controllers

import (
	"SghenApi/models"
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
			}
		} else {
			fmt.Println("loginChan get error")
		}
	}
}

func checkLogin(gameClient *models.GameClient) {
	_, ok := gameManager.gameClientMap.Load(gameClient.ID)
	if !ok {
		addGameUser(gameClient)
	} else {
		gameClient.Conn.WriteJSON(models.GameAction{Msg: "重复登录"})
	}
}

func addGameUser(gameClient *models.GameClient) {
	gameData, err := models.QueryGameData(gameClient.ID)
	if err != nil {
		gameClient.Conn.WriteJSON(models.GameAction{Msg: "该账号下未查询到游戏数据"})
		gameClient.Conn.Close()
		return
	}

	gameClient.GameData = gameData
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