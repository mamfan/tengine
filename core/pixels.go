package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"os"
// )

// type PixelAsset struct {
// 	Width  int         `json:"width"`
// 	Height int         `json:"height"`
// 	Pixels [][][]uint8 `json:"pixels"`
// }

// func LoadPixelsFromFile(path string) ([][]Color, error) {
// 	data, err := os.ReadFile(path)
// 	if err != nil {
// 		return nil, fmt.Errorf("read pixel asset: %w", err)
// 	}

// 	var asset PixelAsset
// 	if err := json.Unmarshal(data, &asset); err != nil {
// 		return nil, fmt.Errorf("parse pixel asset: %w", err)
// 	}

// 	if asset.Width <= 0 || asset.Height <= 0 {
// 		return nil, fmt.Errorf("invalid asset dimensions")
// 	}
// 	if len(asset.Pixels) != asset.Height {
// 		return nil, fmt.Errorf("asset height mismatch")
// 	}

// 	pixels := make([][]Color, asset.Height)
// 	for y := 0; y < asset.Height; y++ {
// 		row := asset.Pixels[y]
// 		if len(row) != asset.Width {
// 			return nil, fmt.Errorf("asset width mismatch at row %d", y)
// 		}
// 		outRow := make([]Color, asset.Width)
// 		for x := 0; x < asset.Width; x++ {
// 			rgb := row[x]
// 			if len(rgb) < 3 {
// 				return nil, fmt.Errorf("invalid color at (%d,%d)", x, y)
// 			}
// 			outRow[x] = Color{R: rgb[0], G: rgb[1], B: rgb[2]}
// 		}
// 		pixels[y] = outRow
// 	}

// 	return pixels, nil
// }

// func NewGameObjectFromFile(path string) (*GameObject, error) {
// 	pixels, err := LoadPixelsFromFile(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return NewGameObject(pixels), nil
// }


import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/klauspost/compress/zstd"
)

func LoadPixelsFromFile(path string) ([][]Color, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open asset: %w", err)
	}
	defer file.Close()

	decoder, err := zstd.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("create zstd decoder: %w", err)
	}
	defer decoder.Close()

	var width uint32
	var height uint32

	if err := binary.Read(decoder, binary.LittleEndian, &width); err != nil {
		return nil, fmt.Errorf("read width: %w", err)
	}
	if err := binary.Read(decoder, binary.LittleEndian, &height); err != nil {
		return nil, fmt.Errorf("read height: %w", err)
	}

	if width == 0 || height == 0 {
		return nil, fmt.Errorf("invalid asset dimensions")
	}

	rawSize := int(width * height * 3)
	rawPixels := make([]byte, rawSize)

	_, err = io.ReadFull(decoder, rawPixels)
	if err != nil {
		return nil, fmt.Errorf("read pixel data: %w", err)
	}

	// تبدیل به [][]Color
	pixels := make([][]Color, height)

	idx := 0
	for y := 0; y < int(height); y++ {
		row := make([]Color, width)
		for x := 0; x < int(width); x++ {
			row[x] = Color{
				R: rawPixels[idx],
				G: rawPixels[idx+1],
				B: rawPixels[idx+2],
			}
			idx += 3
		}
		pixels[y] = row
	}

	return pixels, nil
}

func NewGameObjectFromFile(path string) (*GameObject, error) {
	pixels, err := LoadPixelsFromFile(path)
	if err != nil {
		return nil, err
	}
	return NewGameObject(pixels), nil
}
