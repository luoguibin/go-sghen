package game

import (
	"fmt"
	"go-sghen/helper"
	"go-sghen/models"
	"sort"
	"strconv"
	"sync"

	"github.com/goinggo/mapstructure"
)

type GameMapService struct {
	GameMap       sync.Map
	GameClientMap *sync.Map
}

/*
 * map service start
 */
func (gameMapService *GameMapService) Start() {
	fmt.Println("GameMapService::Start()")

	gameMap := &GameMap{
		Name:   "HumanWorld",
		Width:  10000,
		Height: 10000,
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
	gameMap.InitScreen(gameClient)

	// broadcast the gameClient's data to the clients of the 9 screens
	gameMap.BroadCast9(gameClient.GameData.ScreenId, models.GameOrder{
		OrderType: OT_DataPersonLogin,
		FromID:    IDSYSTEM,
		FromType:  ITSystem,
		Data:      gameClient.GameData,
	}, gameClient.ID)
	// send the client datas of the 9 screens to the gameClient
	gameMap.SendGameDatas9(gameClient)
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
func (gameMapService *GameMapService) DealOrderMsg(gameClient *GameClient, order *models.GameOrder) {
	var orderMsg models.GameOrderMsg
	err := mapstructure.Decode(order.Data, &orderMsg)
	if err != nil {
		models.MConfig.MLogger.Error("DealOrderMsg() mapstructure.Decode error %s", err.Error())
		return
	}
	orderMsg.FromName = gameClient.GameData.Name
	order.Data = orderMsg

	switch order.OrderType {
	case OT_MsgPerson:
		client := gameMapService.GetUserData(orderMsg.ToID)
		if client != nil {
			client.Conn.WriteJSON(order)
		}
	case OT_MsgPersonLogout:
		gameClient.GameStatus = GStatusLogout
		MGameServer.GameAuthorityService.LoginChan <- gameClient
	case OT_MsgNear:
		gameMap := gameMapService.GetGameMap(gameClient.GameData.MapId)
		if gameMap == nil {
			return
		}
		gameMap.BroadCast9(gameClient.GameData.ScreenId, order, 0)
	case OT_MsgAll:
		gameMapService.GameClientMap.Range(func(key, v interface{}) bool {
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
func (gameMapService *GameMapService) DealOrderSkill(gameClient *GameClient, order *models.GameOrder) {
	var orderSkill models.GameOrderSkill
	err := mapstructure.Decode(order.Data, &orderSkill)
	if err != nil {
		models.MConfig.MLogger.Error("mapstructure.Decode error %s", err.Error())
		return
	}
	gameMap := gameMapService.GetGameMap(gameClient.GameData.MapId)
	if gameMap == nil {
		return
	}
	skillID := order.OrderType
	switch skillID / 1000 * 1000 {
	case OT_SkillSingle:
		client := gameMap.GetGameClient(gameClient.GameData.ScreenId, orderSkill.ToID)
		if client == nil {
			return
		}

		data0 := gameClient.GameData
		data1 := client.GameData
		ResetGameDataMove(data0, nil)
		ResetGameDataMove(data1, nil)
		damage := getSkillSingleDamage(skillID, data0, data1, 5)

		if damage < 0 {
			gameClient.Conn.WriteJSON(models.GameOrder{
				OrderType: OT_MsgSystemPerson,
				FromType:  ITSystem,
				FromID:    IDSYSTEM,
				Data: models.GameOrderMsg{
					ToType: ITPerson,
					ToID:   gameClient.ID,
					Msg:    "距离超过" + strconv.Itoa(-damage),
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
			orderSkill.DamageCount = 1
			orderSkill.DamageCountAll = 1
			order.Data = orderSkill
			gameMap.BroadCast9(gameClient.GameData.ScreenId, order, 0)
		}
	case OT_SkillSingleK:
	case OT_SkillNear:
		s := make([]*GameSortItem, 0)
		data0 := gameClient.GameData
		ResetGameDataMove(data0, nil)
		gameMapService.GameClientMap.Range(func(key, client_ interface{}) bool {
			client, ok := (client_).(*GameClient)
			if !ok {
				models.MConfig.MLogger.Error("dataCenter() gameClientMap cast error")
				return true
			}
			if client.ID == gameClient.ID {
				return true
			}

			data1 := client.GameData
			ResetGameDataMove(data1, nil)
			distance := helper.GClientDistance(data0.X, data0.Y, data1.X, data1.Y)
			if distance < 10 {
				s = append(s, &GameSortItem{
					Value:      distance,
					GameClient: client,
				})
			}

			return true
		})

		sort.Sort(GameSort(s))
		count := 0
		orderSkills := make([]models.GameOrderSkill, 0)
		for _, v := range s {
			data1 := v.GameClient.GameData
			damage := getSkillSingleDamage(skillID, data0, data1, 10)
			if data1.Blood <= 0 {
				continue
			}
			if count > 6 {
				break
			}
			count++
			data1.Blood -= damage

			if data1.Blood < 0 {
				damage += data1.Blood
				data1.Blood = 0
			}
			orderSkills = append(orderSkills, models.GameOrderSkill{
				ToID:           data1.ID,
				Damage:         damage,
				DamageAll:      damage,
				DamageCount:    0,
				DamageCountAll: 0,
			})
		}
		order.Data = orderSkills
		gameMap.BroadCast9(gameClient.GameData.ScreenId, order, 0)
	case OT_SkillNearK:
	default:
	}
}

/*
 * deal the action order
 */
func (gameMapService *GameMapService) DealOrderAction(gameClient *GameClient, order *models.GameOrder) {
	skillID := order.OrderType
	gameMap := gameMapService.GetGameMap(gameClient.GameData.MapId)
	if gameMap == nil {
		return
	}

	switch skillID / 1000 * 1000 {
	case OT_ActionDrug:
		data := gameClient.GameData
		addBlood := data.BloodAll / 15
		data.Blood += addBlood
		if data.Blood > data.BloodAll {
			addBlood -= data.Blood - data.BloodAll
			data.Blood = data.BloodAll
		}
		order.Data = addBlood
		gameMap.BroadCast9(gameClient.GameData.ScreenId, order, 0)
	case OT_ActionMove:
		var orderAction []models.GameOrderAction
		err := mapstructure.Decode(order.Data, &orderAction)
		if err != nil {
			models.MConfig.MLogger.Error("mapstructure.Decode error %s", err.Error())
			return
		}
		ResetGameDataMove(gameClient.GameData, orderAction)
		gameMap.BroadCast9(gameClient.GameData.ScreenId, order, 0)
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
		client.Conn.WriteJSON(models.GameOrder{
			OrderType: OT_MsgSystem,
			FromType:  ITSystem,
			FromID:    IDSYSTEM,
			Data: models.GameOrderMsg{
				ToType: ITPerson,
				ToID:   client.ID,
				Msg:    "系统强制离线",
			},
		})

		client.GameStatus = GGStatusLogoutAll
		MGameServer.GameAuthorityService.LoginChan <- client
		return true
	})
}
