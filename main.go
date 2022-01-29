package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
)

func main() {
	// Create game
	game := CreateSnakeGame(13, 13)

	// Create canvas
	s, err := tcell.NewScreen()
	must(err)
	must(s.Init())

	// Set default text style
	styleBG := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	styleBold := styleBG.Bold(true)
	styleSnake := styleBG.Background(tcell.ColorLightGreen)
	styleFood := styleBG.Background(tcell.ColorTeal)
	styleGameOver := styleBold.Background(tcell.ColorGray).Foreground(tcell.ColorWhite)
	s.SetStyle(styleBG)

	// Clear screen
	s.Clear()

	// Event loop
	events := make(chan tcell.Event)
	go func() {
		for {
			events <- s.PollEvent()
		}
	}()
	ticker := time.NewTicker(game.GetTickDuration())

	// Game loop
	moveQ := []Vec{}
	for i := 0; ; i++ {
		// Render
		s.Clear()

		// ...frame
		gw := game.width
		gh := game.height
		for i := 0; i < gw*2+1; i++ {
			s.SetContent(i, 1, '-', nil, styleBG)
			s.SetContent(i, gh+2, '-', nil, styleBG)
		}
		for i := 1; i < gh+2; i++ {
			s.SetContent(0, i, '|', nil, styleBG)
			s.SetContent(gw*2+1, i, '|', nil, styleBG)
		}
		s.SetContent(0, 1, '+', nil, styleBG)
		s.SetContent(0, gh+2, '+', nil, styleBG)
		s.SetContent(gw*2+1, 1, '+', nil, styleBG)
		s.SetContent(gw*2+1, gh+2, '+', nil, styleBG)

		// ...score
		setText(s, gw*2-5, gh+3, fmt.Sprintf("%06d", game.score), styleBold)
		setText(s, 1, gh+3, "SCORE", styleBold)

		// ...snake and food
		for _, vec := range game.snake {
			s.SetContent(vec[0]*2+1, vec[1]+2, ' ', nil, styleSnake)
			s.SetContent(vec[0]*2+2, vec[1]+2, ' ', nil, styleSnake)
		}
		for _, vec := range game.food {
			s.SetContent(vec[0]*2+1, vec[1]+2, ' ', nil, styleFood)
			s.SetContent(vec[0]*2+2, vec[1]+2, ' ', nil, styleFood)
		}

		// ...game over
		if game.isOver {
			setText(s, gw-6, gh/2+1, "              ", styleGameOver)
			setText(s, gw-6, gh/2+2, "  GAME  OVER  ", styleGameOver)
			setText(s, gw-6, gh/2+3, "              ", styleGameOver)
		}

		// Update screen
		s.Show()

		// Poll event
		select {
		case <-ticker.C:
			origVel := game.vel
			for len(moveQ) > 0 && game.vel == origVel {
				game.SetVel(moveQ[0])
				copy(moveQ, moveQ[1:])
				moveQ = moveQ[:len(moveQ)-1]
			}
			game.Step()
			ticker.Reset(game.GetTickDuration())
		case ev := <-events:
			switch e := ev.(type) {
			case *tcell.EventResize:
				s.Sync()
			case *tcell.EventKey:
				switch e.Key() {
				case tcell.KeyEsc:
				case tcell.KeyCtrlC:
					s.Fini()
					return
				case tcell.KeyUp:
					moveQ = append(moveQ, Vec{0, -1})
				case tcell.KeyDown:
					moveQ = append(moveQ, Vec{0, 1})
				case tcell.KeyLeft:
					moveQ = append(moveQ, Vec{-1, 0})
				case tcell.KeyRight:
					moveQ = append(moveQ, Vec{1, 0})
				case tcell.KeyRune:
					switch e.Rune() {
					case 'w':
						moveQ = append(moveQ, Vec{0, -1})
					case 's':
						moveQ = append(moveQ, Vec{0, 1})
					case 'a':
						moveQ = append(moveQ, Vec{-1, 0})
					case 'd':
						moveQ = append(moveQ, Vec{1, 0})
					}
				}
			}
		}
	}
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
