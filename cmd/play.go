package main

import (
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	chip8 "github.com/palSagnik/chip-8"
	"github.com/spf13/cobra"
)

var playFile string

var playCmd = &cobra.Command{
	Use:   "play",
	Short: "Run a CHIP-8 ROM",
	RunE:  runPlay,
}

func init() {
	playCmd.Flags().StringVarP(&playFile, "file", "f", "", "path to ROM file (required)")
	playCmd.MarkFlagRequired("file")
}

func runPlay(cmd *cobra.Command, args []string) error {
	game := Game{
		Cpu: &chip8.CPU{},
	}
	game.Pixel = ebiten.NewImage(1, 1)
	game.Pixel.Fill(color.White)
	game.AudioPlayer = createAudioPlayer(44100)
	game.Cpu.Init()

	mem, err := os.ReadFile(playFile)
	if err != nil {
		log.Fatalf("could not read ROM: %v", err)
	}
	copy(game.Cpu.Memory[0x200:], mem)

	return ebiten.RunGame(&game)
}
