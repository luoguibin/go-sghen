package game

import (
	"SghenApi/models"
	// "encoding/json"
	"container/list"
	"sync"
	"strconv"
	"fmt"
	"github.com/goinggo/mapstructure"
)

type GMapService struct {
	name			string
	gameClientMap  	*sync.Map
	orderList		*list.List
	
	gameMapMap		sync.Map
}

func (gMap *GMapService) Init(name string) {
	fmt.Println("GMapService::Init() " + name)
	gMap.name = name
	gMap.orderList = list.New()
	gMap.gameClientMap = &sync.Map{}
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

func (gMap *GMapService) dealOrderMsg(gameClient *GameClient, order *GameOrder) {
}

func (gMap *GMapService) dealOrderSkill(gameClient *GameClient, order *GameOrder) {
	var orderSkill GameOrderSkill
	err := mapstructure.Decode(order.Data, &orderSkill)
	if err != nil {
		models.MConfig.MLogger.Error("mapstructure.Decode error %s", err.Error())
		return
	}

	skillID := order.OrderType

	switch skillID >> (TYPE_TRANS - 1) << (TYPE_TRANS - 1) {
		case OT_SkillSingle:
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
			damage := getSkillSingleDamage(skillID, data0, data1)
			
			if damage < 0 {
				gameClient.Conn.WriteJSON(GameOrder{
					OrderType: 	OT_MsgSystemPerson,
					FromType:	ITSystem,
					FromID:		IDSYSTEM,
					Data:		GameOrderMsg {
									ToType:		ITPerson,
									ToID:		gameClient.ID,
									Msg: 		"距离超过" + strconv.Itoa(-damage),
								},
				})
			} else {
				data1.GBlood -= damage
	
				if data1.GBlood < 0 {
					data1.GBlood = 0
				}
				orderSkill.Damage = damage
				orderSkill.DamageAll = damage
				orderSkill.DamageCount	= 1
				orderSkill.DamageCountAll = 1
				order.Data = orderSkill
				
				gMap.orderList.PushBack(order)
			}
		case OT_SkillSingleK:
		case OT_SkillNear:
		case OT_SkillNearK:
		default:
	}

	



	// switch client.GameData.GName {
	// 	case "fist":
	// 		data0 := gameClient.GameData
	// 		data1 := client.GameData
	// 		d := helper.GClientDistance(data0.GX, data0.GY, data1.GX, data1.GY)
	// 		if d > 50 {
	// 			gameClient.Conn.WriteJSON(GameOrder{
	// 				OrderType: 	models.OrderMsgFeedback,
	// 				FromType:	models.FromSystem,
	// 				FromID:		models.IDSystem,
	// 				Data:		GameOrderMsg {
	// 								ToType		models.FromUser,
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
	// 			gameClient.Conn.WriteJSON(GameOrder{
	// 				OrderType: 	models.OrderMsgFeedback,
	// 				FromType:	models.FromSystem,
	// 				FromID:		models.IDSystem,
	// 				Data:		GameOrderMsg {
	// 								ToType		models.FromUser,
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
	// 			gameClient.Conn.WriteJSON(GameOrder{
	// 				OrderType: 	models.OrderMsgFeedback,
	// 				FromType:	models.FromSystem,
	// 				FromID:		models.IDSystem,
	// 				Data:		GameOrderMsg {
	// 								ToType		models.FromUser,
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

func (gMap *GMapService) dealOrderAction(gameClient *GameClient, order *GameOrder) {
	skillID := order.OrderType

	switch skillID >> (TYPE_TRANS - 1) << (TYPE_TRANS - 1) {
		case OT_ActionDrug:
		case OT_ActionMove:
		default:
	}
}

func (gMap *GMapService) dataCenter() {	
	// ①读取ws数据
	orders := make([]interface{}, 0)
	for e := gMap.orderList.Front(); e != nil; e = e.Next() {
		orders = append(orders, e.Value.(*GameOrder))
	}

	// ②计算
	wsDatas := make([]interface{}, 0)
	gMap.gameClientMap.Range(func(key, client_ interface{}) bool {
		client, ok := (client_).(*GameClient)
		if !ok {
			models.MConfig.MLogger.Error("dataCenter() gameClientMap cast error")
			return true
		}

		client.GameData.GOrders = make([]*interface{}, 0)
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
		err := client.Conn.WriteJSON(GameOrder{
			OrderType:		OT_DataAll,
			FromType:		ITSystem,
			FromID:			IDSYSTEM,
			Data:			GameOrderData {
								Orders:		orders,
								Data:		wsDatas,
							},
		}) 
		if err != nil { 
			// models.MConfig.MLogger.Debug(err.Error())
			// client.GameStatus = GStatusErrorLogout
			// GLoginChan <- client
		} 
		return true
	})

	for e := gMap.orderList.Front(); e != nil;  {
		e_ := e.Next()
		gMap.orderList.Remove(e)
		e = e_
	}

	// fmt.Println("gameClientMap count=", count)
}

func (gMap *GMapService) logoutAll() {
	gMap.gameClientMap.Range(func(key, client_ interface{}) bool {
		client, ok := (client_).(*GameClient)
		if !ok {
			models.MConfig.MLogger.Error("dataCenter() gameClientMap cast error")
			return true
		}
		client.Conn.WriteJSON(GameOrder{
			OrderType:		OT_MsgSystem,
			FromType:		ITSystem,
			FromID:			IDSYSTEM,
			Data:			GameOrderMsg {
								ToType:		ITPerson,
								ToID:		client.ID,
								Msg: 		"系统强制离线",
							},
		})

		client.GameStatus = GGStatusLogoutAll
		GLoginChan <- client
		return true
	})
}