package render

import (
	"bytes"
	"fmt"
	"strconv"

	"tengine/structs"
)

type Color = structs.Color
type World = structs.World
type Camera = structs.Camera

// NewWorld is a convenience wrapper.
func NewWorld(w, h int) *World {
	return structs.NewWorld(w, h)
}

// Render outputs the world to the terminal.
func Render(w *World) {
	var buf bytes.Buffer
	buf.Grow(w.W * w.H * 25)

	buf.WriteString("\x1b[H")

	var lastColor Color
	firstPixel := true

	for y := 0; y < w.H; y++ {
		row := w.Px[y]
		for x := 0; x < w.W; x++ {
			c := row[x]
			if firstPixel || c != lastColor {
				buf.WriteString("\x1b[48;2;")
				buf.WriteString(strconv.Itoa(int(c.R)))
				buf.WriteByte(';')
				buf.WriteString(strconv.Itoa(int(c.G)))
				buf.WriteByte(';')
				buf.WriteString(strconv.Itoa(int(c.B)))
				buf.WriteByte('m')
				lastColor = c
				firstPixel = false
			}
			buf.WriteByte(' ')
		}
		buf.WriteString("\x1b[0m\n")
		firstPixel = true
	}

	fmt.Print(buf.String())
}
