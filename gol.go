package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

type GOL struct {
	w        int
	h        int
	dataLen  int
	dataLast []byte
	dataNext []byte
	writer   *bufio.Writer
}

func (p *GOL) Play() {
	for {
		time.Sleep(delay)
		p.calcNextStep()
		p.draw()
	}
}

func (p *GOL) Init(w int, h int) {
	p.writer = bufio.NewWriter(os.Stdout)
	p.w = w
	p.h = h * 2
	p.dataLen = p.w * p.h
	p.dataLast = make([]byte, p.dataLen)
	p.dataNext = make([]byte, p.dataLen)
	p.randomData()
}

func (p *GOL) randomData() {
	for i := 0; i < p.dataLen; i++ {
		if rand.Intn(100) >= initSurvivePercent {
			p.dataNext[i] = Dead
		} else {
			p.dataNext[i] = Alive
		}
	}
}

func (p *GOL) calcNextStep() {
	wg := &sync.WaitGroup{}

	for i := 0; i < p.h; i++ {
		wg.Add(1)
		go func(internalI int) {
			defer wg.Done()
			s := internalI * p.w
			p.calcGroup(s, s+p.w)
		}(i)
	}

	wg.Wait()

}

func (p *GOL) calcGroup(sptr int, eptr int) {
	for i := sptr; i < eptr; i++ {
		if p.isAlive(i) == 1 {
			p.dieIfNeeds(i)
		} else {
			p.bornIfPossible(i)
		}
	}
}

func (p *GOL) draw() {
	MoveCursor(1, 1, p.writer)
	copy(p.dataLast, p.dataNext)
	fmt.Fprint(p.writer, p.compressField())
	// fmt.Fprint(p.writer, strings.ReplaceAll(string(p.dataLast), "8", "█"))
	p.writer.Flush()

}

func (p *GOL) compressField() string {
	var out string = ""
	for i := 0; i < p.h; i += 2 {
		for j := 0; j < p.w; j++ {
			upper := p.dataLast[(i*p.w)+j]
			down := p.dataLast[(i*p.w)+j+p.w]
			if upper == Alive && down == Alive {
				out += "█"
			} else if upper == Alive {
				out += "▀"
			} else if down == Alive {
				out += "▄"
			} else {
				out += " "
			}

		}
	}
	// for i := 0; i < len(p.dataLast); i++ {
	// 	out += string(p.dataLast[i])
	// }
	return out
}

func compressString() {

}

func (p *GOL) dieIfNeeds(ptr int) {
	neighbours := p.countNeighbours(ptr)
	if neighbours > maxNeigh || neighbours < minNeigh {
		p.dataNext[ptr] = Dead
	}
}

func (p *GOL) bornIfPossible(ptr int) {
	neighbours := p.countNeighbours(ptr)
	if neighbours == bornNeigh {
		p.dataNext[ptr] = Alive
	}
}

func (p *GOL) countNeighbours(ptr int) int {
	n := p.isAlive(ptr - 1)
	n += p.isAlive(ptr + 1)
	n += p.isAlive(ptr - p.w)
	n += p.isAlive(ptr - p.w - 1)
	n += p.isAlive(ptr - p.w + 1)
	n += p.isAlive(ptr + p.w)
	n += p.isAlive(ptr + p.w - 1)
	n += p.isAlive(ptr + p.w + 1)
	return n
}

func (p *GOL) isAlive(ptr int) int {
	if ptr < 0 || ptr >= p.dataLen {
		return 0
	}
	if p.dataLast[ptr] == Dead {
		return 0
	} else {
		return 1
	}
}
