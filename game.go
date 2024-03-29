package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/nsf/termbox-go"
)

type Game_t struct {
	Score              int
	Streak             int
	moveInterval       time.Duration
	lastMoveTime       time.Time
	lastSpeedIncrement time.Time
	gamestartTime      time.Time
	board              Board
	MinoReal           *Mino
	MinoGhost          *Mino
	MinoReserved       *Mino
	MinoFuture         *Mino
	ResEN              bool
}

func NewGame() *Game_t {
	return &Game_t{
		Score:              0,
		Streak:             0,
		moveInterval:       1000 * time.Millisecond,
		lastMoveTime:       time.Now(),
		lastSpeedIncrement: time.Now(),
		gamestartTime:      time.Now(),
		board:              *NewBoard(),
		MinoReal:           nil,
		MinoGhost:          nil,
		MinoReserved:       nil,
		MinoFuture:         nil,
		ResEN:              false,
	}
}

func (g *Game_t) ReserveMino() {
	if g.MinoReal == nil {
		return
	}
	if !g.ResEN && g.MinoReserved == nil {
		return
	}

	if g.ResEN && g.MinoReserved == nil && g.MinoReal != nil {
		// Store current mino
		g.MinoReserved = g.MinoReal
		g.MinoReal = nil
		return
	}

	if g.ResEN && g.MinoReserved != nil && g.MinoReal != nil {
		// Exchange current mino with reserved one
		t := g.MinoReal
		g.MinoReal = g.MinoReserved
		g.MinoReserved = t
		g.ResEN = false
		return
	}
}

func (g *Game_t) GameOver() {
	tbprint(0, 1, termbox.ColorRed, termbox.ColorBlack, "Game_t over")
	time.Sleep(1 * time.Second)
	gracefulStop()
	os.Exit(0)
}

func (g *Game_t) Update() {
	var err error
	if g.MinoReal == nil {
		g.NewMino()
		g.ResEN = true
	}

	defer func() {
		defer g.board.Draw()
		tbprint(1, 1, termbox.ColorRed, termbox.ColorDefault, fmt.Sprintf("Score: %d, streak: %d, moveInterval: %s", g.Score, g.Streak, g.moveInterval))
		if g.MinoReal == nil {
			tbprint(1, 0, termbox.ColorRed, termbox.ColorDefault, "mino = NIL")
			return
		}
		_, g.MinoGhost = g.MaxFall(g.MinoReal)
		g.board.ProjectMino(g.MinoFuture, g.board.futureMatrix, 0, 0, true)
		g.board.ProjectMino(g.MinoReserved, g.board.reservedMatrix, 0, 0, true)
		g.board.ProjectMino(g.MinoGhost, g.board.ghostMatrix, g.MinoGhost.position.y, g.MinoGhost.position.x, true)
		err = g.board.ProjectMino(g.MinoReal, g.board.upperMatrix, g.MinoReal.position.y, g.MinoReal.position.x, true)
		if err != nil {
			log.Print(err)
			g.GameOver()
		}
	}()

	if time.Since(g.lastMoveTime) > g.moveInterval {
		g.lastMoveTime = time.Now()
		if time.Since(g.lastSpeedIncrement) > 60*time.Second {
			// Descrease moveInterval to make it harder every 60 second(s)
			g.lastSpeedIncrement = time.Now()
			g.moveInterval -= 300 * time.Millisecond
			if g.moveInterval < 200*time.Millisecond {
				g.moveInterval = 200 * time.Millisecond
			}
		}

		if !g.PermittedMoves(g.MinoReal)["down"] {
			// Add stuck blocks to the solid board
			g.board.ProjectMino(g.MinoReal, g.board.solidMatrix, g.MinoReal.position.y, g.MinoReal.position.x, false)
			// Destroy solidified mino in order to allow the creation of a new one later
			g.MinoReal = nil
			// clear eventually complete rows
			g.clearCompleteRows()
		} else {
			// Make the mino fall by one unit
			g.MinoReal.position.y += 1
		}
	}

}

