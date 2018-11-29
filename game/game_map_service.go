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

	switch skillID / 1000 * 1000 {
		case OT_SkillSingle:
			client_, ok := gMap.gameClientMap.Load(orderSkill.ToID)
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
			resetGameDataMove(data0, nil)
			resetGameDataMove(data1, nil)
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
					damage += data1.GBlood
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
}

func (gMap *GMapService) dealOrderAction(gameClient *GameClient, order *GameOrder) {
	skillID := order.OrderType

	switch skillID / 1000 * 1000 {
		case OT_ActionDrug:
			data := gameClient.GameData
			data.GBlood += data.GBloodAll / 15
			if data.GBlood > data.GBloodAll {
				data.GBlood = data.GBloodAll
			}
		case OT_ActionMove:
			var orderAction GameOrderAction
			err := mapstructure.Decode(order.Data, &orderAction)
			if err != nil {
				models.MConfig.MLogger.Error("mapstructure.Decode error %s", err.Error())
				return
			}
			order.Data = orderAction
			resetGameDataMove(gameClient.GameData, &orderAction)
			pushCenterOrder(order)
			// fmt.Println(gameClient.GameData)
			// data.GX = orderAction.X
			// data.GY = orderAction.Y
		default:
	}
}

func (gMap *GMapService) dataCenter() {	
	// ①读取ws数据
	orders := make([]interface{}, 0)
	for e := gMap.orderList.Front(); e != nil; e = e.Next() {
		order, ok := e.Value.(*GameOrder)
		if !ok {
			models.MConfig.MLogger.Error("dataCenter() GameOrder cast error")
			continue
		}
		orders = append(orders, order)
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
			models.MConfig.MLogger.Debug("dataCenter() " + err.Error())
			client.GameStatus = GStatusErrorLogout
			GLoginChan <- client
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