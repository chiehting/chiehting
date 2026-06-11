package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"math"
	"os"

	"github.com/golang/freetype"
)

type Circle struct {
	X, Y, R float64
}

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
		color.RGBA{0x00, 0x00, 0x00, 0xff}, // 黑
		color.RGBA{0x00, 0x00, 0xff, 0xff}, // 藍
		color.RGBA{0x00, 0xff, 0x00, 0xff}, // 綠
		color.RGBA{0x00, 0xff, 0xff, 0xff}, // 青
		color.RGBA{0xff, 0x00, 0x00, 0xff}, // 紅
		color.RGBA{0xff, 0x00, 0xff, 0xff}, // 品紅
		color.RGBA{0xff, 0xff, 0x00, 0xff}, // 黃
		color.RGBA{0xff, 0xff, 0xff, 0xff}, // 白
	}

	fontPath := "/usr/share/fonts/truetype/ubuntu/UbuntuSansMono-Italic[wght].ttf"
	fontBytes, err := os.ReadFile(fontPath)
	if err != nil {
		fmt.Printf(fontPath, err)
		return
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	var images []*image.Paletted
	var delays []int
	steps := 240

	for step := 0; step <= steps; step++ {
		rgbaImg := image.NewRGBA(image.Rect(0, 0, w, h))
		delays = append(delays, 4)

		currentX := float64(w) / float64(steps) * float64(step)

		for i := 0; i < 3; i++ {
			circles[i] = &Circle{
				X: currentX,
				Y: math.Sin(currentX*float64(i+1)*math.Pi/180)*float64((h-4)/2) + float64(h/2),
				R: float64(5 * (i + 1)),
			}
		}

		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				r := circles[0].Brightness(float64(x), float64(y))
				g := circles[1].Brightness(float64(x), float64(y))
				b := circles[2].Brightness(float64(x), float64(y))

				if r > 0 || g > 0 || b > 0 {
					outR, outG, outB := r, g, b
					rgbaImg.Set(x, y, color.RGBA{outR, outG, outB, 255})

				} else if y == h/2 {
					rgbaImg.Set(x, y, color.RGBA{255, 255, 255, 255})
				} else {
					rgbaImg.Set(x, y, color.RGBA{0, 0, 0, 255})
				}
			}

			for i := 1; i <= 3; i++ {
				y := math.Sin(float64(x*i)*math.Pi/180) * float64((h-4)/2)
				rgbaImg.Set(x, int(y)+(h/2), color.RGBA{
					uint8((i + 1) % 3 * 255),
					uint8((i + 2) % 3 * 255),
					uint8((i + 3) % 3 * 255),
					255,
				})
			}
		}

		c := freetype.NewContext()
		c.SetDPI(96)
		c.SetFont(font)
		c.SetFontSize(50)
		c.SetClip(rgbaImg.Bounds())
		c.SetDst(rgbaImg)
		c.SetSrc(image.NewUniform(color.RGBA{255, 255, 255, 255}))

		_, _ = c.DrawString("H e l l o", freetype.Pt(50, 65))
		_, _ = c.DrawString("W o r l d", freetype.Pt(460, 65))

		palettedImg := image.NewPaletted(image.Rect(0, 0, w, h), palette)
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				palettedImg.Set(x, y, rgbaImg.At(x, y))
			}
		}
		images = append(images, palettedImg)
	}

	f, err := os.OpenFile("sine.gif", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	err = gif.EncodeAll(f, &gif.GIF{
		Image: images,
		Delay: delays,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

}
