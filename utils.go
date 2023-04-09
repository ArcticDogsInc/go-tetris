package main

import (
	"log"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

func gracefulStop() {
	log.Print("stop signal received")
	termbox.Close()
}

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}
