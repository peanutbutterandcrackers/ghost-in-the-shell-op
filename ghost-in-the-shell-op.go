package main

import (
	"math/rand"
	"time"
	"flag"

	"github.com/nsf/termbox-go"
)

// for decrease animation
type Position struct {
	x int
	y int
}
var positions = []Position{}

var title_string string
var subtitle_string string

var w int  = 0 // width
var h int = 0 // height

var frame_count int = 0

const DURATION = 70

// threshold
const DRAWBKGD_TH  = 10
const DRAWTITLE_TH  = 20
const DEC_TH = 50

func draw() {
	// clear screen
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	draw_background()
	draw_subtitle(subtitle_string)
	draw_title(title_string)

	termbox.Flush()
}

func draw_subtitle(name string) {
	name_len := len(name)
	pos_y := h / 2 - 1

	if frame_count < DRAWBKGD_TH {
		return
	}

	if frame_count < DRAWTITLE_TH {
		for _, char := range name {
			rand_x := rand.Intn(w)
			termbox.SetCell(rand_x, pos_y, char, termbox.ColorGreen, termbox.ColorDefault)
		}
	} else {
		first_pos_x := w / 2 - name_len / 2
		for index, char := range name {
			termbox.SetCell(first_pos_x + index, pos_y, char, termbox.ColorGreen, termbox.ColorDefault)
		}
	}
}

func draw_title(name string) {
	name_len := len(name)
	pos_y := h / 2

	if frame_count < DRAWBKGD_TH {
		return
	}

	if frame_count < DRAWTITLE_TH {
		for _, char := range name {
			rand_x := rand.Intn(w)
			termbox.SetCell(rand_x, pos_y, char, termbox.ColorGreen, termbox.ColorDefault)
		}
	} else {
		first_pos_x := w / 2 - name_len / 2
		for index, char := range name {
			termbox.SetCell(first_pos_x + index, pos_y, char, termbox.ColorGreen, termbox.ColorDefault)
		}
	}
}

func draw_background() {
	if len(positions) != 0 && frame_count > DRAWBKGD_TH {
		shuffle(positions)
		positions = decrease_cells(positions)
	}

	for _, position := range positions {
		// The digits 0 through 9 are denoted by ASCII code points 48 through 57; 0 is 48, 1 is 49...
		// The same code is used by UTF-8 (for backwards compatibility), the default encoding used by Go.
		// The function termbox.SetCell() takes the 'rune' of a character to be placed in a given cell.
		// Rune is Go-lingo for 'Code Point'.
		// The following line adds a random integer from the range [0, 10) and to 48 (the code point for 0),
		// essentially giving us a random number in the range [48, 57], a code points for a random digit.
		// rune() is then called for type-coercion to change the random code point from type int to type rune
		cell_char := rune(48 + rand.Intn(10))

		// SetCell(x, y, Character, ForegroundColor, BackgroundColor)		
		termbox.SetCell(position.x, position.y, cell_char, termbox.ColorGreen, termbox.ColorDefault)
	}
}

func init_positions() {
	w, h = termbox.Size()

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			positions = append(positions, Position{x, y})
		}
	}
}

func decrease_cells(data []Position) (result []Position) {
	n := 0
	size := len(data) * 60 / 100 // decrease data size to 60% of it's initial size

	if size < DEC_TH {
		return result
	}

	for _, value := range data {
		if n > size {
			break
		}
		result = append(result, value)
		n++
	}

	return result
}

func shuffle(data []Position) {
	n := len(data)
	for i := n - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		data[i], data[j] = data[j], data[i]
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	flag.StringVar(&title_string, "t", "TAKUMA MATSUSHITA", "set title")
	flag.StringVar(&subtitle_string, "st", "Written by", "set subtitle")
	flag.Parse()

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	event_queue := make(chan termbox.Event)
	go func() {
		for {
			event_queue <- termbox.PollEvent()
		}
	}()

	init_positions()
	draw()

	// termbox
loop:
	for {
		select {
		case ev := <-event_queue:
			if ev.Type == termbox.EventKey&& ev.Key == termbox.KeyEsc {
				break loop
			}
		default:
			draw()
			time.Sleep(DURATION * time.Millisecond)
		}
		frame_count++
	}
}

