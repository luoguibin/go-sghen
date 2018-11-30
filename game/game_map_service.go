package game

import (
	"SghenApi/models"
	// "SghenApi/helper"
	"fmt"
	// "sort"
	"sync"
	"strconv"
	"github.com/goinggo/mapstructure"
)

type GameMapService struct {
	GameMap			sync.Map
	GameClientMap	*sync.Map
}

/*
 * map service start
 */
func (gameMapService *GameMapService) Start() {
	fmt.Println("GameMapService::Start()")

	gameMap	:= &GameMap{
		Name:		"HumanWorld",
		Width:		10000,
		Height:		10000,
	}
	gameMap.Init()
	gameMapService.GameMap.Store(0, gameMap)
	gameMapService.GameClientMap = &sync.Map{}
}

/*
 * add a gameClient into the gameMapService
 */
func (gameMapService *GameMapService) AddGameClient(gameClient *GameClient) {
	gameMapService.GameClientMap.Store(gameClient.ID, gameClient)
	gameMap_, ok := gameMapService.GameMap.Load(gameClient.GameData.MapId)
	if !ok {
		models.MConfig.MLogger.Error("AddGameClient() GameMap load error %v", gameClient.GameData.MapId)
		return
	}
	gameMap := gameMap_.(*GameMap)
	gameMap.ChangeScreen(gameClient)
	gameMap.BroadCast9(gameClient, GameOrder {
		OrderType:		OT_DataPersonLogin,
		FromID:			IDSYSTEM,
		FromType:		ITSystem,
		Data:			gameClient.GameData,
	})
}

/*
 * gameClient go into another map
 */
func (gameMapService *GameMapService) ChangeMapId(mapId int, gameClient *GameClient) {
	gameClient.GameData.MapId = mapId
	gameClient.GameData.X = 0
	gameClient.GameData.Y = 0
	gameMapService.AddGameClient(gameClient)
}

/*
 * remove a gameClient from the gameMapService
 */
 func (gameMapService *GameMapService) RemoveGameClient(gameClient *GameClient) {
	gameMapService.GameClientMap.Delete(gameClient.ID)
	ResetGameDataMove(gameClient.GameData, nil)

	gameMap_, ok := gameMapService.GameMap.Load(gameClient.GameData.MapId)
	if !ok {
		models.MConfig.MLogger.Error("AddGameClient() GameMap load error %v", gameClient.GameData.MapId)
		return
	}
	gameMap := gameMap_.(*GameMap)
	gameMap.RemoveScreen(gameClient)
}

/*
 * get the client by `id`
 */
func (gameMapService *GameMapService) GetUserData(id int64) *GameClient {
	v, ok := gameMapService.GameClientMap.Load(id)
	if !ok {
		return nil
	}
	client, ok := v.(*GameClient)
	if !ok { 
		return nil
	} else {
		return client
	}
}

/*
 * deal the msg order
 */
func (gameMapService *GameMapService) DealOrderMsg(gameClient *GameClient, order *GameOrder) {
	var orderMsg GameOrderMsg
	err := mapstructure.Decode(order.Data, &orderMsg)
	if err != nil {
		models.MConfig.MLogger.Error("DealOrderMsg() mapstructure.Decode error %s", err.Error())
		return
	}
	order.Data = orderMsg

	switch order.OrderType {
	case OT_MsgPerson:
		client := gameMapService.GetUserData(orderMsg.ToID)
		if client != nil {
			client.Conn.WriteJSON(order)
		}
	case OT_MsgNear:
		gameMap := gameMapService.GetGameMap(gameClient.GameData.MapId)
		if gameMap == nil {
			return
		}
		gameMap.BroadCast9(gameClient, order)
	case OT_MsgAll:
		gameMapService.GameClientMap.Range(func (key, v interface{}) bool {
			client, ok := v.(*GameClient)
			if !ok {
				models.MConfig.MLogger.Error("DealOrderMsg() GameClient cast error %s", key)
			} else {
				client.Conn.WriteJSON(order)
			}
			return true
		})
	default:
	}
}

/*
 * deal the skill order
 */
