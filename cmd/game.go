package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	chip8 "github.com/palSagnik/chip-8"
)

type Game struct{
	Cpu *chip8.CPU
	Pixel *ebiten.Image
	AudioPlayer *audio.Player
}

func (g *Game) Update() error {
	// check input per tick
	// we need the input state, before the FDE cycle
	// putting it after, would mean that input arrives 1 tick late
	g.keyInput()

	// Approximating Cpu speed: 600Hz
	// This implies 600 / 60 => 10 FDE cycles
	for range 10 {
		g.Cpu.Fetch()
		g.Cpu.Decode()
	}

	// Playing sound
	if g.Cpu.SoundTimer > 0 {
		if !g.AudioPlayer.IsPlaying() {
			g.AudioPlayer.Rewind()
			g.AudioPlayer.Play()
		}
	} else {
		g.AudioPlayer.Pause()
	}

	// Decrementing timers
	if g.Cpu.DelayTimer > 0 { g.Cpu.DelayTimer-- }
	if g.Cpu.SoundTimer > 0 { g.Cpu.SoundTimer-- }

	return nil
}

// Draw handles the display by going over the output buffer
// then rendering a 10X10 rectangular block over every output[x][y] = 1
func (g *Game) Draw(screen *ebiten.Image) {
	for row := range 64 {
		for col := range 32 {
			if g.Cpu.Output[row][col] == 1 {
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

func (g *Game) drawRect(screen *ebiten.Image, x, y, w, h float64, c color.Color) {
    op := &ebiten.DrawImageOptions{}
    op.GeoM.Scale(w, h)
    op.GeoM.Translate(x, y)
    op.ColorScale.ScaleWithColor(c)
    screen.DrawImage(g.Pixel, op)
}