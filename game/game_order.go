package game

type GameOrder struct {
	// who makes the `order`
	FromType int   `json:"fromType"`
	FromID   int64 `json:"fromId"`

	// the `order` type
	OrderType int `json:"order"`

	// the `order` data than need to be translated and executed,
	// use `mapstructure.Decode`, and should make the json label letter same as defined ignore uppercase
	Data interface{} `json:"data"`
}

type GameOrderMsg struct {
	ToType   int    `json:"toType"`
	ToID     int64  `json:"toId"`
	FromName string `json:"fromName,omitempty"`
	Msg      string `json:"msg"`
}

type GameOrderSkill struct {
	ToID int64 `json:"toId"`
	// Targets			[]int64		`json:"targets"`
	Damage         int `json:"damage"`
	DamageAll      int `json:"damageAll"`
	DamageCount    int `json:"damageCount,omitempty"`
	DamageCountAll int `json:"damageCountAll,omitempty"`
}

type GameOrderAction struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}
