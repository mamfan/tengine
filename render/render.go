package render

import (
	"bytes"
	"os"
	"strconv"
	"time"

	"tengine/structs"
)

type Color = structs.Color
type World = structs.World

// NewWorld is a convenience wrapper.
func NewWorld(w, h int) *World {
	return structs.NewWorld(w, h)
}

var prev [][]Color
var initialized bool
var firstFrame bool = true

// allocate previous frame buffer
func initPrev(w *World) {
	prev = make([][]Color, w.H)
	for y := range prev {
		prev[y] = make([]Color, w.W)
	}
	initialized = true
}

// Render outputs only changed pixels to terminal.
func Render(w *World) {

	if !initialized || len(prev) != w.H || len(prev[0]) != w.W {
		initPrev(w)
		firstFrame = true // Force full redraw on size change
	}

	var buf bytes.Buffer

	// First frame: do full redraw
	if firstFrame {
		buf.Grow(w.W * w.H * 25)
		buf.WriteString("\x1b[H") // Home cursor

		var lastColor Color
		isFirst := true

		for y := 0; y < w.H; y++ {
			row := w.Px[y]
			for x := 0; x < w.W; x++ {
				c := row[x]
				if isFirst || c != lastColor {
					buf.WriteString("\x1b[48;2;")
					buf.WriteString(strconv.Itoa(int(c.R)))
					buf.WriteByte(';')
					buf.WriteString(strconv.Itoa(int(c.G)))
					buf.WriteByte(';')
					buf.WriteString(strconv.Itoa(int(c.B)))
					buf.WriteByte('m')
					lastColor = c
					isFirst = false
				}
				buf.WriteByte(' ')
				prev[y][x] = c // Update prev buffer
			}
			buf.WriteString("\x1b[0m\n")
			isFirst = true
		}

		firstFrame = false
		os.Stdout.Write(buf.Bytes())
		return
	}

	// Differential rendering for subsequent frames
	buf.Grow(w.W * w.H * 4)

	var lastColor Color
	colorSet := false

	for y := 0; y < w.H; y++ {
		row := w.Px[y]

		x := 0
		for x < w.W {

			// skip unchanged pixels
			if prev[y][x] == row[x] {
				x++
				continue
			}

			startX := x
			c := row[x]

			// find run of same-color changed pixels
			for x < w.W &&
				prev[y][x] != row[x] &&
				row[x] == c {
				x++
			}

			runLen := x - startX

			// move cursor once
			buf.WriteString("\x1b[")
			buf.WriteString(strconv.Itoa(y + 1))
			buf.WriteByte(';')
			buf.WriteString(strconv.Itoa(startX + 1))
			buf.WriteByte('H')

			// change color only if needed
			if !colorSet || c != lastColor {
				buf.WriteString("\x1b[48;2;")
				buf.WriteString(strconv.Itoa(int(c.R)))
				buf.WriteByte(';')
				buf.WriteString(strconv.Itoa(int(c.G)))
				buf.WriteByte(';')
				buf.WriteString(strconv.Itoa(int(c.B)))
				buf.WriteByte('m')

				lastColor = c
				colorSet = true
			}

			// draw run in one write
			for i := 0; i < runLen; i++ {
				buf.WriteByte(' ')
				prev[y][startX+i] = c
			}
		}
	}

	if colorSet {
		buf.WriteString("\x1b[0m")
	}

	os.Stdout.Write(buf.Bytes())
	time.Sleep(time.Millisecond) // give terminal time to process
}
