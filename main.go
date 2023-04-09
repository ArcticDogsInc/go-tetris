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

	GameObj := NewGame()

	evChan := make(chan termbox.Event)

	go func() {
		for {
			evChan <- termbox.PollEvent()
		}
	}()

	go func() {
		for {
			select {
			case <-time.After(16 * time.Millisecond):
				DrawScreen(GameObj)
			case ev := <-evChan:
				switch ev.Type {
				case termbox.EventKey:
					if ev.Key == termbox.KeyEsc || (ev.Key == termbox.KeyCtrlC && ev.Mod == 0) {
						gracefulStop()
						os.Exit(0)
					}
					mino := GameObj.Mino
					switch ev.Key {
					case termbox.KeyArrowUp:
						GameObj.Rotate(mino)
					case termbox.KeyArrowRight:
						if mino != nil && GameObj.PermittedMoves(mino)["right"] {
							mino.position.x += 1
						}
					case termbox.KeyArrowLeft:
						if mino != nil && GameObj.PermittedMoves(mino)["left"] {
							mino.position.x -= 1
						}
					case termbox.KeyArrowDown:
						if mino != nil {
							if GameObj.PermittedMoves(mino)["down"] {
								GameObj.lastMoveTime = time.Now()
								mino.position.y += 1
							} else {
								GameObj.lastMoveTime = time.Time{}.Add(1)
							}
						}
					case termbox.KeySpace:
						if mino != nil {
							if GameObj.PermittedMoves(mino)["down"] {
								_, ghost := GameObj.MaxFall(mino)
								GameObj.Mino = ghost
								GameObj.lastMoveTime = time.Time{}.Add(1)
							}
						}

					case termbox.KeyEnter:
						if mino != nil {
							GameObj.ReserveMino()
						}
					}
				}
			}
		}
	}()

	for {
		time.Sleep(1 * time.Second)
	}
}

func DrawScreen(g *Game) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	w, h := termbox.Size()
	if w < 20 || h < 40 {
		tbprint(0, 0, termbox.ColorRed, termbox.ColorDefault, "Terminal is too small")
	} else {
		g.Update()
	}

	termbox.Flush()
}
