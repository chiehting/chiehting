package main

import (
	"image"
	"image/color"
	"image/gif"
	"io/ioutil"
	"math"
	"os"

	"github.com/golang/freetype"
)

// Circle position
type Circle struct {
	X, Y, R float64
}

// Brightness is return
func (c *Circle) Brightness(x, y float64) uint8 {
	var dx, dy float64 = c.X - x, c.Y - y
	d := math.Sqrt(dx*dx+dy*dy) / c.R
	if d > 1 {
		return 0
	}
	return 200
}

func main() {
	var w, h int = 840, 90
	circles := make(map[int]*Circle)
	var palette = []color.Color{
		color.RGBA{0x00, 0x00, 0x00, 0xff},
		color.RGBA{0x00, 0x00, 0xff, 0xff},
		color.RGBA{0x00, 0xff, 0x00, 0xff},
		color.RGBA{0x00, 0xff, 0xff, 0xff},
		color.RGBA{0xff, 0x00, 0x00, 0xff},
		color.RGBA{0xff, 0x00, 0xff, 0xff},
		color.RGBA{0xff, 0xff, 0x00, 0xff},
		color.RGBA{0xff, 0xff, 0xff, 0xff},
	}

	var images []*image.Paletted
	var delays []int
	steps := 70

	for step := 0; step <= steps; step++ {
		img := image.NewPaletted(image.Rect(0, 0, w, h), palette)
		images = append(images, img)
		delays = append(delays, 0)

		for i := 0; i < 3; i++ {
			circles[i] = &Circle{
				X: float64(w / steps * step),
				Y: math.Sin(float64(w/steps*step*(i+1))*math.Pi/180)*float64((h-4)/2) + float64(h/2),
				R: float64(5 * (i + 1)),
			}
		}

		for x := 0; x <= w; x++ {
			for y := 0; y < h; y++ {
				img.Set(x, y, color.RGBA{
					circles[0].Brightness(float64(x), float64(y)),
					circles[1].Brightness(float64(x), float64(y)),
					circles[2].Brightness(float64(x), float64(y)),
					255,
				})
			}

			for i := 1; i <= 3; i++ {
				y := math.Sin(float64(x*i)*math.Pi/180) * float64((h-4)/2)
				img.Set(x, int(y)+(h/2), color.RGBA{
					uint8((i + 1) % 3 * 255),
					uint8((i + 2) % 3 * 255),
					uint8((i + 3) % 3 * 255),
					255,
				})
			}
		}

		fontBytes, _ := ioutil.ReadFile("/System/Library/Fonts/Supplemental/Chalkduster.ttf")
		font, _ := freetype.ParseFont(fontBytes)
		c := freetype.NewContext()
		c.SetDPI(float64(96))
		c.SetFont(font)
		c.SetFontSize(34)
		c.SetClip(img.Bounds())
		c.SetDst(img)
		c.SetSrc(image.NewUniform(color.RGBA{255, 255, 255, 255}))
		//draw the label on image
		c.DrawString("Hello", freetype.Pt(270, 65))
		c.DrawString("World", freetype.Pt(430, 65))
	}

	f, _ := os.OpenFile("sine.gif", os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	gif.EncodeAll(f, &gif.GIF{
		Image: images,
		Delay: delays,
	})
}
