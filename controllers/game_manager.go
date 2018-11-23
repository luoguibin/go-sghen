package controllers

import (
	"SghenApi/models"
	// "SghenApi/helper"
	// "encoding/json"
	"sync"
	// "math/rand"
	"time"
	"fmt"
	"github.com/goinggo/mapstructure"
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
	preTime := time.Now().UnixNano() / 1e6
	for {
		// 获取指令 
		var order models.GameOrder
		err := gameClient.Conn.ReadJSON(&order)
		if err != nil {
			models.MConfig.MLogger.Error("ws read msg error: " + err.Error())
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
		v := order.OrderType >> 3 << 3
		switch v {
			case models.OrderMsg:
				manager.dealOrderMsg(gameClient, order)
			// case models.OrderSkill:
			// 	manager.dealOrderSkill(gameClient, order)
			// case models.OrderNormal:
			// 	manager.dealOrderNormal(gameClient, order)
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
func (manager *GameManager) dealOrderMsg(gameClient *models.GameClient, order models.GameOrder) bool{
	switch order.OrderType {
		case models.OrderMsgPerson:
			var orderMsg models.GameOrderMsg
			err := mapstructure.Decode(order.Data, &orderMsg)
			if err != nil {
				models.MConfig.MLogger.Error("mapstructure.Decode error %s", err.Error())
				return true
			}
			client_, ok := manager.gameClientMap.Load(orderMsg.ToID)
			if !ok {
				fmt.Println(orderMsg)
				models.MConfig.MLogger.Error("gameClientMap.Load error")
				return true;
			}
			
			client, ok := (client_).(*models.GameClient)
			if !ok {
				models.MConfig.MLogger.Error("gameClientMap cast error")
				return true;
			}
			client.Conn.WriteJSON(order)
			return true;
		case models.OrderMsgGroup:
			return true;
		case models.OrderMsgAll:
			return false;
		case models.OrderMsgSystem:
			return false;
		default:
			return false;
	}
}

// func (manager *GameManager) dealOrderSkill(gameClient *models.GameClient, order models.GameOrder) {
// 	client_, ok := manager.gameClientMap.Load(order.Target)
// 	if !ok {
// 		models.MConfig.MLogger.Error("gameClientMap.Load error")
// 		return;
// 	}
	
// 	client, ok := (client_).(*models.GameClient)
// 	if !ok {
// 		models.MConfig.MLogger.Error("gameClientMap cast error")
// 		return;
// 	}

// 	switch order.Msg {
// 		case "fist":
// 			data0 := gameClient.GameData
// 			data1 := client.GameData
// 			d := helper.GClientDistance(data0.GX, data0.GY, data1.GX, data1.GY)
// 			if d > 50 {
// 				gameClient.Conn.WriteJSON(models.GameOrder{
// 					OrderType: 	models.OrderMsg,
// 					Target:		-1,
// 					Msg: 		"距离超过50",
// 				})
// 				break
// 			}
// 			ran := rand.Intn(200)
// 			if rand.Intn(10) < 5 {
// 				ran = data0.GPower + ran
// 			} else {
// 				ran = data0.GPower - ran
// 			}
// 			data1.GBlood -= ran

// 			if data1.GBlood < 0 {
// 				data1.GBlood = 0
// 			}
// 		case "skill":
// 			data0 := gameClient.GameData
// 			data1 := client.GameData
// 			d := helper.GClientDistance(data0.GX, data0.GY, data1.GX, data1.GY)
// 			if d > 50 {
// 				gameClient.Conn.WriteJSON(models.GameOrder{
// 					OrderType: 	models.OrderMsg,
// 					Target:		-1,
// 					Msg: 		"距离超过50",
// 				})
// 				break
// 			}
// 			ran := rand.Intn(200)
// 			ran = int(float32(data0.GPower) * 1.3) + ran
// 			data1.GBlood -= ran

// 			if data1.GBlood < 0 {
// 				data1.GBlood = 0
// 			}
// 		case "skill_big":
// 			data0 := gameClient.GameData
// 			data1 := client.GameData
// 			d := helper.GClientDistance(data0.GX, data0.GY, data1.GX, data1.GY)
// 			if d > 80 {
// 				gameClient.Conn.WriteJSON(models.GameOrder{
// 					OrderType: 	models.OrderMsg,
// 					Target:		-1,
// 					Msg: 		"距离超过80",
// 				})
// 				break
// 			}
// 			ran := rand.Intn(10000)
// 			ran = int(float32(data0.GPower) * 3.3) + ran
// 			data1.GBlood -= ran

// 			if data1.GBlood < 0 {
// 				data1.GBlood = 0
// 			}
// 		default:
// 	}
// }

// func (manager *GameManager) dealOrderNormal(gameClient *models.GameClient, order models.GameOrder) {
// 	// client_, ok := manager.gameClientMap.Load(order.Target)
// 	// if !ok {
// 	// 	models.MConfig.MLogger.Error("gameClientMap.Load error")
// 	// 	return;
// 	// }
	
// 	// client, ok := (client_).(*models.GameClient)
// 	// if !ok {
// 	// 	models.MConfig.MLogger.Error("gameClientMap cast error")
// 	// 	return;
// 	// }

// 	switch order.Msg {
// 		case "drug":
// 			data0 := gameClient.GameData
// 			data0.GBlood += data0.GBloodAll / 10
// 			if data0.GBlood > data0.GBloodAll {
// 				data0.GBlood = data0.GBloodAll
// 			}
// 		case "action":
// 			b := []byte(order.Data)
// 			action := models.GameAction{}
// 			err := json.Unmarshal(b, &action)
// 			if err != nil {
// 				models.MConfig.MLogger.Error("dealOrderNormal() action parse err " + err.Error())
// 				break
// 			}

// 			data0 := gameClient.GameData
// 			data0.GX = action.GX
// 			data0.GY = action.GY
// 		default:
// 	}
// }

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
			err := client.Conn.WriteJSON(models.GameOrder{
				OrderType:		models.OrderGameData,
				FromType:		models.FromSystem,
				FromID:			models.IDSystem,
				Data:			wsDatas,
			}) 
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
