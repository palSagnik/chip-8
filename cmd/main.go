package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	chip8 "github.com/palSagnik/chip-8"
)

type Game struct{
	cpu *chip8.CPU
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// easy to map output to screen size
	// scaling of 10
	return 640, 320
}

func main()  {
	game := Game{}

	err := ebiten.RunGame(&game)
	if err != nil {
		panic(err)
	}
}
