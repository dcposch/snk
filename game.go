package main

import (
	"math/rand"
	"time"
)

// Vec is a basic (x,y) int vector. Supports addition and ==
type Vec [2]int

func (a Vec) Plus(b Vec) Vec {
	return Vec{a[0] + b[0], a[1] + b[1]}
}

// SnakeGame stores the state of the board and current score.
// Implements the rules. Not responsible for user input or output.
type SnakeGame struct {
	width, height int
	snake         []Vec
	food          []Vec
	vel           Vec // snake velocity: up, down, left, or right
	isOver        bool
	score         int
}

func CreateSnakeGame(width, height int) *SnakeGame {
	// snake tail starts at 0,0, with the snake moving down
	snake := make([]Vec, 4)
	for i := 0; i < len(snake); i++ {
		snake[i] = Vec{0, i}
	}
	vel := Vec{0, 1}

	// place random food pieces, not overlapping each other or the snake
	food := make([]Vec, 2)

	ret := &SnakeGame{
		width:  width,
		height: height,
		snake:  snake,
		food:   food,
		vel:    vel,
	}

	for i := range food {
		food[i] = ret.placeRandomFood()
	}

	return ret
}

// SetVel sets the new snake velocity, in squares per tick.
// The only allowed directions are up, down (0, 1), left, and right (1, 0).
func (game *SnakeGame) SetVel(vel Vec) {
	if vel.Plus(game.vel) == (Vec{0, 0}) {
		return // ignore direction reversal
	}
	game.vel = vel
}

// Step advances the game by one step. Moves the snake, etc.
func (game *SnakeGame) Step() {
	if game.isOver {
		return
	}

	// find the new head of the snake
	tip := game.snake[len(game.snake)-1]
	nt := tip.Plus(game.vel)
	// bounds check
	if nt[0] < 0 || nt[0] >= game.width || nt[1] < 0 || nt[1] >= game.height {
		game.isOver = true
		return
	}
	// collision check
	if indexOf(nt, game.snake) >= 0 {
		game.isOver = true
		return
	}

	foodIndex := indexOf(nt, game.food)
	if foodIndex < 0 {
		// rotate in place, snake length remains the same
		for i := 0; i < len(game.snake)-1; i++ {
			game.snake[i] = game.snake[i+1]
		}
		game.snake[len(game.snake)-1] = nt
	} else {
		// append tip to snake
		game.snake = append(game.snake, nt)
		// eat the food, place new food at random location
		game.food[foodIndex] = game.placeRandomFood()
		// incrmeent score
		game.score++
	}
}

// Gets the current game tick speed. The game speeds up as the score increases.
func (game *SnakeGame) GetTickDuration() time.Duration {
	return time.Second * 3 / time.Duration(10+game.score/3)
}

func (game *SnakeGame) placeRandomFood() (ret Vec) {
	for {
		ret[0] = rand.Int() % game.width
		ret[1] = rand.Int() % game.height
		if indexOf(ret, game.snake) < 0 && indexOf(ret, game.food) < 0 {
			return
		}
	}
}

func indexOf(loc Vec, arr []Vec) int {
	for i := 0; i < len(arr); i++ {
		if arr[i] == loc {
			return i
		}
	}
	return -1
}
