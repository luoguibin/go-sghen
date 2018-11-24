package models

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
	// who makes the `order`
	FromType	int			`json:"fromType"`
	FromID		int64		`json:"fromId"`

	// the `order` type
	OrderType 	int			`json:"order"`

	// the `order` data than need to be translated and executed
	Data		interface{}	`json:"data"`
}

type GameOrderData struct {
	Orders		[]interface{}	`json:"orders"`
	Data		[]interface{}	`json:"data"`
}

type GameOrderMsg struct {
	ToType		int 		`json:"toType"`
	ToID		int64		`json:"toId"`
	Msg			string		`json:"msg"`
}

type GameOrderSkill struct {
	ToID			int64 		`json:"toId"`
	SkillID			int			`json:"skillId"`
	Targets			[]int64		`json:"targets"`
	Damage			int 		`json:"damage"`
	DamageAll 		int 		`json:"damageAll"`
	DamageCount		int			`json:"damageCount,omitempty"`
	DamageCountAll	int			`json:"damageCountAll,omitempty"`
}

type GameAction struct {
	Action  int     `json:"action"`
	GX		int 	`json:"x"`
	GY		int 	`json:"y"`
}