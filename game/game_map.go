package game

/* 
 * game map size define, base on pixel 
 */
const (
	// max map size
    GMapMaxWidth        =   12000
	GMapMaxHeight       =   12000

	// csreen size unit, which is used to sort a map into many screens
	GMapSreenUnit       =   1200

	// min map size
    GMapMinWidth        =   GMapSreenUnit
	GMapMinHeight       =   GMapSreenUnit
	
	// a unit of the minimal object, which is used to calculate some collision detections
	// but only to check the collision between roles and buildings
    GMapUnit            =   20
)

type GameMap struct {
	Name				string
	Width				int
	Height				int

	ScreenXCount		int
	ScreenYCount		int
	ScreenCount			int

	Screens				[]map[int64]*GameClient

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
	index := i + j * gameMap.ScreenXCount
	return index
}

/*
 * add a reflect between game client and the map
 */
func (gameMap *GameMap) ChangeScreen(gameClient *GameClient) {
	index := gameMap.GetScreenIndex(gameClient.GameData.X, gameClient.GameData.Y)
	if index != gameClient.GameData.ScreenId {
		delete(gameMap.Screens[gameClient.GameData.ScreenId], gameClient.ID)
	}
	gameMap.Screens[index][gameClient.ID] = gameClient
}

/*
 * remove the screen reflect
 */
func (gameMap *GameMap) RemoveScreen(gameClient *GameClient) {
	delete(gameMap.Screens[gameClient.GameData.ScreenId], gameClient.ID)
	gameMap.BroadCast9(gameClient, nil)
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
 * broadcast the `order` to the clients of the map
 */
func (gameMap *GameMap) BroadCast(order interface{}) {
	for i := 0; i < gameMap.ScreenCount; i++ {
		gameMap.BroadCast1(i, order)
	}
}

/*
 * broadcast the `order` to the clients of the nine screens
 */
func (gameMap *GameMap) BroadCast9(gameClient *GameClient, order interface{}) {
	screenId := gameClient.GameData.ScreenId
	i := screenId % gameMap.ScreenXCount

	gameMap.BroadCast1(screenId, order)

	topId := screenId - gameMap.ScreenXCount
	if (topId >= 0) {
		gameMap.BroadCast1(topId, order)
		if i > 0 {
			// left and left-top
			gameMap.BroadCast1(screenId - 1, order)
			gameMap.BroadCast1(topId - 1, order)
		} else if i < (gameMap.ScreenXCount - 1) {
			// right and right-top
			gameMap.BroadCast1(screenId + 1, order)
			gameMap.BroadCast1(topId + 1, order)
		}
	}

	bottomId := screenId + gameMap.ScreenXCount
	if (bottomId < gameMap.ScreenCount) {
		gameMap.BroadCast1(bottomId, order)
		if i > 0 {
			// left-bottom
			gameMap.BroadCast1(bottomId - 1, order)
		} else if i < (gameMap.ScreenXCount - 1) {
			// right-bottom
			gameMap.BroadCast1(bottomId + 1, order)
		}
	}
	
}

/*
 * broadcast order in a screen
 */
func (gameMap *GameMap) BroadCast1(screenId int, order interface{}) {
	for _, client := range gameMap.Screens[screenId] {
		client.Conn.WriteJSON(order)
	}
}