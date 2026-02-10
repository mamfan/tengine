package main

import (
	"fmt"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"

	"tengine/render"
	"tengine/structs"
)

type Vec2 = structs.Vec2
type GameObject = structs.GameObject
type InputState = structs.InputState
type KeyboardEvent = structs.KeyboardEvent
type KeyboardAction = structs.KeyboardAction
type KeyMod = structs.KeyMod
type Color = structs.Color
type World = structs.World
type Camera = structs.Camera

const (
	KeyDown = structs.KeyDown
	KeyUp   = structs.KeyUp

	ModShift = structs.ModShift
	ModCtrl  = structs.ModCtrl
	ModAlt   = structs.ModAlt
	ModMeta  = structs.ModMeta
)

// NewGameObject creates a new GameObject (convenience wrapper).
func NewGameObject(pixels [][]Color) *GameObject {
	return structs.NewGameObject(pixels)
}

// DrawObjects draws all objects sorted by Z-index.
func DrawObjects(w *World, objects []*GameObject) {
	sort.SliceStable(objects, func(i, j int) bool {
		return objects[i].ZIndex < objects[j].ZIndex
	})
	for _, obj := range objects {
		obj.Draw(w)
	}
}

func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func pollInput(events <-chan KeyboardEvent, camera *Camera, world *World) bool {
	for {
		select {
		case ev, ok := <-events:
			if !ok {
				return true
			}
			if ev.Modifiers.Has(ModCtrl) && ev.Key == int('C') {
				return true
			}
			if ev.Action != KeyDown {
				continue
			}
			switch ev.Key {
			case int('W'):
				camera.Y = clamp(camera.Y-1, 0, world.H-camera.H)
			case int('S'):
				camera.Y = clamp(camera.Y+1, 0, world.H-camera.H)
			case int('A'):
				camera.X = clamp(camera.X-1, 0, world.W-camera.W)
			case int('D'):
				camera.X = clamp(camera.X+1, 0, world.W-camera.W)
			}
		default:
			return false
		}
	}
}

func main() {
	world := render.NewWorld(720, 270)
	bg := Color{R: 255, G: 255, B: 255}

	fmt.Print("\x1b[2J\x1b[H\x1b[?25l")
	defer fmt.Print("\x1b[0m\x1b[?25h")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	inputEvents, restoreInput, err := StartInput()
	if err != nil {
		fmt.Println("Input init error:", err)
		return
	}
	defer restoreInput()

	red := Color{R: 255, G: 0, B: 0}
	green := Color{R: 0, G: 200, B: 0}

	player := NewGameObject([][]Color{
		{red, red, red, red},
		{red, red, red, red},
	})
	player.Position = Vec2{X: 2, Y: 2}
	player.Pivot = Vec2{X: 1, Y: 1}
	player.ZIndex = 2

	enemy := NewGameObject([][]Color{
		{green, green, green},
		{green, green, green},
		{green, green, green},
	})
	enemy.Position = Vec2{X: world.W - 3, Y: world.H - 3}
	enemy.Pivot = Vec2{X: 1, Y: 1}
	enemy.ZIndex = 0

	alex, err := NewGameObjectFromFile("../assets/aaa.timg")
	if err != nil {
		fmt.Println("Error loading game object:", err)
		return
	}
	alex.Position = Vec2{X: 1, Y: 1}
	alex.Pivot = Vec2{X: 1, Y: 1}
	alex.ZIndex = 1

	objects := []*GameObject{player, enemy, alex}

	camera := Camera{X: 0, Y: 0, W: 480, H: 270}
	if camera.W > world.W {
		camera.W = world.W
	}
	if camera.H > world.H {
		camera.H = world.H
	}

	// vx, vy := 3, 1
	ticker := time.NewTicker((1000 / 30) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			if pollInput(inputEvents, &camera, world) {
				return
			}
			world.Clear(bg)

			// player.Position.X += vx
			// player.Position.Y += vy
			// if player.Position.X <= 1 || player.Position.X >= world.W-2 {
			// 	vx = -vx
			// }
			// if player.Position.Y <= 1 || player.Position.Y >= world.H-2 {
			// 	vy = -vy
			// }

			// alex.Position.X += vx
			// if alex.Position.X <= 1 || alex.Position.X >= world.W-2 {
			// 	vx = -vx
			// }

			DrawObjects(world, objects)

			view := camera.View(world)
			render.Render(view)
		}
	}
}
