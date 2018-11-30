package game

import(
	"SghenApi/models"
	"github.com/gorilla/websocket"
)

type GameClient struct {
	ID  			int64
	Conn			*websocket.Conn
	GameData		*models.GameData
	GMap			*GameMapService
	GameStatus		int
}