package controllers

import (
	"go-sghen/game"
	"strconv"
)

type GameController struct {
	BaseController
}

/**
 * WebSocket连接入口
 * 在BeforeRouter检测jwt中的合法后才给予长连接
 */
func (c *GameController) Get() {
	uId, _ := strconv.ParseInt(c.Ctx.Input.Query("uId"), 10, 64)
	game.AddToServer(c.Ctx, uId)
}
