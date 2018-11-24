package controllers

import (
	"SghenApi/models"
	"SghenApi/helper"
	// "encoding/json"
	"container/list"
	"sync"
	"math/rand"
	"time"
	"fmt"
	"github.com/goinggo/mapstructure"
)

type GMapService struct {
	name			string
	gameClientMap  	sync.Map
	orderList		*list.List
	
	gameMapMap		sync.Map
}

func (gMap *GMapService) Init(name string) {
	fmt.Println("GMapService::Init() " + name)
	gMap.name = name
	gMap.orderList = list.New()

	go gMap.dataCenter()
}

func (gMap *GMapService) getUserData(id int64) *models.GameData {
	client_, ok := gMap.gameClientMap.Load(id)
	if !ok {
		return nil
	}
	client, ok := (client_).(*GameClient)
	if !ok { 
		return nil
	}
	return client.GameData
}

func (gMap *GMapService) gameClientHandle(gameClient *GameClient) {
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
			case OTypeMsg:
				ok := gMap.dealOrderMsg(gameClient, &order)
				if !ok {
					gMap.orderList.PushBack(&order)
				}
			case OTypeSkill:
				gMap.dealOrderSkill(gameClient, &order)
			// case models.OrderNormal:
			// 	gMap.dealOrderNormal(gameClient, order)
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
func (gMap *GMapService) dealOrderMsg(gameClient *GameClient, order *models.GameOrder) bool{
	switch order.OrderType {
		case OTypeMsgPerson:
			var orderMsg models.GameOrderMsg
			err := mapstructure.Decode(order.Data, &orderMsg)
			if err != nil {
				models.MConfig.MLogger.Error("mapstructure.Decode error %s", err.Error())
				return true
			}
			client_, ok := gMap.gameClientMap.Load(orderMsg.ToID)
			if !ok {
				// fmt.Println(orderMsg)
				models.MConfig.MLogger.Error("gameClientMap.Load error")
				return true;
			}
			
			client, ok := (client_).(*GameClient)
			if !ok {
				models.MConfig.MLogger.Error("gameClientMap cast error")
				return true;
			}
			client.Conn.WriteJSON(order)
			return true;
		case OTypeMsgGroup:
			return true;
		case OTypeMsgAll:
			var orderMsg models.GameOrderMsg
			err := mapstructure.Decode(order.Data, &orderMsg)
			if err != nil {
				models.MConfig.MLogger.Error("mapstructure.Decode error %s", err.Error())
				return true
			}
			order.Data = orderMsg
			return false;
		case OTypeMsgSystem:
			return false;
		default:
			return false;
	}
}

func (gMap *GMapService) dealOrderSkill(gameClient *GameClient, order *models.GameOrder) {
	var orderSkill models.GameOrderSkill
	err := mapstructure.Decode(order.Data, &orderSkill)
	if err != nil {
		models.MConfig.MLogger.Error("mapstructure.Decode error %s", err.Error())
		return
	}

	skillID := orderSkill.SkillID
	if skillID < 1000 {
		models.MConfig.MLogger.Error("skillID < 1000")
		return
	}

	switch skillID >> 3 << 3 {
		case models.SkillS:
			// single oject skill
			client_, ok := gMap.gameClientMap.Load(order.FromID)
			if !ok {
				models.MConfig.MLogger.Error("gameClientMap.Load error")
				return;
			}
			
			client, ok := (client_).(*GameClient)
			if !ok {
				models.MConfig.MLogger.Error("gameClientMap cast error", client.ID)
				return;
			}
			data0 := gameClient.GameData
			data1 := client.GameData
			switch skillID {
				case models.SkillS1:
					d := helper.GClientDistance(data0.GX, data0.GY, data1.GX, data1.GY)
					if d > 50 {
						gameClient.Conn.WriteJSON(models.GameOrder{
							OrderType: 	OTypeMsgSystem,
							FromType:	ITypeSystem,
							FromID:		IDSYSTEM,
							Data:		models.GameOrderMsg {
											ToType:		ITypePerson,
											ToID:		gameClient.ID,
											Msg: 		"距离超过50",
										},
						})
						break
					}
					ran := rand.Intn(100)
					if rand.Intn(10) < 5 {
						ran = data0.GSpear.SStrength + ran
					} else {
						ran = data0.GSpear.SStrength - ran
					}
					data1.GBlood -= ran
		
					if data1.GBlood < 0 {
						data1.GBlood = 0
					}
					orderSkill.Damage = ran
					orderSkill.DamageAll = ran
					orderSkill.DamageCount	= 1
					orderSkill.DamageCountAll = 1
					order.Data = orderSkill
					
					gMap.orderList.PushBack(order)
				default:	
			}
		case models.SkillG:
		default:
	}

	



	// switch client.GameData.GName {
	// 	case "fist":
	// 		data0 := gameClient.GameData
	// 		data1 := client.GameData
	// 		d := helper.GClientDistance(data0.GX, data0.GY, data1.GX, data1.GY)
	// 		if d > 50 {
	// 			gameClient.Conn.WriteJSON(models.GameOrder{
	// 				OrderType: 	models.OrderMsgFeedback,
	// 				FromType:	models.FromSystem,
	// 				FromID:		models.IDSystem,
	// 				Data:		models.GameOrderMsg {
	// 								ToType:		models.FromUser,
	// 								ToID:		gameClient.ID,
	// 								Msg: 		"距离超过50",
	// 							},
	// 			})
	// 			break
	// 		}
	// 		ran := rand.Intn(200)
	// 		if rand.Intn(10) < 5 {
	// 			ran = data0.GSpear.SStrength + ran
	// 		} else {
	// 			ran = data0.GSpear.SStrength - ran
	// 		}
	// 		data1.GBlood -= ran

	// 		if data1.GBlood < 0 {
	// 			data1.GBlood = 0
	// 		}
	// 	case "skill":
	// 		data0 := gameClient.GameData
	// 		data1 := client.GameData
	// 		d := helper.GClientDistance(data0.GX, data0.GY, data1.GX, data1.GY)
	// 		if d > 50 {
	// 			gameClient.Conn.WriteJSON(models.GameOrder{
	// 				OrderType: 	models.OrderMsgFeedback,
	// 				FromType:	models.FromSystem,
	// 				FromID:		models.IDSystem,
	// 				Data:		models.GameOrderMsg {
	// 								ToType:		models.FromUser,
	// 								ToID:		gameClient.ID,
	// 								Msg: 		"距离超过50",
	// 							},
	// 			})
	// 			break
	// 		}
	// 		ran := rand.Intn(200)
	// 		ran = int(float32(data0.GSpear.SStrength) * 1.3) + ran
	// 		data1.GBlood -= ran

	// 		if data1.GBlood < 0 {
	// 			data1.GBlood = 0
	// 		}
	// 	case "skill_big":
	// 		data0 := gameClient.GameData
	// 		data1 := client.GameData
	// 		d := helper.GClientDistance(data0.GX, data0.GY, data1.GX, data1.GY)
	// 		if d > 80 {
	// 			gameClient.Conn.WriteJSON(models.GameOrder{
	// 				OrderType: 	models.OrderMsgFeedback,
	// 				FromType:	models.FromSystem,
	// 				FromID:		models.IDSystem,
	// 				Data:		models.GameOrderMsg {
	// 								ToType:		models.FromUser,
	// 								ToID:		gameClient.ID,
	// 								Msg: 		"距离超过80",
	// 							},
	// 			})
	// 			break
	// 		}
	// 		ran := rand.Intn(10000)
	// 		ran = int(float32(data0.GSpear.SStrength) * 3.3) + ran
	// 		data1.GBlood -= ran

	// 		if data1.GBlood < 0 {
	// 			data1.GBlood = 0
	// 		}
	// 	default:
	// }
}

// func (gMap *GMapService) dealOrderNormal(gameClient *GameClient, order models.GameOrder) {
// 	// client_, ok := gMap.gameClientMap.Load(order.Target)
// 	// if !ok {
// 	// 	models.MConfig.MLogger.Error("gameClientMap.Load error")
// 	// 	return;
// 	// }
	
// 	// client, ok := (client_).(*GameClient)
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

func (gMap *GMapService) dataCenter() {	
	for {
		// ①读取ws数据
		orders := make([]interface{}, 0)
		for e := gMap.orderList.Front(); e != nil; e = e.Next() {
			orders = append(orders, e.Value.(*models.GameOrder))
		}

		// ②计算
		wsDatas := make([]interface{}, 0)
		gMap.gameClientMap.Range(func(key, client_ interface{}) bool {
			client, ok := (client_).(*GameClient)
			if !ok {
				models.MConfig.MLogger.Error("dataCenter() gameClientMap cast error")
				return true
			}

			client.GameData.GOrders = make([]*models.GameOrder, 0)

			wsDatas = append(wsDatas, client.GameData)
			return true
		})

		// ③发送ws数据
		count := 0
		gMap.gameClientMap.Range(func(key, client_ interface{}) bool {
			client, ok := (client_).(*GameClient)
			count++
			if !ok {
				models.MConfig.MLogger.Error("dataCenter() gameClientMap cast error")
				return true
			}
			err := client.Conn.WriteJSON(models.GameOrder{
				OrderType:		OTypeDataAll,
				FromType:		ITypeSystem,
				FromID:			IDSYSTEM,
				Data:			models.GameOrderData {
									Orders:		orders,
									Data:		wsDatas,
								},
			}) 
			if err != nil { 
				models.MConfig.MLogger.Debug(err.Error())
				client.GameStatus = GStatusLogout
				GLoginChan <- client
			} 
			return true
		})

		for e := gMap.orderList.Front(); e != nil;  {
			e_ := e.Next()
			gMap.orderList.Remove(e)
			e = e_
		}

		time.Sleep(time.Second * 1)
		// fmt.Println("gameClientMap count=", count)
	}
}

func (gMap *GMapService) logoutAll() {
	gMap.gameClientMap.Range(func(key, client_ interface{}) bool {
		client, ok := (client_).(*GameClient)
		if !ok {
			models.MConfig.MLogger.Error("dataCenter() gameClientMap cast error")
			return true
		}
		client.Conn.WriteJSON(models.GameOrder{
			OrderType:		OTypeMsgSystem,
			FromType:		ITypeSystem,
			FromID:			IDSYSTEM,
			Data:			models.GameOrderMsg {
								ToType:		ITypePerson,
								ToID:		client.ID,
								Msg: 		"系统强制离线",
							},
		})

		client.GameStatus = GGStatusLogoutAll
		GLoginChan <- client
		return true
	})
}