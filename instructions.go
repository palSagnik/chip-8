package chip8

import (
	"fmt"
	"math/rand"
)

// 00E0 - CLS
// Clear the display.
func (cpu *CPU) __00E0() {
	cpu.Output = [64][32]byte{}
}

// 00EE - RET
// Return from a subroutine.
func (cpu *CPU) __00EE() {
	cpu.PC = cpu.Stack[cpu.SP]
	cpu.SP--
}

// 1nnn - JP addr
// Jump to location nnn.
func (cpu *CPU) __1nnn() {
	cpu.PC = cpu.Instruction & 0xFFF
}

// 2nnn - CALL addr
// Call subroutine at nnn.
func (cpu *CPU) __2nnn() {
	cpu.SP++
	cpu.Stack[cpu.SP] = cpu.PC
	cpu.PC = cpu.Instruction & 0xFFF
}

// 3xkk - SE Vx, byte
// Skip next instruction if Vx = kk.
func (cpu *CPU) __3xkk() {
	x := (cpu.Instruction & 0x0F00) >> 8
	kk := cpu.Instruction & 0xFF
	if cpu.V[x] == byte(kk) {
		cpu.PC += 2
	}
}

// 4xkk - SNE Vx, byte
// Skip next instruction if Vx != kk.
func (cpu *CPU) __4xkk() {
	x := (cpu.Instruction & 0x0F00) >> 8
	kk := cpu.Instruction & 0xFF
	if cpu.V[x] != byte(kk) {
		cpu.PC += 2
	}
}

// 5xy0 - SE Vx, Vy
// Skip next instruction if Vx = Vy.
func (cpu *CPU) __5xy0() {
	x := (cpu.Instruction & 0x0F00) >> 8
	y := (cpu.Instruction & 0x00F0) >> 4
	if cpu.V[x] == cpu.V[y] {
		cpu.PC += 2
	}
}

// 6xkk - LD Vx, byte
// Set Vx = kk.
func (cpu *CPU) __6xkk() {
	x := (cpu.Instruction & 0x0F00) >> 8
	kk := cpu.Instruction & 0xFF
	cpu.V[x] = byte(kk)
}

// 7xkk - ADD Vx, byte
// Set Vx = Vx + kk.
func (cpu *CPU) __7xkk() {
	x := (cpu.Instruction & 0x0F00) >> 8
	kk := cpu.Instruction & 0xFF
	cpu.V[x] += byte(kk)
}

// 8xy0 - LD Vx, Vy
// Set Vx = Vy.
func (cpu *CPU) __8xy0() {
	x := (cpu.Instruction & 0x0F00) >> 8
	y := (cpu.Instruction & 0x00F0) >> 4
	cpu.V[x] = cpu.V[y]
}

// 8xy1 - OR Vx, Vy
// Set Vx = Vx OR Vy.
func (cpu *CPU) __8xy1() {
	x := (cpu.Instruction & 0x0F00) >> 8
	y := (cpu.Instruction & 0x00F0) >> 4
	cpu.V[x] |= cpu.V[y]
}

// 8xy2 - AND Vx, Vy
// Set Vx = Vx AND Vy.
func (cpu *CPU) __8xy2() {
	x := (cpu.Instruction & 0x0F00) >> 8
	y := (cpu.Instruction & 0x00F0) >> 4
	cpu.V[x] &= cpu.V[y]
}

// 8xy3 - XOR Vx, Vy
// Set Vx = Vx XOR Vy.
func (cpu *CPU) __8xy3() {
	x := (cpu.Instruction & 0x0F00) >> 8
	y := (cpu.Instruction & 0x00F0) >> 4
	cpu.V[x] ^= cpu.V[y]
}

// 8xy4 - ADD Vx, Vy
// Set Vx = Vx + Vy, set VF = carry.
func (cpu *CPU) __8xy4() {
	x := (cpu.Instruction & 0x0F00) >> 8
	y := (cpu.Instruction & 0x00F0) >> 4
	result := uint16(cpu.V[x]) + uint16(cpu.V[y])
	if result > 255 {
		cpu.V[0xF] = 1
	} else {
		cpu.V[0xF] = 0
	}
	cpu.V[x] = byte(result)
}

// 8xy5 - SUB Vx, Vy
// Set Vx = Vx - Vy, set VF = NOT borrow.
func (cpu *CPU) __8xy5() {
	x := (cpu.Instruction & 0x0F00) >> 8
	y := (cpu.Instruction & 0x00F0) >> 4
	result := uint16(cpu.V[x]) - uint16(cpu.V[y])
	if cpu.V[x] > cpu.V[y] {
		cpu.V[0xF] = 1
	} else {
		cpu.V[0xF] = 0
	}
	cpu.V[x] = byte(result)
}

// 8xy6 - SHR Vx {, Vy}
// Set Vx = Vx SHR 1.
func (cpu *CPU) __8xy6() {
	x := (cpu.Instruction & 0x0F00) >> 8
	cpu.V[0xF] = cpu.V[x] & 0x1
	cpu.V[x] >>= 1
}

// 8xy7 - SUBN Vx, Vy
// Set Vx = Vy - Vx, set VF = NOT borrow.
func (cpu *CPU) __8xy7() {
	x := (cpu.Instruction & 0x0F00) >> 8
	y := (cpu.Instruction & 0x00F0) >> 4
	result := uint16(cpu.V[y]) - uint16(cpu.V[x])
	if cpu.V[y] > cpu.V[x] {
		cpu.V[0xF] = 1
	} else {
		cpu.V[0xF] = 0
	}
	cpu.V[x] = byte(result)
}

