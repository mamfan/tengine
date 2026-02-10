package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "image"
    _ "image/gif"
    _ "image/jpeg"
    _ "image/png"
    "os"
    "path/filepath"
    "strings"
)

type PixelAsset struct {
    Width  int         `json:"width"`
    Height int         `json:"height"`
    Pixels [][][]uint8 `json:"pixels"`
}

func main() {
    input := flag.String("in", "", "input image path (png/jpg/gif)")
    output := flag.String("out", "", "output file path (.timg.json)")
    flag.Parse()

    if *input == "" {
        fmt.Fprintln(os.Stderr, "missing -in")
        os.Exit(1)
    }

    inFile, err := os.Open(*input)
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to open image: %v\n", err)
        os.Exit(1)
    }
    defer inFile.Close()

    img, _, err := image.Decode(inFile)
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to decode image: %v\n", err)
        os.Exit(1)
    }

    bounds := img.Bounds()
    width := bounds.Max.X - bounds.Min.X
    height := bounds.Max.Y - bounds.Min.Y

    asset := PixelAsset{
        Width:  width,
        Height: height,
        Pixels: make([][][]uint8, height),
    }

    for y := 0; y < height; y++ {
        row := make([][]uint8, width)
        for x := 0; x < width; x++ {
            r, g, b, _ := img.At(x+bounds.Min.X, y+bounds.Min.Y).RGBA()
            row[x] = []uint8{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)}
        }
        asset.Pixels[y] = row
    }

    outPath := *output
    if outPath == "" {
        base := strings.TrimSuffix(*input, filepath.Ext(*input))
        outPath = base + ".timg.json"
    }

    data, err := json.Marshal(asset)
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to encode json: %v\n", err)
        os.Exit(1)
    }

    if err := os.WriteFile(outPath, data, 0644); err != nil {
        fmt.Fprintf(os.Stderr, "failed to write file: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("saved %s (%dx%d)\n", outPath, width, height)
}