func (g *Game_t) clearCompleteRows() {
	for y := 0; y < len(g.board.solidMatrix); y++ {
		rawCompleted := true
		for x := 0; x < len(g.board.solidMatrix[y]); x++ {
			if g.board.solidMatrix[y][x] == 0 {
				rawCompleted = false
				break
			}
		}
		if rawCompleted {
			g.board.ClearRow(y)
			g.Streak++
			if g.Streak > 5 {
				g.Streak = 5
			}
			g.Score += (g.Streak + 1)
		} else if g.Streak > 1 {
			g.Streak = 0
		}
	}
}

func (g *Game_t) MaxFall(m *Mino) (int, *Mino) {
	maxfall := 0
	shadowMino := Mino{
		shape: m.shape,
		position: &pos_t{
			x: m.position.x,
			y: m.position.y,
		},
		rotation: m.rotation,
		matrix:   m.matrix,
		stuck:    m.stuck,
	}

	for g.PermittedMoves(&shadowMino)["down"] {
		maxfall++
		shadowMino.position.y += 1
	}
	return maxfall, &shadowMino
}

type moves_t map[string]bool

func (g *Game_t) PermittedMoves(mino *Mino) moves_t {
	moves := moves_t{
		"left":  true,
		"right": true,
		"down":  true,
	}
	ny := len(mino.matrix)
	nx := len(mino.matrix[0])

	// tbprint(1, 1, termbox.ColorRed, termbox.ColorDefault, fmt.Sprintf("nx: %d, ny: %d", nx, ny))

	if mino.position.x+nx > board_width_blocks-1 {
		moves["right"] = false
	}
	if mino.position.x < 1 {
		moves["left"] = false
	}
	if mino.position.y+ny > board_height_blocks-1 {
		moves["down"] = false
	}

	for my := 0; my < len(mino.matrix); my++ {
		for mx := 0; mx < len(mino.matrix[my]); mx++ {
			if mino.matrix[my][mx] == 0 {
				continue
			}
			if moves["right"] && g.board.solidMatrix[mino.position.y+my][mino.position.x+mx+1] > 0 {
				moves["right"] = false
			}
			if moves["left"] && g.board.solidMatrix[mino.position.y+my][mino.position.x+mx-1] > 0 {
				moves["left"] = false
			}
			if moves["down"] && g.board.solidMatrix[mino.position.y+my+1][mino.position.x+mx] > 0 {
				moves["down"] = false
			}
		}
	}
	return moves
}

func (g *Game_t) NewMino() {
	g.MinoReal = g.MinoFuture
	m := Mino{
		shape:    mino_t(rand.Intn(7)),
		position: &pos_t{x: 4, y: 0},
		rotation: None,
		matrix:   [][]int{},
		stuck:    false,
	}
	m.setMatrix()
	g.MinoFuture = &m
}

func (g *Game_t) Rotate(mino *Mino) {
	if mino == nil {
		return
	}

	n := len(mino.matrix)
	m := len(mino.matrix[0])
	rotated := make([][]int, m)

	for i := range rotated {
		rotated[i] = make([]int, n)
		for j := range rotated[i] {
			rotated[i][j] = mino.matrix[n-j-1][i]
		}
	}

	newTiles := make([][]int, 4)
	for i := range newTiles {
		newTiles[i] = make([]int, 4)
	}
	for i := range newTiles {
		for j := range newTiles[i] {
			if i > len(rotated)-1 || j > len(rotated[0])-1 {
				newTiles[i][j] = 0
				continue
			}
			if i > len(mino.matrix)-1 || j > len(mino.matrix[0])-1 {
				newTiles[i][j] = rotated[i][j]
				continue
			}
			newTiles[i][j] = int(math.Max(0, float64(rotated[i][j])-float64(mino.matrix[i][j])))
		}
	}

	for i := range newTiles {
		for j := range newTiles[i] {
			if newTiles[i][j] == 1 {
				if mino.position.y+i > board_height_blocks-1 || mino.position.x+j > board_width_blocks-1 {
					// log.Print("No rotation allowed")
					return
				}
				if g.board.solidMatrix[mino.position.y+i][mino.position.x+j] == 1 {
					// log.Print("No rotation allowed")
					return
				}
			}
		}
	}

	mino.matrix = rotated
}
