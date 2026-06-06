package main

import (
	"bytes"
	"image/color"
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	chip8 "github.com/palSagnik/chip-8"
)

type Game struct{
	cpu *chip8.CPU
	pixel *ebiten.Image
	audioPlayer *audio.Player
}

func (g *Game) Update() error {
	// check input per tick
	// we need the input state, before the FDE cycle
	// putting it after, would mean that input arrives 1 tick late
	g.keyInput()

	// Approximating CPU speed: 600Hz
	// This implies 600 / 60 => 10 FDE cycles
	for range 10 {
		g.cpu.Fetch()
		g.cpu.Decode()
	}

	// Playing sound
	if g.cpu.SoundTimer > 0 {
		if !g.audioPlayer.IsPlaying() {
			g.audioPlayer.Rewind()
			g.audioPlayer.Play()
		}
	} else {
		g.audioPlayer.Pause()
	}

	// Decrementing timers
	if g.cpu.DelayTimer > 0 { g.cpu.DelayTimer-- }
	if g.cpu.SoundTimer > 0 { g.cpu.SoundTimer-- }

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

	// audio
	audioCtx := audio.NewContext(44100)
	beepSamples := generateBeepSamples()
	loop := audio.NewInfiniteLoop(bytes.NewReader(beepSamples), int64(len(beepSamples)))
	player, _ := audioCtx.NewPlayer(loop)
	game.audioPlayer = player

	// Initialising the CPU
	game.cpu.Init()

	// Copy rom into memory
	mem, err := os.ReadFile("test/7-beep.ch8")
	if err != nil {
		log.Fatalf("Something went wrong: %v", err)
	}
	copy(game.cpu.Memory[0x200:], mem)

	err = ebiten.RunGame(&game)
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

func (g *Game) keyInput() {
	g.cpu.Input[0x0] = boolToByte(ebiten.IsKeyPressed(ebiten.KeyX))
	g.cpu.Input[0x1] = boolToByte(ebiten.IsKeyPressed(ebiten.Key1))
	g.cpu.Input[0x2] = boolToByte(ebiten.IsKeyPressed(ebiten.Key2))
	g.cpu.Input[0x3] = boolToByte(ebiten.IsKeyPressed(ebiten.Key3))
	g.cpu.Input[0x4] = boolToByte(ebiten.IsKeyPressed(ebiten.KeyQ))
	g.cpu.Input[0x5] = boolToByte(ebiten.IsKeyPressed(ebiten.KeyW))
	g.cpu.Input[0x6] = boolToByte(ebiten.IsKeyPressed(ebiten.KeyE))
	g.cpu.Input[0x7] = boolToByte(ebiten.IsKeyPressed(ebiten.KeyA))
	g.cpu.Input[0x8] = boolToByte(ebiten.IsKeyPressed(ebiten.KeyS))
	g.cpu.Input[0x9] = boolToByte(ebiten.IsKeyPressed(ebiten.KeyD))
	g.cpu.Input[0xA] = boolToByte(ebiten.IsKeyPressed(ebiten.KeyZ))
	g.cpu.Input[0xB] = boolToByte(ebiten.IsKeyPressed(ebiten.KeyC))
	g.cpu.Input[0xC] = boolToByte(ebiten.IsKeyPressed(ebiten.Key4))
	g.cpu.Input[0xD] = boolToByte(ebiten.IsKeyPressed(ebiten.KeyR))
	g.cpu.Input[0xE] = boolToByte(ebiten.IsKeyPressed(ebiten.KeyF))
	g.cpu.Input[0xF] = boolToByte(ebiten.IsKeyPressed(ebiten.KeyV))
}

func boolToByte(b bool) byte {
	if b { return 1 }
	return 0
}

func generateBeepSamples() []byte {
	sampleRate := 44100.0
	frequency := 440.0
	numSamples := int(sampleRate * 1)
	buf := make([]byte, numSamples*4)   // 4 bytes per sample

	for i := 0; i < numSamples; i++ {
		t := float64(i) / sampleRate
		sample := math.Sin(2 * math.Pi * frequency * t)
		s := int16(sample * 32767)

		buf[i*4]   = byte(s)        // left low
		buf[i*4+1] = byte(s >> 8)   // left high
		buf[i*4+2] = byte(s)        // right low
		buf[i*4+3] = byte(s >> 8)   // right high
	}

    return buf
}