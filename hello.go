package main

import tl "github.com/JoelOtter/termloop"
import "math"

func main() {

	game := tl.NewGame()
	level := tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorYellow,
		Fg: tl.ColorBlack,
		Ch: ' ',
	})

	IniGame(level, game)
	game.Screen().SetLevel(level)
	game.Start()

}

type Player struct {
	entity *tl.Entity
}

// Here, Draw simply tells the Entity ent to handle its own drawing.
// We don't need to do anything.
func (player *Player) Draw(screen *tl.Screen) {
	player.entity.Draw(screen)
}

func (player *Player) Tick(event tl.Event) {
	if event.Type == tl.EventKey { // Is it a keyboard event?
		x, y := player.entity.Position()
		switch event.Key { // If so, switch on the pressed key.
		case tl.KeyArrowRight:
			player.entity.SetPosition(x+1, y)
		case tl.KeyArrowLeft:
			player.entity.SetPosition(x-1, y)
		case tl.KeyArrowUp:
			player.entity.SetPosition(x, y-1)
		case tl.KeyArrowDown:
			player.entity.SetPosition(x, y+1)
		case tl.KeyEnter:
			player.EnterHandle()
		case tl.KeyCtrlA:
			panic(GetInstance().startCount)

		}
	}
}

func IniGame(level *tl.BaseLevel, game *tl.Game) {
	GetInstance().nests = [7][7]int{
		{1, 0, 0, 1, 0, 0, 1},
		{0, 1, 0, 1, 0, 1, 0},
		{0, 0, 1, 1, 1, 0, 0},
		{1, 1, 1, 2, 1, 1, 1},
		{0, 0, 1, 1, 1, 0, 0},
		{0, 1, 0, 1, 0, 1, 0},
		{1, 0, 0, 1, 0, 0, 1}}
	nests := &GetInstance().nests
	GetInstance().level = level

	(*level).AddEntity(tl.NewRectangle(0, 0, 7, 1, tl.ColorBlue)) // -----
	(*level).AddEntity(tl.NewRectangle(2, 1, 3, 1, tl.ColorBlue)) //  ---
	//
	(*level).AddEntity(tl.NewRectangle(2, 5, 3, 1, tl.ColorBlue)) //  ---
	(*level).AddEntity(tl.NewRectangle(0, 6, 7, 1, tl.ColorBlue)) // -----
	//same but verticallly
	(*level).AddEntity(tl.NewRectangle(0, 0, 1, 7, tl.ColorBlue)) // -----
	(*level).AddEntity(tl.NewRectangle(1, 2, 1, 3, tl.ColorBlue)) //  ---
	//
	(*level).AddEntity(tl.NewRectangle(5, 2, 1, 3, tl.ColorBlue)) //  ---
	(*level).AddEntity(tl.NewRectangle(6, 0, 1, 7, tl.ColorBlue)) // -----

	for i := 0; i < 7; i++ {
		for j := 0; j < 7; j++ {
			if nests[i][j] == 1 {
				e := tl.NewEntity(i, j, 1, 1)
				e.SetCell(0, 0, &tl.Cell{Fg: tl.ColorRed, Ch: 'O'})
				(*level).AddEntity(e)
			}
		}
	}
	player := Player{
		entity: tl.NewEntity(1, 1, 1, 1),
	}
	player.entity.SetCell(0, 0, &tl.Cell{Fg: tl.ColorRed, Ch: 'X'})
	(*level).AddEntity(&player)
	(*game).Screen().SetLevel(level)

}
func (p *Player) ChipsCount() int {
	var side int
	if GetInstance().Turn {
		side = 6
	} else {
		side = 9
	}
	count := 0
	for x := 0; x < 7; x++ {
		for y := 0; y < 7; y++ {
			if GetInstance().nests[x][y] == side {
				count++
			}
		}
	}
	return count
}

type singleton struct {
	Turn       bool
	Fg         tl.Attr
	startCount int
	nests      [7][7]int
	isRemoving bool
	isMoving   bool
	xMoving    int
	yMoving    int
	level      *tl.BaseLevel
	prevPlaceX int
	prevPlaceY int
	countRemov int
}

