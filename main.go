package main

import (
	"fmt"
	"os"
	"time"

	"github.com/nsf/termbox-go"
)

func main() {
	if err := termbox.Init(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer termbox.Close()

	evChan := make(chan termbox.Event)
	go func() {
		for {
			evChan <- termbox.PollEvent()
		}
	}()

	game := NewGame()
	go func() {
		for {
			select {
			case <-time.After(16 * time.Millisecond):
				UpdateScreen(game)

			case ev := <-evChan:
				switch ev.Type {
				case termbox.EventKey:
					if ev.Key == termbox.KeyEsc || (ev.Key == termbox.KeyCtrlC && ev.Mod == 0) {
						gracefulStop()
						os.Exit(0)
					}
					HandleKey(ev, game)
				}
			}
		}
	}()

	for {
		time.Sleep(1 * time.Second)
	}
}

func HandleKey(ev termbox.Event, game *Game_t) {
	switch ev.Key {
	case termbox.KeyArrowUp:
		game.Rotate(game.MinoReal)

	case termbox.KeyArrowRight:
		if game.MinoReal != nil && game.PermittedMoves(game.MinoReal)["right"] {
			game.MinoReal.position.x += 1
		}

	case termbox.KeyArrowLeft:
		if game.MinoReal != nil && game.PermittedMoves(game.MinoReal)["left"] {
			game.MinoReal.position.x -= 1
		}

	case termbox.KeyArrowDown:
		if game.MinoReal != nil {
			if game.PermittedMoves(game.MinoReal)["down"] {
				game.lastMoveTime = time.Now()
				game.MinoReal.position.y += 1
			} else {
				game.lastMoveTime = time.Time{}.Add(1)
			}
		}

	case termbox.KeySpace:
		if game.MinoReal != nil {
			if game.PermittedMoves(game.MinoReal)["down"] {
				_, ghost := game.MaxFall(game.MinoReal)
				game.MinoReal = ghost
				game.lastMoveTime = time.Time{}.Add(1)
			}
		}

	case termbox.KeyEnter:
		if game.MinoReal != nil {
			game.ReserveMino()
		}
	}
}

func UpdateScreen(game *Game_t) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	defer termbox.Flush()

	w, h := termbox.Size()
	if w < 20 || h < 40 {
		tbprint(0, 0, termbox.ColorRed, termbox.ColorDefault, "Terminal is too small")
		return
	}
	game.Update()
}
