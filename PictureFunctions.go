package main

import (
	"image"
	"image/png"
	"log"
	"os"

	"github.com/fogleman/gg"
)

// ClearPic is used to create an empty white bacground
func ClearPic(dc *gg.Context) {
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
}

// SaveGrayPic converts image to 8 bit gray scale png and saves it to a given file
func SaveGrayPic(pic image.Image, imagePath string) {
	gray := convert(pic)
	outfile, err := os.Create(imagePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer outfile.Close()
	png.Encode(outfile, gray)
}

func convert(m image.Image) *image.Gray {
	b := m.Bounds()
	gray := image.NewGray(b)
	pos := 0
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, _, _, _ := m.At(x, y).RGBA() // no gray scale juts blac & white algo, indicator is Red
			gray.Pix[pos] = uint8(r >> 8)
			pos++
		}
	}
	return gray
}
