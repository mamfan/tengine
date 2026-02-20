package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func LoadPixelsFromFile(path string) ([][]Color, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width == 0 || height == 0 {
		return nil, fmt.Errorf("invalid image dimensions")
	}

	// Convert to [][]Color
	pixels := make([][]Color, height)

	for y := 0; y < height; y++ {
		row := make([]Color, width)
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(bounds.Min.X+x, bounds.Min.Y+y).RGBA()
			r8, g8, b8 := unpremultiplyRGBA(r, g, b, a)
			row[x] = Color{
				R: r8,
				G: g8,
				B: b8,
				A: uint8(a >> 8),
			}
		}
		pixels[y] = row
	}

	return pixels, nil
}

func NewGameObjectFromFile(path string) (*GameObject, error) {
	pixels, err := LoadPixelsFromFile("../assets/" + path)
	if err != nil {
		return nil, err
	}
	return NewGameObject(pixels), nil
}

func unpremultiplyRGBA(r, g, b, a uint32) (uint8, uint8, uint8) {
	if a == 0 {
		return 0, 0, 0
	}
	const max = 0xFFFF
	r8 := uint8((r * max / a) >> 8)
	g8 := uint8((g * max / a) >> 8)
	b8 := uint8((b * max / a) >> 8)
	return r8, g8, b8
}
