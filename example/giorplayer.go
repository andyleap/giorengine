package main

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"time"

	"github.com/andyleap/tinyfb"

	"github.com/andyleap/giorengine"
)

var PlayerColors = []color.RGBA{
	color.RGBA{0xff, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0xFF, 0xff},
	color.RGBA{0x00, 0x80, 0x00, 0xff},
	color.RGBA{0x80, 0x00, 0x80, 0xff},
	color.RGBA{0x00, 0x80, 0x80, 0xff},
	color.RGBA{0x00, 0x46, 0x00, 0xff},
	color.RGBA{0xff, 0xa5, 0x00, 0xff},
	color.RGBA{0xa5, 0x2a, 0x2a, 0xff},
}

func main() {
	file, _ := os.Open("SJi00J1Nx.gior")

	g, _ := giorengine.LoadReplay(file)

	t := tinyfb.New("test", int32(g.Match.Width*10), int32(g.Match.Height*10))

	quit := false
	go func() {
		t.Run()
		quit = true
	}()

	frame := time.Now()
	i := image.NewRGBA(image.Rect(0, 0, g.Match.Width*10, g.Match.Height*10))
	fmt.Println(int32(g.Match.Width*10), int32(g.Match.Height*10))
	time.Sleep(5 * time.Second)
	for !quit {
		for x := 0; x < g.Match.Width; x++ {
			for y := 0; y < g.Match.Height; y++ {
				tileColor := color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
				cell := g.Cells[y*g.Match.Width+x]
				if cell.Faction != nil {
					tileColor = PlayerColors[cell.Faction.ID]
				}
				if cell.Type == giorengine.Mountain {
					tileColor = color.RGBA{0x00, 0x00, 0x00, 0xFF}
				}
				for x2 := 0; x2 < 10; x2++ {
					for y2 := 0; y2 < 10; y2++ {
						if cell.Type == giorengine.City && x2 > 2 && x2 < 7 && y2 > 2 && y2 < 7 {
							i.SetRGBA(x*10+x2, y*10+y2, color.RGBA{0x80, 0x80, 0x80, 0xFF})
						} else if cell.Type == giorengine.General && x2 > 2 && x2 < 7 && y2 > 2 && y2 < 7 {
							i.SetRGBA(x*10+x2, y*10+y2, color.RGBA{0xFF, 0xFF, 0xFF, 0xFF})
						} else {
							i.SetRGBA(x*10+x2, y*10+y2, tileColor)
						}
					}
				}

			}
		}
		g.Step()
		t.Update(i)
		end := time.Now()
		delta := end.Sub(frame).Nanoseconds() - (time.Second / 12).Nanoseconds()
		if delta < 0 {
			time.Sleep(time.Duration(-delta))
		}
		frame = time.Now()
	}

}
