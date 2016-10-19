package main

import "testing"
import tl "github.com/JoelOtter/termloop"

func Test(t *testing.T) {
	game := tl.NewGame()
	level := tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorYellow,
		Fg: tl.ColorBlack,
		Ch: ' ',
	})
	panic(1)
	IniGame(level, game)
	game.Screen().SetLevel(level)
	GetInstance().prevPlaceX = 0
	GetInstance().prevPlaceY = 0
	GetInstance().nests[0][0] = 9
	GetInstance().nests[0][3] = 6
	GetInstance().nests[0][6] = 6
	if GetInstance().CheckThird() == 1 {

	} else {
		t.Error("NE OK")
	}

}
