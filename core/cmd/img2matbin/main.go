package main

import (
	"encoding/binary"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"flag"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/klauspost/compress/zstd"
)

func encodeImage(inputPath, outputPath string) error {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	img, _, err := image.Decode(inFile)
	if err != nil {
		return err
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	rawPixels := make([]byte, width*height*3)

	idx := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			rawPixels[idx] = uint8(r >> 8)
			rawPixels[idx+1] = uint8(g >> 8)
			rawPixels[idx+2] = uint8(b >> 8)
			idx += 3
		}
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	encoder, err := zstd.NewWriter(outFile, zstd.WithEncoderLevel(zstd.SpeedBestCompression))
	if err != nil {
		return err
	}
	defer encoder.Close()

	// نوشتن header
	if err := binary.Write(encoder, binary.LittleEndian, uint32(width)); err != nil {
		return err
	}
	if err := binary.Write(encoder, binary.LittleEndian, uint32(height)); err != nil {
		return err
	}

	// نوشتن پیکسل‌ها
	_, err = encoder.Write(rawPixels)
	return err
}

func main() {
    input := flag.String("in", "", "input image path (png/jpg/gif)")
    output := flag.String("out", "", "output file path (.timg)")
    flag.Parse()

    if *input == "" {
        fmt.Fprintln(os.Stderr, "missing -in")
        os.Exit(1)
    }

    outPath := *output
    if outPath == "" {
        base := strings.TrimSuffix(*input, filepath.Ext(*input))
        outPath = base + ".timg"
    }

    if err := encodeImage(*input, outPath); err != nil {
        fmt.Fprintf(os.Stderr, "failed to encode image: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("saved %s\n", outPath)
}
