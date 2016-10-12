package main

import tl "github.com/JoelOtter/termloop"

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
}

var instance *singleton

func GetInstance() *singleton {
	if instance == nil {
		instance = &singleton{}
	}
	return instance
}

func (p *Player) EnterHandle() {
	gameMode := GetInstance()
	if gameMode.startCount != 18 {

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

			(*gameMode.level).AddEntity(entity)
		}
		gameMode.startCount += 1
		gameMode.Turn = !(GetInstance().Turn)
		p.ChangeColor()
		return
	}

}

func (p *Player) ChangeColor() {
	if GetInstance().Turn == true {
		p.entity.SetCell(0, 0, &tl.Cell{Fg: tl.ColorWhite, Ch: 'X'})
	} else {
		p.entity.SetCell(0, 0, &tl.Cell{Fg: tl.ColorBlack, Ch: 'X'})
	}

}
