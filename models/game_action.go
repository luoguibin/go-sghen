package models

type GameAction struct {
	Action 		string		`json:"action"`
	Target		int64		`json:"target"`
	Msg         string		`json:"msg"`
}