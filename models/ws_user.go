package models

import(
	"github.com/gorilla/websocket"
)

type WsUser struct {
	ID  	int64
	Conn	*websocket.Conn
	WsData	Game0
}