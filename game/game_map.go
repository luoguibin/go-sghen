package game

import (
	"go-sghen/helper"
	"go-sghen/models"
)

/*
 * game map size define, base on pixel
 */
const (
	// max map size
	GMapMaxWidth  = 12000
	GMapMaxHeight = 12000

	// csreen size unit, which is used to sort a map into many screens
	GMapSreenUnit = 1200

	// min map size
	GMapMinWidth  = GMapSreenUnit
	GMapMinHeight = GMapSreenUnit

	// a unit of the minimal object, which is used to calculate some collision detections
	// but only to check the collision between roles and buildings
	GMapUnit = 20
)

type GameMap struct {
	Name   string
	Width  int
	Height int

	ScreenXCount int
	ScreenYCount int
	ScreenCount  int

	Screens []map[int64]*GameClient

	// will use to store the obstruction and empty place
	// Data				[][]int
}

/*
 * init the other peoperties
 */
func (gameMap *GameMap) Init() {
	if gameMap.Width < GMapMinWidth {
		gameMap.Width = GMapMinWidth
	} else if gameMap.Width > GMapMaxWidth {
		gameMap.Width = GMapMaxWidth
	}

	if gameMap.Height < GMapMinHeight {
		gameMap.Height = GMapMinHeight
	} else if gameMap.Height > GMapMaxHeight {
		gameMap.Height = GMapMaxHeight
	}

	i, j := gameMap.GetScreenIndex2(gameMap.Width, gameMap.Height)
	gameMap.ScreenXCount = i + 1
	gameMap.ScreenYCount = j + 1
	gameMap.ScreenCount = gameMap.ScreenXCount * gameMap.ScreenYCount

	gameMap.Screens = make([]map[int64]*GameClient, gameMap.ScreenCount)
	for i := gameMap.ScreenCount - 1; i >= 0; i-- {
		gameMap.Screens[i] = make(map[int64]*GameClient, 0)
	}
}

/*
 * get the sreen index point(i, j) of the point(x, y)
 */
func (gameMap *GameMap) GetScreenIndex2(x, y int) (int, int) {
	if x < 0 {
		x = 0
	} else if x >= gameMap.Width {
		x = gameMap.Width - 1
	}

	if y < 0 {
		y = 0
	} else if y >= gameMap.Height {
		y = gameMap.Height - 1
	}

	i := x / GMapSreenUnit
	j := y / GMapSreenUnit

	return i, j
}

/*
 * get the sreen index of the point(x, y)
 */
func (gameMap *GameMap) GetScreenIndex(x, y int) int {
	i, j := gameMap.GetScreenIndex2(x, y)
	return gameMap.Index2ToIndex(i, j)
}

/*
 * change the (i, j) to screenId
 */
func (gameMap *GameMap) Index2ToIndex(i, j int) int {
	return i + j*gameMap.ScreenXCount
}

/*
 * change the `screenId` to (i, j)
 */
func (gameMap *GameMap) IndexToIndex2(screenId int) (int, int) {
	i := screenId % gameMap.ScreenXCount
	j := screenId / gameMap.ScreenXCount
	return i, j
}

/*
 * get the 9 screen increased ids
 */
func (gameMap *GameMap) GetScreenIndexs(screenId int) []int {
	i, _ := gameMap.IndexToIndex2(screenId)
	screenIds := make([]int, 0)

	topId := screenId - gameMap.ScreenXCount
	if topId >= 0 {

		if i > 0 {
			// left-top
			screenIds = append(screenIds, topId-1)
		}
		screenIds = append(screenIds, topId)
		if i < gameMap.ScreenXCount-1 {
			// right-top
			screenIds = append(screenIds, topId+1)
		}
	}
	// left
	if i > 0 {
		screenIds = append(screenIds, screenId-1)
	}
	// self
	screenIds = append(screenIds, screenId)
	// right
	if i < gameMap.ScreenXCount-1 {
		screenIds = append(screenIds, screenId+1)
	}

	bottomId := screenId + gameMap.ScreenXCount
	if bottomId < gameMap.ScreenCount {
		if i > 0 {
			// left-bottom
			screenIds = append(screenIds, bottomId-1)
		}
		// bottom
		screenIds = append(screenIds, bottomId)
		if i < gameMap.ScreenXCount-1 {
			// right-bottom
			screenIds = append(screenIds, bottomId+1)
		}
	}
	return screenIds
}

