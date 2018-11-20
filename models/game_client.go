package models

import(
	"github.com/gorilla/websocket"
)

type GameClient struct {
	ID  			int64
	Conn			*websocket.Conn
	GameData		*GameData
	GameStatus		int
}

var (
	StatusLogin 	=	0
	StatusInGame	=   1
	StatusLogout	= 	-1
)