func (gameMapService *GameMapService) DealOrderSkill(gameClient *GameClient, order *GameOrder) {
	var orderSkill GameOrderSkill
	err := mapstructure.Decode(order.Data, &orderSkill)
	if err != nil {
		models.MConfig.MLogger.Error("mapstructure.Decode error %s", err.Error())
		return
	}

	skillID := order.OrderType
	switch skillID / 1000 * 1000 {
		case OT_SkillSingle:
			gameMap := gameMapService.GetGameMap(gameClient.GameData.MapId)
			if gameMap == nil {
				return
			}
			client := gameMap.GetGameClient(gameClient.GameData.ScreenId, orderSkill.ToID)
			if client == nil {
				return
			}
			
			data0 := gameClient.GameData
			data1 := client.GameData
			ResetGameDataMove(data0, nil)
			ResetGameDataMove(data1, nil)
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
				data1.Blood -= damage
	
				if data1.Blood < 0 {
					damage += data1.Blood
					data1.Blood = 0
				}
				orderSkill.Damage = damage
				orderSkill.DamageAll = damage
				orderSkill.DamageCount	= 1
				orderSkill.DamageCountAll = 1
				order.Data = orderSkill
			}
			gameMap.BroadCast9(gameClient, order)
		case OT_SkillSingleK:
		case OT_SkillNear:
			// s := make([]*GameSortItem, 0)
			// data0 := gameClient.GameData
			// ResetGameDataMove(data0, nil)
			// gameMapService.GameClientMap.Range(func(key, client_ interface{}) bool {
			// 	client, ok := (client_).(*GameClient)
			// 	if !ok {
			// 		models.MConfig.MLogger.Error("dataCenter() gameClientMap cast error")
			// 		return true
			// 	}
			// 	if client.ID == gameClient.ID {
			// 		 return true
			// 	}

			// 	data1 := client.GameData
			// 	ResetGameDataMove(data1, nil)
			// 	distance := helper.GClientDistance(data0.GX, data0.GY, data1.GX0, data1.GY)
			// 	if (distance < 180) {
			// 		s = append(s, &GameSortItem{
			// 			Value:			distance,
			// 			GameClient:		client,
			// 		})
			// 	}

			// 	return true
			// })

			// sort.Sort(GameSort(s)) 
			// count := 0
			// for _, v := range s {
			// 	data1 := v.GameClient.GameData
			// 	damage := getSkillSingleDamage(skillID, data0, data1)
			// 	if data1.GBlood <= 0 {
			// 		continue
			// 	}
			// 	if count > 6 {
			// 		break
			// 	}
			// 	count++
			// 	data1.GBlood -= damage
	
			// 	if data1.GBlood < 0 {
			// 		damage += data1.GBlood
			// 		data1.GBlood = 0
			// 	}
			// 	pushCenterOrder(&GameOrder {
			// 		OrderType:		order.OrderType,
			// 		FromID:			order.FromID,
			// 		FromType:		order.FromType,
			// 		Data:			GameOrderSkill{
			// 							ToID:			v.GameClient.ID,
			// 							Damage:			damage,
			// 							DamageAll:		damage,
			// 							DamageCount:	1,
			// 							DamageCountAll:	1,
			// 						},
			// 	})
			// }
		case OT_SkillNearK:
		default:
	}
}

/*
 * deal the action order
 */
func (gameMapService *GameMapService) DealOrderAction(gameClient *GameClient, order *GameOrder) {
	skillID := order.OrderType

	switch skillID / 1000 * 1000 {
		case OT_ActionDrug:
			data := gameClient.GameData
			data.Blood += data.BloodAll / 15
			if data.Blood > data.BloodAll {
				data.Blood = data.BloodAll
			}
		case OT_ActionMove:
			var orderAction GameOrderAction
			err := mapstructure.Decode(order.Data, &orderAction)
			if err != nil {
				models.MConfig.MLogger.Error("mapstructure.Decode error %s", err.Error())
				return
			}
			order.Data = orderAction
			ResetGameDataMove(gameClient.GameData, &orderAction)
			// fmt.Println(gameClient.GameData)
			// data.GX = orderAction.X
			// data.GY = orderAction.Y
		default:
	}
}

/*
 * get the game map by `id`
 */
func (gameMapService *GameMapService) GetGameMap(mapId int) *GameMap {
	v, ok := gameMapService.GameMap.Load(mapId)
	if !ok {
		models.MConfig.MLogger.Error("DealOrderMsg() load gameMap error %s", mapId)
		return nil
	} 
	gameMap, ok := v.(*GameMap)
	if !ok {
		models.MConfig.MLogger.Error("DealOrderMsg() gameMap cast error %s", mapId)
		return nil
	} 
	return gameMap
}

/*
 * logout all of the clients
 */
func (gameMapService *GameMapService) LogoutAll() {
	gameMapService.GameClientMap.Range(func(key, v interface{}) bool {
		client, ok := v.(*GameClient)
		if !ok {
			models.MConfig.MLogger.Error("LogoutAll() GameClient cast error: key=%v", key)
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
		MGameServer.GameAuthorityService.LoginChan <- client
		return true
	})
}