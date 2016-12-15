package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"path"

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
	color.RGBA{0, 0, 0, 0},             //12 transparent

	color.RGBA{0x00, 0x00, 0x00, 0xFF}, //13
	color.RGBA{0x11, 0x11, 0x11, 0xFF},
	color.RGBA{0x22, 0x22, 0x22, 0xFF},
	color.RGBA{0x33, 0x33, 0x33, 0xFF},
	color.RGBA{0x44, 0x44, 0x44, 0xFF},
	color.RGBA{0x55, 0x55, 0x55, 0xFF},
	color.RGBA{0x66, 0x66, 0x66, 0xFF},
	color.RGBA{0x77, 0x77, 0x77, 0xFF},
	color.RGBA{0x88, 0x88, 0x88, 0xFF},
	color.RGBA{0x99, 0x99, 0x99, 0xFF},
	color.RGBA{0xaa, 0xaa, 0xaa, 0xFF},
	color.RGBA{0xbb, 0xbb, 0xbb, 0xFF},
	color.RGBA{0xcc, 0xcc, 0xcc, 0xFF},
	color.RGBA{0xdd, 0xdd, 0xdd, 0xFF},
	color.RGBA{0xee, 0xee, 0xee, 0xFF},
	color.RGBA{0xff, 0xff, 0xff, 0xFF},
}

func main() {
	filename := "match.gior"
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}
	fmt.Println("Processing", filename)
	file, _ := os.Open(filename)

	g, _ := giorengine.LoadReplay(file)

	frames := []*image.Paletted{}
	delays := []int{}
	disposal := []byte{}
	running := true
	lastframe := image.NewPaletted(image.Rect(0, 0, g.Match.Width*10, g.Match.Height*10), Palette)
	for running {
		i := image.NewPaletted(image.Rect(0, 0, g.Match.Width*10, g.Match.Height*10), Palette)
		optimized := image.NewPaletted(image.Rect(0, 0, g.Match.Width*10, g.Match.Height*10), Palette)
		maxArmy := 0
		for x := 0; x < g.Match.Width; x++ {
			for y := 0; y < g.Match.Height; y++ {
				if g.Cells[y*g.Match.Width+x].Pop > maxArmy {
					maxArmy = g.Cells[y*g.Match.Width+x].Pop
				}
			}
		}
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

				popVal := 0
				if cell.Pop > 0 {
					logVal := math.Log(float64(cell.Pop)) / math.Log(float64(maxArmy))
					linVal := float64(cell.Pop) / float64(maxArmy)
					popVal = int(math.Floor((logVal + linVal) * 8))
				}
				//size := popVal / 4
				popColor := Palette[13+(15-popVal)]
				for x2 := 0; x2 < 10; x2++ {
					for y2 := 0; y2 < 10; y2++ {
						color := tileColor
						if x2 > 2 && x2 < 7 && y2 > 2 && y2 < 7 && cell.Type != giorengine.Mountain {
							color = popColor
						}
						if cell.Type == giorengine.City && ((x2 == 0 || x2 == 9) || (y2 == 0 || y2 == 9)) {
							color = Palette[10]
						} else if cell.Type == giorengine.General && ((x2 == 0 || x2 == 9) || (y2 == 0 || y2 == 9)) {
							color = Palette[8]
						}

						i.Set(x*10+x2, y*10+y2, color)
						if lastframe.At(x*10+x2, y*10+y2) == color {
							color = Palette[12]
						}
						optimized.Set(x*10+x2, y*10+y2, color)
					}
				}
			}
		}
		lastframe = i
		delays = append(delays, 0)
		frames = append(frames, optimized)
		disposal = append(disposal, gif.DisposalNone)
		running = g.Step()
		//fmt.Print("Turn ", g.Turn, "\r")
		/*f, _ := os.OpenFile(fmt.Sprintf("turn%d.gif", g.Turn), os.O_WRONLY|os.O_CREATE, 0600)
		defer f.Close()
		gif.Encode(f, i, nil)*/
	}

	var extension = path.Ext(filename)
	var name = filename[0 : len(filename)-len(extension)]
	fmt.Println("Saving", name+".gif")
	f, _ := os.Create(name + ".gif")
	defer f.Close()
	/*delays = append(delays, 100)
	frames = append(frames, frames[len(frames)-1])
	disposal = append(disposal, gif.DisposalNone)*/
	delays[len(delays)-1] = 100
	gif.EncodeAll(f, &gif.GIF{
		Image:    frames,
		Delay:    delays,
		Disposal: disposal,
		Config: image.Config{
			Width:      g.Match.Width * 10,
			Height:     g.Match.Height * 10,
			ColorModel: color.Palette(Palette),
		},
		BackgroundIndex: 8,
	})
}
