package main

import (
	"bufio"
	"fmt"
	"time"

	"golang.org/x/term"
)

const Dead byte = 32
const Alive byte = 56

const minNeigh = 2
const maxNeigh = 3
const bornNeigh = 3

const delay = 50 * time.Millisecond

const initSurvivePercent = 25

func main() {
	if !term.IsTerminal(0) {
		return
	}
	width, height, err := term.GetSize(0)
	if err != nil {
		return
	}
	pg := GOL{}
	pg.Init(width, height)
	pg.draw()
	pg.Play()
}

func MoveCursor(x int, y int, w *bufio.Writer) {
	fmt.Fprintf(w, "\033[%d;%dH", y, x)
}