var instance *singleton

func GetInstance() *singleton {
	if instance == nil {
		instance = &singleton{}
		instance.startCount = 0
	}
	return instance
}

func (gm *singleton) CheckThird() int {
	nests := GetInstance().nests
	var side int
	if gm.Turn {
		side = 6
	} else {
		side = 9
	}
	actualCounter := 0
	counter := 0

	for y := 0; y < 7; y++ {
		if nests[gm.prevPlaceX][y] == 2 {
			counter = 0
		} else if nests[gm.prevPlaceX][y] == 0 {
			continue
		} else if nests[gm.prevPlaceX][y] != side {
			counter = 0
		} else if nests[gm.prevPlaceX][y] == side {
			counter++
		} else {
			panic("in the panthagram")
		}
	}
	if counter == 3 {
		actualCounter += 1
	}
	counter = 0
	for x := 0; x < 7; x++ {
		if nests[x][gm.prevPlaceY] == 2 {
			counter = 0
			panic(counter)
		} else if nests[x][gm.prevPlaceY] == 0 {
			continue
		} else if nests[x][gm.prevPlaceY] != side {
			counter = 0
		} else if nests[x][gm.prevPlaceY] == side {
			counter++
		} else {
			panic("in the panthagram")
		}
	}
	if counter == 3 {
		actualCounter += 1
	}

	return actualCounter

}

