package main

import (
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	chip8 "github.com/palSagnik/chip-8"
)

func main()  {
	game := Game{
		Cpu: &chip8.CPU{},
	}
	game.Pixel = ebiten.NewImage(1, 1)
	game.Pixel.Fill(color.White)

	// audio
	game.AudioPlayer = createAudioPlayer(44100)

	// Initialising the CPU
	game.Cpu.Init()

	if len(os.Args) < 2 {
		log.Fatal("usage: <rom.ch8>")
	}

	path := os.Args[1]

	// Copy rom into memory
	mem, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Something went wrong: %v", err)
	}
	copy(game.Cpu.Memory[0x200:], mem)

	err = ebiten.RunGame(&game)
	if err != nil {
		panic(err)
	}
}