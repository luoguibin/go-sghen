package controllers

import (
	"SghenApi/models"
	"sync"
	"math/rand"
	"time"
	"fmt"
)

type GameManager struct {
	gameClientMap  	sync.Map

	loginChan		chan *models.GameClient
	gLoginService 	GLoginService
}

func (manager *GameManager) Init() {
	fmt.Println("GameManager::Init()")

	manager.loginChan = make(chan *models.GameClient)
	manager.gLoginService = GLoginService{}
	go manager.gLoginService.start()

	go manager.dataCenter()
}

func (manager *GameManager) gameClientHandle(gameClient *models.GameClient) {
	for { 
		var action models.GameAction
		err := gameClient.Conn.ReadJSON(&action)
		if err != nil {
			models.MConfig.MLogger.Error("ws read msg error: " + err.Error())
			break;
		}
		client_, ok := gameManager.gameClientMap.Load(action.Target)
		if !ok {
			models.MConfig.MLogger.Error("gameClientMap.Load error")
			break;
		}
		
		client, ok := (client_).(*models.GameClient)
		if !ok {
			models.MConfig.MLogger.Error("gameClientMap cast error")
			break;
		}

		switch action.Action {
			case "fist":
				ran := rand.Intn(200)
				if rand.Intn(10) < 5 {
					ran = gameClient.GameData.GPower + ran
				} else {
					ran = gameClient.GameData.GPower - ran
				}
				client.GameData.GBlood -= ran

				if client.GameData.GBlood < 0 {
					client.GameData.GBlood = 0
				}
			case "skill":
				ran := rand.Intn(200)
				ran = int(float32(gameClient.GameData.GPower) * 1.3) + ran
				client.GameData.GBlood -= ran

				if client.GameData.GBlood < 0 {
					client.GameData.GBlood = 0
				}
			case "skill_big":
				ran := rand.Intn(10000)
				ran = int(float32(gameClient.GameData.GPower) * 3.3) + ran
				client.GameData.GBlood -= ran

				if client.GameData.GBlood < 0 {
					client.GameData.GBlood = 0
				}
			case "msg":
				action.Target = gameClient.ID
				client.Conn.WriteJSON(action)
			default:
		}
	} 
}

func (manager *GameManager) dataCenter() {	
	for {
		// ①读取ws数据
		

		// ②计算
		wsDatas := make([]interface{}, 0)
		manager.gameClientMap.Range(func(key, client_ interface{}) bool {
			client, ok := (client_).(*models.GameClient)
			if !ok {
				models.MConfig.MLogger.Error("dataCenter() gameClientMap cast error")
				return true
			}
			wsDatas = append(wsDatas, client.GameData)
			return true
		})

		// ③发送ws数据
		count := 0
		manager.gameClientMap.Range(func(key, client_ interface{}) bool {
			client, ok := (client_).(*models.GameClient)
			count++
			if !ok {
				models.MConfig.MLogger.Error("dataCenter() gameClientMap cast error")
				return true
			}
			err := client.Conn.WriteJSON(wsDatas) 
			if err != nil { 
				models.MConfig.MLogger.Debug(err.Error())
				client.GameStatus = models.StatusLogout
				manager.loginChan <- client
			} 
			return true
		})

		time.Sleep(time.Second * 1)
		// fmt.Println("gameClientMap count=", count)
	}
}
