package main

import "github.com/hajimehoshi/ebiten/v2"

func (g *Game) keyInput() {
	g.Cpu.Input[0x0] = boolToByte(ebiten.IsKeyPressed(ebiten.KeyX))
	g.Cpu.Input[0x1] = boolToByte(ebiten.IsKeyPressed(ebiten.Key1))
	g.Cpu.Input[0x2] = boolToByte(ebiten.IsKeyPressed(ebiten.Key2))
	g.Cpu.Input[0x3] = boolToByte(ebiten.IsKeyPressed(ebiten.Key3))
	g.Cpu.Input[0x4] = boolToByte(ebiten.IsKeyPressed(ebiten.KeyQ))
	g.Cpu.Input[0x5] = boolToByte(ebiten.IsKeyPressed(ebiten.KeyW))
	g.Cpu.Input[0x6] = boolToByte(ebiten.IsKeyPressed(ebiten.KeyE))
	g.Cpu.Input[0x7] = boolToByte(ebiten.IsKeyPressed(ebiten.KeyA))
	g.Cpu.Input[0x8] = boolToByte(ebiten.IsKeyPressed(ebiten.KeyS))
	g.Cpu.Input[0x9] = boolToByte(ebiten.IsKeyPressed(ebiten.KeyD))
	g.Cpu.Input[0xA] = boolToByte(ebiten.IsKeyPressed(ebiten.KeyZ))
	g.Cpu.Input[0xB] = boolToByte(ebiten.IsKeyPressed(ebiten.KeyC))
	g.Cpu.Input[0xC] = boolToByte(ebiten.IsKeyPressed(ebiten.Key4))
	g.Cpu.Input[0xD] = boolToByte(ebiten.IsKeyPressed(ebiten.KeyR))
	g.Cpu.Input[0xE] = boolToByte(ebiten.IsKeyPressed(ebiten.KeyF))
	g.Cpu.Input[0xF] = boolToByte(ebiten.IsKeyPressed(ebiten.KeyV))
}

func boolToByte(b bool) byte {
	if b { return 1 }
	return 0
}