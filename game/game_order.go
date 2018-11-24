package game


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
	Targets			[]int64		`json:"targets"`
	Damage			int 		`json:"damage"`
	DamageAll 		int 		`json:"damageAll"`
	DamageCount		int			`json:"damageCount,omitempty"`
	DamageCountAll	int			`json:"damageCountAll,omitempty"`
}