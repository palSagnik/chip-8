package main

import (
	"bytes"
	"math"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

func createAudioPlayer(sampleRate int) *audio.Player {
	audioCtx := audio.NewContext(sampleRate)
	beepSamples := generateBeepSamples()
	loop := audio.NewInfiniteLoop(bytes.NewReader(beepSamples), int64(len(beepSamples)))

	player, _ := audioCtx.NewPlayer(loop)
	return player
}

func generateBeepSamples() []byte {
	sampleRate := 44100.0
	frequency := 440.0
	numSamples := int(sampleRate * 1)	// per second
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