/*
 * init the screenId when first come into the map
 */
func (gameMap *GameMap) InitScreen(gameClient *GameClient) {
	index := gameMap.GetScreenIndex(gameClient.GameData.X, gameClient.GameData.Y)
	gameClient.GameData.ScreenId = index
	gameMap.Screens[index][gameClient.ID] = gameClient
}

/*
 * add a reflect between game client and the map
 */
func (gameMap *GameMap) ChangeScreen(gameClient *GameClient) {
	index := gameMap.GetScreenIndex(gameClient.GameData.X, gameClient.GameData.Y)
	preIndex := gameClient.GameData.ScreenId
	if index != preIndex {
		delete(gameMap.Screens[preIndex], gameClient.ID)
	}
	gameMap.Screens[index][gameClient.ID] = gameClient

	preIds := gameMap.GetScreenIndexs(preIndex)
	ids := gameMap.GetScreenIndexs(index)

	for i := len(ids) - 1; i >= 0; i-- {
		for j := len(preIds) - 1; j >= 0; j-- {
			if preIds[j] == ids[i] {
				ids = helper.GSliceRemove(ids, i)
				preIds = helper.GSliceRemove(preIds, j)
				break
			}
		}
	}

	// call the new screens clients that the client is new
	for _, screenId := range ids {
		gameMap.BroadCast1(screenId, GameOrder{
			OrderType: OT_ActionMoveAdd,
			FromID:    IDSYSTEM,
			FromType:  ITSystem,
			Data:      gameClient.ID,
		}, gameClient.ID)
	}
	// call the old screens clients that the client is old
	for _, screenId := range preIds {
		gameMap.BroadCast1(screenId, GameOrder{
			OrderType: OT_ActionMoveRemove,
			FromID:    IDSYSTEM,
			FromType:  ITSystem,
			Data:      gameClient.ID,
		}, gameClient.ID)
	}
}

/*
 * remove the screen reflect
 */
func (gameMap *GameMap) RemoveScreen(gameClient *GameClient) {
	delete(gameMap.Screens[gameClient.GameData.ScreenId], gameClient.ID)
	gameMap.BroadCast9(gameClient.GameData.ScreenId, GameOrder{
		OrderType: OT_ActionMoveRemove,
		FromID:    IDSYSTEM,
		FromType:  ITSystem,
		Data:      gameClient.ID,
	}, gameClient.ID)
}

/*
 * get the screen client
 */
func (gameMap *GameMap) GetGameClient(screenId int, id int64) *GameClient {
	client, ok := gameMap.Screens[screenId][id]
	if !ok {
		return nil
	} else {
		return client
	}
}

/*
 * send the 9 screen user datas to the `gameClient`
 */
func (gameMap *GameMap) SendGameDatas9(gameClient *GameClient) {
	screenIds := gameMap.GetScreenIndexs(gameClient.GameData.ScreenId)
	datas := make([]*models.GameData, 0)
	for _, screenId := range screenIds {
		for _, client := range gameMap.Screens[screenId] {
			if client.ID == gameClient.ID {
				continue
			}
			datas = append(datas, client.GameData)
		}
	}
	if len(datas) == 0 {
		return
	}
	gameClient.Conn.WriteJSON(GameOrder{
		OrderType: OT_DataAll,
		FromID:    IDSYSTEM,
		FromType:  ITSystem,
		Data:      datas,
	})
}

/*
 * broadcast the `order` to the clients of the map
 */
func (gameMap *GameMap) BroadCast(order interface{}) {
	for i := 0; i < gameMap.ScreenCount; i++ {
		gameMap.BroadCast1(i, order, 0)
	}
}

/*
 * broadcast the `order` to the clients of the nine screens
 */
func (gameMap *GameMap) BroadCast9(screenId int, order interface{}, exceptID int64) {
	screenIds := gameMap.GetScreenIndexs(screenId)
	for _, screenId := range screenIds {
		gameMap.BroadCast1(screenId, order, exceptID)
	}
}

/*
 * broadcast order in a screen
 */
func (gameMap *GameMap) BroadCast1(screenId int, order interface{}, exceptID int64) {
	for _, client := range gameMap.Screens[screenId] {
		if client.ID == exceptID {
			continue
		}
		client.Conn.WriteJSON(order)
	}
}
