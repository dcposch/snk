package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
)

var (
	styleBG       = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	styleBold     = styleBG.Bold(true)
	styleSnake    = styleBG.Background(tcell.ColorLightGreen)
	styleFood     = styleBG.Background(tcell.ColorTeal)
	styleGameOver = styleBold.Background(tcell.ColorGray).Foreground(tcell.ColorWhite)
	gw            = 13
	gh            = 13
)

func main() {
	// Create game
	game := CreateSnakeGame(gw, gh)

	// Create canvas
	s, err := tcell.NewScreen()
	must(err)
	must(s.Init())
	s.SetStyle(styleBG)

	// Poll events, allowing a single synchronous game loop using select{}
	events := make(chan tcell.Event)
	go func() {
		for {
			events <- s.PollEvent()
		}
	}()

	// Game loop
	ticker := time.NewTicker(game.GetTickDuration())
	moveQ := []Vec{}
	for {
		// Render
		s.Clear()
		drawFrame(s)
		drawScore(s, game.score)
		drawGame(s, game)
		s.Show()

		// Poll event
		select {
		case <-ticker.C:
			// game tick. start by turning the snake, if applicable.
			// to make the game more fun, player can queue up a few moves
			// that run on subsequent game ticks
			origVel := game.vel
			for len(moveQ) > 0 && game.vel == origVel {
				game.SetVel(moveQ[0])
				copy(moveQ, moveQ[1:])
				moveQ = moveQ[:len(moveQ)-1]
			}
			// next, advance the game by one turn
			game.Step()
			// speed up the game tick as the score increases
			ticker.Reset(game.GetTickDuration())

		case ev := <-events:
			// no game gick. instead, just deal with events
			switch e := ev.(type) {
			case *tcell.EventResize:
				s.Sync()
			case *tcell.EventKey:
				switch e.Key() {
				case tcell.KeyEsc:
				case tcell.KeyCtrlC:
					// exit
					s.Fini()
					return
				default:
					// player turns the snake at next tick, tick after etc.
					moveQ = append(moveQ, getMove(e))
				}
			}
		}
	}
}

func getMove(e *tcell.EventKey) Vec {
	switch e.Key() {
	case tcell.KeyUp:
		return Vec{0, -1}
	case tcell.KeyDown:
		return Vec{0, 1}
	case tcell.KeyLeft:
		return Vec{-1, 0}
	case tcell.KeyRight:
		return Vec{1, 0}
	case tcell.KeyRune:
		switch e.Rune() {
		case 'w':
			return Vec{0, -1}
		case 's':
			return Vec{0, 1}
		case 'a':
			return Vec{-1, 0}
		case 'd':
			return Vec{1, 0}
		}
	}
	return Vec{0, 0}
}

func drawFrame(s tcell.Screen) {
	ox, oy := getCenterOffset(s)
	for i := 0; i < gw*2+1; i++ {
		s.SetContent(ox+i, oy+1, '-', nil, styleBG)
		s.SetContent(ox+i, oy+gh+2, '-', nil, styleBG)
	}
	for i := 1; i < gh+2; i++ {
		s.SetContent(ox, oy+i, '|', nil, styleBG)
		s.SetContent(ox+gw*2+1, oy+i, '|', nil, styleBG)
	}
	s.SetContent(ox, oy+1, '+', nil, styleBG)
	s.SetContent(ox, oy+gh+2, '+', nil, styleBG)
	s.SetContent(ox+gw*2+1, oy+1, '+', nil, styleBG)
	s.SetContent(ox+gw*2+1, oy+gh+2, '+', nil, styleBG)
}

func drawScore(s tcell.Screen, score int) {
	ox, oy := getCenterOffset(s)
	setText(s, ox+gw*2-5, oy+gh+3, fmt.Sprintf("%06d", score), styleBold)
	setText(s, ox+1, oy+gh+3, "SCORE", styleBold)
}

func drawGame(s tcell.Screen, game *SnakeGame) {
	ox, oy := getCenterOffset(s)

	// ...snake and food
	for _, vec := range game.snake {
		s.SetContent(ox+vec[0]*2+1, oy+vec[1]+2, ' ', nil, styleSnake)
		s.SetContent(ox+vec[0]*2+2, oy+vec[1]+2, ' ', nil, styleSnake)
	}
	for _, vec := range game.food {
		s.SetContent(ox+vec[0]*2+1, oy+vec[1]+2, ' ', nil, styleFood)
		s.SetContent(ox+vec[0]*2+2, oy+vec[1]+2, ' ', nil, styleFood)
	}

	// ...game over
	if game.isOver {
		setText(s, ox+gw-6, oy+gh/2+1, "              ", styleGameOver)
		setText(s, ox+gw-6, oy+gh/2+2, "  GAME  OVER  ", styleGameOver)
		setText(s, ox+gw-6, oy+gh/2+3, "              ", styleGameOver)
	}
}

func getCenterOffset(s tcell.Screen) (ox, oy int) {
	sx, sy := s.Size()
	ox = sx/2 - gw
	oy = (sy - gh) / 2
	return
}

func setText(s tcell.Screen, x, y int, text string, style tcell.Style) {
	for i := 0; i < len(text); i++ {
		s.SetContent(x+i, y, rune(text[i]), nil, style)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
