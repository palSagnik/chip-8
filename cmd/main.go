package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	chip8 "github.com/palSagnik/chip-8"
)

type Game struct{
	cpu *chip8.CPU
	pixel *ebiten.Image
}

func (g *Game) Update() error {
	return nil
}

// Draw handles the display by going over the output buffer
// then rendering a 10X10 rectangular block over every output[x][y] = 1
func (g *Game) Draw(screen *ebiten.Image) {
	for row := range 64 {
		for col := range 32 {
			if g.cpu.Output[row][col] == 1 {
				g.drawRect(screen, float64(row * 10), float64(col * 10), 10, 10, color.White)
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// easy to map output to screen size
	// scaling of 10
	return 640, 320
}

func main()  {
	game := Game{
		cpu: &chip8.CPU{},
	}
	game.pixel = ebiten.NewImage(1, 1)
	game.pixel.Fill(color.White)

	// Initialising the CPU
	game.cpu.Init()

	err := ebiten.RunGame(&game)
	if err != nil {
		panic(err)
	}
}

func (g *Game) drawRect(screen *ebiten.Image, x, y, w, h float64, c color.Color) {
    op := &ebiten.DrawImageOptions{}
    op.GeoM.Scale(w, h)
    op.GeoM.Translate(x, y)
    op.ColorScale.ScaleWithColor(c)
    screen.DrawImage(g.pixel, op)
}