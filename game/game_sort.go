package game

type GameSortItem struct {
	Value 			interface{}
	GameClient		*GameClient
}

type GameSort []*GameSortItem

func (s GameSort) Len() int { 
	return len(s) 
} 

func (s GameSort) Swap(i, j int) { 
	s[i], s[j] = s[j], s[i] 
} 

func (s GameSort) Less(i, j int) bool {
	switch s[i].Value.(type) {
	case int:
		return s[i].Value.(int) < s[j].Value.(int)
	case int64:
		return s[i].Value.(int64) < s[j].Value.(int64)
	case float64:
		return s[i].Value.(float64) < s[j].Value.(float64)
	default:
		return false
	}
}