// 8xyE - SHL Vx {, Vy}
// Set Vx = Vx SHL 1.
func (cpu *CPU) __8xyE() {
	x := (cpu.Instruction & 0x0F00) >> 8
	cpu.V[0xF] = (cpu.V[x] >> 7) & 0x1
	cpu.V[x] <<= 1
}

// 9xy0 - SNE Vx, Vy
// Skip next instruction if Vx != Vy.
func (cpu *CPU) __9xy0() {
	x := (cpu.Instruction & 0x0F00) >> 8
	y := (cpu.Instruction & 0x00F0) >> 4
	if cpu.V[x] != cpu.V[y] {
		cpu.PC += 2
	}
}

// Annn - LD I, addr
// Set I = nnn.
func (cpu *CPU) __Annn() {
	cpu.IR = cpu.Instruction & 0xFFF
}

// Bnnn - JP V0, addr
// Jump to location nnn + V0.
func (cpu *CPU) __Bnnn() {
	cpu.PC = (cpu.Instruction & 0xFFF) + uint16(cpu.V[0])
}

// Cxkk - RND Vx, byte
// Set Vx = random byte AND kk.
func (cpu *CPU) __Cxkk() {
	x := (cpu.Instruction & 0x0F00) >> 8
	kk := cpu.Instruction & 0xFF
	cpu.V[x] = byte(rand.Intn(256)) & byte(kk)
}

// Dxyn - DRW Vx, Vy, nibble
// Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.
func (cpu *CPU) __Dxyn() {
	x := (cpu.Instruction & 0x0F00) >> 8
	y := (cpu.Instruction & 0x00F0) >> 4
	n := cpu.Instruction & 0x000F

	Vx := uint16(cpu.V[x])
	Vy := uint16(cpu.V[y])
	cpu.V[0xF] = 0

	for row := range n {
		spriteMem := cpu.Memory[cpu.IR+row]
		for col := range uint16(8) {
			bit := (spriteMem >> (7 - col)) & 1
			if bit == 1 {
				px := (Vx + col) % 64
				py := (Vy + row) % 32
				if cpu.Output[px][py] == 1 {
					cpu.V[0xF] = 1
				}
				cpu.Output[px][py] ^= bit
			}
		}
	}
}

// Ex9E - SKP Vx
// Skip next instruction if key with value Vx is pressed.
func (cpu *CPU) __Ex9E() {
	x := (cpu.Instruction & 0x0F00) >> 8
	if cpu.Input[cpu.V[x]] == 1 {
		cpu.PC += 2
	}
}

// ExA1 - SKNP Vx
// Skip next instruction if key with value Vx is not pressed.
func (cpu *CPU) __ExA1() {
	x := (cpu.Instruction & 0x0F00) >> 8
	if cpu.Input[cpu.V[x]] == 0 {
		cpu.PC += 2
	}
}

// Fx07 - LD Vx, DT
// Set Vx = delay timer value.
func (cpu *CPU) __Fx07() {
	x := (cpu.Instruction & 0x0F00) >> 8
	cpu.V[x] = cpu.DelayTimer
}

// Fx0A - LD Vx, K
// Wait for a key press, store the value of the key in Vx.
func (cpu *CPU) __Fx0A() {
	x := (cpu.Instruction & 0x0F00) >> 8
	for key := range 16 {
		if cpu.Input[key] == 1 {
			cpu.V[x] = byte(key)
			return
		}
	}
	cpu.PC -= 2
}

// Fx15 - LD DT, Vx
// Set delay timer = Vx.
func (cpu *CPU) __Fx15() {
	x := (cpu.Instruction & 0x0F00) >> 8
	cpu.DelayTimer = cpu.V[x]
}

// Fx18 - LD ST, Vx
// Set sound timer = Vx.
func (cpu *CPU) __Fx18() {
	x := (cpu.Instruction & 0x0F00) >> 8
	cpu.SoundTimer = cpu.V[x]
}

// Fx1E - ADD I, Vx
// Set I = I + Vx.
func (cpu *CPU) __Fx1E() {
	x := (cpu.Instruction & 0x0F00) >> 8
	cpu.IR += uint16(cpu.V[x])
}

// Fx29 - LD F, Vx
// Set I = location of sprite for digit Vx.
func (cpu *CPU) __Fx29() {
	x := (cpu.Instruction & 0x0F00) >> 8
	cpu.IR = uint16(cpu.V[x]) * 5
}

// Fx33 - LD B, Vx
// Store BCD representation of Vx in memory locations I, I+1, I+2.
func (cpu *CPU) __Fx33() {
	x := (cpu.Instruction & 0x0F00) >> 8
	num := cpu.V[x]
	for i := 2; i >= 0; i-- {
		cpu.Memory[cpu.IR+uint16(i)] = num % 10
		num /= 10
	}
}

// Fx55 - LD [I], Vx
// Store registers V0 through Vx in memory starting at location I.
func (cpu *CPU) __Fx55() {
	x := (cpu.Instruction & 0x0F00) >> 8
	for i := uint16(0); i <= x; i++ {
		cpu.Memory[cpu.IR+i] = cpu.V[i]
	}
}

// Fx65 - LD Vx, [I]
// Read registers V0 through Vx from memory starting at location I.
func (cpu *CPU) __Fx65() {
	x := (cpu.Instruction & 0x0F00) >> 8
	for i := uint16(0); i <= x; i++ {
		cpu.V[i] = cpu.Memory[cpu.IR+i]
	}
}

func (cpu *CPU) __invalid() {
	fmt.Printf("Invalid opcode: 0x%X\n", cpu.Instruction)
}
