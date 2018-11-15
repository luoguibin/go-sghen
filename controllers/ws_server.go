package controllers

import(
	"SghenApi/models"
	"time"
	"log"
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
)

type WSServerController struct {
	BaseController
}

func init() {
    go handleMessages()
}

var (
	upgrader 	= websocket.Upgrader{
		ReadBufferSize: 	2048,
		WriteBufferSize: 	2048,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	clients   	= make(map[*websocket.Conn]bool)
	broadcast	= make(chan models.WSMessage)
)


func (c *WSServerController) Get() { 
	ws, err := upgrader.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil) 
	if err != nil { 
		log.Fatal(err) 
	} 
	// defer ws.Close() 
	clients[ws] = true 
	//不断的广播发送到页面上 
	for { 
		//目前存在问题 定时效果不好 需要在业务代码替换时改为beego toolbox中的定时器 
		time.Sleep(time.Second * 3) 
		msg := models.WSMessage{Message: "这是向页面发送的数据 " + time.Now().Format("2006-01-02 15:04:05")} 
		broadcast <- msg
	} 
}

//广播发送至页面 
func handleMessages() { 
	for { 
		msg := <-broadcast 
		fmt.Println("clients len ", len(clients)) 
		for client := range clients { 
			err := client.WriteJSON(msg) 
			if err != nil { 
				log.Printf("client.WriteJSON error: %v", err) 
				client.Close() 
				delete(clients, client) 
			} 
		} 
	} 
}
