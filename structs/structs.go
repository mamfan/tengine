package structs

// Color represents an RGB color.
type Color struct {
	R, G, B uint8
}

// World represents the game world grid.
type World struct {
	W, H int
	Px   [][]Color
}

// Camera represents a viewport into the world.
type Camera struct {
	X, Y int
	W, H int
}

// Vec2 represents a 2D integer vector.
type Vec2 struct {
	X, Y int
}

// GameObject represents a drawable entity.
type GameObject struct {
	Position Vec2
	Pivot    Vec2
	Pixels   [][]Color
	ZIndex   int
	Visible  bool
}

// InputState holds keyboard events for the current frame.
type InputState struct {
	Events []KeyboardEvent
}

// KeyboardAction describes the kind of key event.
type KeyboardAction uint8

const (
	KeyDown KeyboardAction = iota
	KeyUp
)

// KeyMod is a bitmask of modifier keys.
type KeyMod uint16

const (
	ModShift KeyMod = 1 << iota
	ModCtrl
	ModAlt
	ModMeta
)

// KeyboardEvent represents a keyboard input event.
type KeyboardEvent struct {
	Key       int
	Action    KeyboardAction
	Modifiers KeyMod
}

func (m KeyMod) Has(flag KeyMod) bool {
	return m&flag != 0
}

// NewWorld creates a new World with given dimensions.
func NewWorld(w, h int) *World {
	px := make([][]Color, h)
	for y := 0; y < h; y++ {
		row := make([]Color, w)
		px[y] = row
	}
	return &World{W: w, H: h, Px: px}
}

// Clear sets all pixels in the world to the given color.
func (w *World) Clear(c Color) {
	for y := 0; y < w.H; y++ {
		row := w.Px[y]
		for x := 0; x < w.W; x++ {
			row[x] = c
		}
	}
}

// Set sets a pixel at (x, y) to the given color.
func (w *World) Set(x, y int, c Color) {
	if x < 0 || y < 0 || x >= w.W || y >= w.H {
		return
	}
	w.Px[y][x] = c
}

// View returns a view of the world from the camera's perspective.
func (cam *Camera) View(w *World) *World {
	if cam.W <= 0 || cam.H <= 0 {
		return NewWorld(0, 0)
	}
	view := NewWorld(cam.W, cam.H)
	for y := 0; y < cam.H; y++ {
		worldY := cam.Y + y
		if worldY < 0 || worldY >= w.H {
			continue
		}
		viewRow := view.Px[y]
		worldRow := w.Px[worldY]
		for x := 0; x < cam.W; x++ {
			worldX := cam.X + x
			if worldX < 0 || worldX >= w.W {
				continue
			}
			viewRow[x] = worldRow[worldX]
		}
	}
	return view
}

// NewGameObject creates a new GameObject with the given pixel data.
func NewGameObject(pixels [][]Color) *GameObject {
	return &GameObject{
		Pixels:  pixels,
		Visible: true,
	}
}

// Width returns the width of the GameObject.
func (g *GameObject) Width() int {
	if len(g.Pixels) == 0 {
		return 0
	}
	return len(g.Pixels[0])
}

// Height returns the height of the GameObject.
func (g *GameObject) Height() int {
	return len(g.Pixels)
}

// Draw renders the GameObject onto the world.
func (g *GameObject) Draw(w *World) {
	if !g.Visible {
		return
	}
	for y := 0; y < len(g.Pixels); y++ {
		row := g.Pixels[y]
		for x := 0; x < len(row); x++ {
			worldX := g.Position.X + x - g.Pivot.X
			worldY := g.Position.Y + y - g.Pivot.Y
			w.Set(worldX, worldY, row[x])
		}
	}
}
