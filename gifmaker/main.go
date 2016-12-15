package main

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path"
	"time"

	"image/gif"

	"github.com/andyleap/giorengine"
)

var Palette = []color.Color{
	color.RGBA{0xff, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0xFF, 0xff},
	color.RGBA{0x00, 0x80, 0x00, 0xff},
	color.RGBA{0x80, 0x00, 0x80, 0xff},
	color.RGBA{0x00, 0x80, 0x80, 0xff},
	color.RGBA{0x00, 0x46, 0x00, 0xff},
	color.RGBA{0xff, 0xa5, 0x00, 0xff},
	color.RGBA{0xa5, 0x2a, 0x2a, 0xff},

	color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}, //8
	color.RGBA{0x00, 0x00, 0x00, 0xFF}, //9
	color.RGBA{0x80, 0x80, 0x80, 0xFF}, //10
	color.RGBA{0xa0, 0xa0, 0xa0, 0xFF}, //11
}

func main() {
	filename := "match.gior"
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}
	fmt.Println("Processing", filename)
	file, _ := os.Open(filename)

	g, _ := giorengine.LoadReplay(file)

	time.Sleep(5 * time.Second)
	frames := []*image.Paletted{}
	delays := []int{}
	disposal := []byte{}
	running := true
	for running {
		i := image.NewPaletted(image.Rect(0, 0, g.Match.Width*10, g.Match.Height*10), Palette)
		for x := 0; x < g.Match.Width; x++ {
			for y := 0; y < g.Match.Height; y++ {
				tileColor := Palette[8]
				cell := g.Cells[y*g.Match.Width+x]
				if cell.Faction != nil {
					tileColor = Palette[cell.Faction.ID]
				} else if cell.Pop > 0 {
					tileColor = Palette[11]
				}
				if cell.Type == giorengine.Mountain {
					tileColor = Palette[9]
				}
				for x2 := 0; x2 < 10; x2++ {
					for y2 := 0; y2 < 10; y2++ {
						if cell.Type == giorengine.City && x2 > 2 && x2 < 7 && y2 > 2 && y2 < 7 {
							i.Set(x*10+x2, y*10+y2, Palette[10])
						} else if cell.Type == giorengine.General && x2 > 2 && x2 < 7 && y2 > 2 && y2 < 7 {
							i.Set(x*10+x2, y*10+y2, Palette[8])
						} else {
							i.Set(x*10+x2, y*10+y2, tileColor)
						}
					}
				}
			}
		}
		delays = append(delays, 0)
		frames = append(frames, i)
		disposal = append(disposal, gif.DisposalNone)
		running = g.Step()
		fmt.Print("Turn ", g.Turn, "\r")
		/*f, _ := os.OpenFile(fmt.Sprintf("turn%d.gif", g.Turn), os.O_WRONLY|os.O_CREATE, 0600)
		defer f.Close()
		gif.Encode(f, i, nil)*/
	}

	var extension = path.Ext(filename)
	var name = filename[0 : len(filename)-len(extension)]
	fmt.Println("Saving", name+".gif")
	time.Sleep(10 * time.Second)
	f, _ := os.Create(name + ".gif")
	defer f.Close()
	gif.EncodeAll(f, &gif.GIF{
		Image:    frames,
		Delay:    delays,
		Disposal: disposal,
		Config: image.Config{
			Width:      g.Match.Width * 10,
			Height:     g.Match.Height * 10,
			ColorModel: color.Palette(Palette),
		},
	})
}
