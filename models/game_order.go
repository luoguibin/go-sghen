package models

var (
	OrderMsg 		= 	1
	OrderSkill		=	2
	OrderNormal		=	3
)
/*
 * 指令集合：
 *		①消息指令：
 *			消息源：消息产生来源。用ID代表
 *		 	消费对象：消息消费的对象，或一个体，或一群体
 *			消息内容：暂不考虑过滤，表情
 *		②技能指令：
 *			技能源：技能产生来源。用ID代表
 *			消费对象：技能消费的对象，或一个体，或一范围群体
 *			技能：某个技能。用SkillID代表
 *		③综合指令：
 *			
 */

type GameOrder struct {
	OrderType 	int			`json:"order"`
	Target		int64		`json:"target"`
	Msg         string		`json:"msg"`
	Data		string		`json:"data"`
}

type GameAction struct {
	Action  int     `json:"action"`
	GX		int 	`json:"x"`
	GY		int 	`json:"y"`
}