func (p *Player) EnterHandle() {
	gameMode := GetInstance()
	if gameMode.isRemoving {
		var antiside int
		if gameMode.Turn {
			antiside = 9
		} else {
			antiside = 6
		}
		x, y := p.entity.Position()
		if x < 0 || y < 0 || x > 7 || y > 7 || gameMode.nests[x][y] != antiside {
			return
		}
		gameMode.nests[x][y] = 1
		e := tl.NewEntity(x, y, 1, 1)
		e.SetCell(0, 0, &tl.Cell{Fg: tl.ColorRed, Ch: 'O'})
		(*gameMode.level).AddEntity(e)
		gameMode.countRemov--
		if gameMode.countRemov != 0 {
			return
		} else {
			gameMode.isRemoving = false
		}

	} else if gameMode.startCount < 18 {
		x, y := p.entity.Position()
		if x < 0 || y < 0 || x > 7 || y > 7 || gameMode.nests[x][y] != 1 {
			return
		} else {
			c := tl.Cell{Ch: ' '}
			if gameMode.Turn == true {
				gameMode.nests[x][y] = 6
				c.Bg = tl.ColorWhite
			} else {

				gameMode.nests[x][y] = 9
				c.Bg = tl.ColorBlack
			}
			entity := tl.NewEntity(x, y, 1, 1)

			entity.SetCell(0, 0, &c)
			gameMode.prevPlaceX = x
			gameMode.prevPlaceY = y
			(*gameMode.level).AddEntity(entity)
			cnt := gameMode.CheckThird()
			if cnt > 0 {
				gameMode.isRemoving = true
				gameMode.countRemov = cnt
				gameMode.startCount += 1
				return
			}
			gameMode.startCount += 1
		}
	} else if gameMode.isMoving {
		if p.ChipsCount() == 3 {
			x, y := p.entity.Position()
			if x < 0 || y < 0 || x > 7 || y > 7 || gameMode.nests[x][y] != 1 {
				return
			}
			gameMode.nests[x][y] = gameMode.nests[gameMode.prevPlaceX][gameMode.prevPlaceY]
			gameMode.nests[gameMode.prevPlaceX][gameMode.prevPlaceY] = 1
			e := tl.NewEntity(x, y, 1, 1)
			var ch tl.Attr
			if gameMode.Turn {
				ch = tl.ColorWhite
			} else {
				ch = tl.ColorBlack
			}
			e.SetCell(0, 0, &tl.Cell{Bg: ch, Ch: ' '})
			(*gameMode.level).AddEntity(e)
			ent := tl.NewEntity(gameMode.prevPlaceX, gameMode.prevPlaceY, 1, 1)
			ent.SetCell(0, 0, &tl.Cell{Bg: ch, Ch: 'O'})
			(*gameMode.level).AddEntity(ent)
			gameMode.isMoving = false
			gameMode.prevPlaceX = x
			gameMode.prevPlaceY = y
			cnt := gameMode.CheckThird()

			if cnt > 0 {
				gameMode.isRemoving = true
				gameMode.countRemov = cnt
				gameMode.startCount += 1
				return
			}
			p.ChangeColor()
			return
		}
		x, y := p.entity.Position()
		if x < 0 || y < 0 || x > 7 || y > 7 || gameMode.nests[x][y] != 1 {
			return
		}

		if (math.Abs(float64(x-gameMode.prevPlaceX)) == 2 && y == gameMode.prevPlaceY && gameMode.nests[(x+gameMode.prevPlaceX)/2][y] != 1) ||
			(math.Abs(float64(y-gameMode.prevPlaceY)) == 2 && x == gameMode.prevPlaceX && gameMode.nests[x][(y+gameMode.prevPlaceY)/2] != 1) ||
			(math.Abs(float64(x-gameMode.prevPlaceX)) == 3 && y == gameMode.prevPlaceY && gameMode.nests[(x+gameMode.prevPlaceX)/2+1][y] != 1) ||
			(math.Abs(float64(y-gameMode.prevPlaceY)) == 3 && x == gameMode.prevPlaceX && gameMode.nests[x][(y+gameMode.prevPlaceY)/2+1] != 1) ||
			(math.Abs(float64(y-gameMode.prevPlaceY)) == 1 && x == gameMode.prevPlaceX) ||
			(math.Abs(float64(x-gameMode.prevPlaceX)) == 1 && y == gameMode.prevPlaceY) {

			gameMode.nests[x][y] = gameMode.nests[gameMode.prevPlaceX][gameMode.prevPlaceY]
			gameMode.nests[gameMode.prevPlaceX][gameMode.prevPlaceY] = 1
			e := tl.NewEntity(x, y, 1, 1)
			var ch tl.Attr
			if gameMode.Turn {
				ch = tl.ColorWhite
			} else {
				ch = tl.ColorBlack
			}
			e.SetCell(0, 0, &tl.Cell{Bg: ch, Ch: ' '})
			(*gameMode.level).AddEntity(e)
			ent := tl.NewEntity(gameMode.prevPlaceX, gameMode.prevPlaceY, 1, 1)
			ent.SetCell(0, 0, &tl.Cell{Bg: ch, Ch: 'O'})
			(*gameMode.level).AddEntity(ent)
			gameMode.isMoving = false
			gameMode.prevPlaceX = x
			gameMode.prevPlaceY = y
			cnt := gameMode.CheckThird()

			if cnt > 0 {
				gameMode.isRemoving = true
				gameMode.countRemov = cnt
				gameMode.startCount += 1
				return
			}
		} else {
			return
		}
	} else {
		if p.ChipsCount() < 3 {

			(*GetInstance().level).AddEntity(tl.NewText(0, 0, "Well played", tl.ColorBlack, tl.ColorGreen))
		}
		var side int

		if GetInstance().Turn {
			side = 6
		} else {
			side = 9
		}
		x, y := p.entity.Position()
		if gameMode.nests[x][y] != side {
			return
		} else {
			gameMode.isMoving = true
			gameMode.prevPlaceX = x
			gameMode.prevPlaceY = y
			return
		}
	}

	p.ChangeColor()
	return

}

func (p *Player) ChangeColor() {
	GetInstance().Turn = !(GetInstance().Turn)

	if GetInstance().Turn == true {
		p.entity.SetCell(0, 0, &tl.Cell{Fg: tl.ColorWhite, Ch: 'X'})
	} else {
		p.entity.SetCell(0, 0, &tl.Cell{Fg: tl.ColorBlack, Ch: 'X'})
	}

}
