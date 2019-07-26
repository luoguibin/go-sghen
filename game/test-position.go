package game

import (
	"fmt"
	"math"
	"time"
)

// Order ...
type Order struct {
	X float64
	Y float64
}

// Player ...
type Player struct {
	Order *Order
	X     float64
	Y     float64
}

var (
	//MPlayer ...
	MPlayer *Player

	// MStep ...
	MStep = 3.0

	// MBoxSize ...
	MBoxSize = MStep * 10.0
)

// test for the box index of moving position
func positionService() {
	for {
		time.Sleep(time.Second)

		if MPlayer.Order != nil {
			tempX := MPlayer.Order.X - MPlayer.X
			tempY := MPlayer.Order.Y - MPlayer.Y

			tempVal := math.Sqrt(math.Pow(tempX, 2) + math.Pow(tempY, 2))
			if tempVal <= MStep {
				MPlayer.Order = nil
			} else {
				ratio := MStep / tempVal
				tempX = tempX * ratio
				tempY = tempY * ratio
			}

			boxTempX := math.Floor(MPlayer.X / MBoxSize)
			boxTempY := math.Floor(MPlayer.Y / MBoxSize)

			MPlayer.X = MPlayer.X + tempX
			MPlayer.Y = MPlayer.Y + tempY

			boxX := math.Floor(MPlayer.X / MBoxSize)
			boxY := math.Floor(MPlayer.Y / MBoxSize)

			fmt.Print(MPlayer, " ")
			if boxX != boxTempX || boxY != boxTempY {
				fmt.Print(boxTempX, boxTempY, ",", boxX, boxY)
			}
			fmt.Println()
		}
	}
}

func main() {
	MPlayer = &Player{
		X: 0,
		Y: 0,
	}
	fmt.Println(MPlayer)

	go positionService()

	for {
		var str string
		fmt.Scan(&str)
		switch str {
		case "move":
			var x float64
			var y float64

			fmt.Scan(&x)
			fmt.Scan(&y)

			order := &Order{
				X: x,
				Y: y,
			}
			MPlayer.Order = order
		}
	}
}
