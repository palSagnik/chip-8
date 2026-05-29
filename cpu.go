package chip8

import (
	"fmt"
	"math/rand"
)

type CPU struct {
	V           		[16]byte            // 16 8-bit registers (V0 - VF)
	PC          		uint16              // Program Counter: stores the address of the current instruction
	SP          		byte                // Stack Pointer: refers to the top of the stack
	IR          		uint16              // I register: used to store memory address, last 12 bits used in general
	Stack       		[16]uint16          // Stack
	DelayTimer  		byte                // Delay Register
	SoundTimer       	byte                // Sound Register

	Memory 				[4096]byte          // Memory
	Output     		    [64][32]byte        // Output buffer: working as display
	Input 				[16]byte			// Input: Keypad mapping

	Opcode 				uint16				// current opcode
}

func (cpu *CPU) Fetch() {
	// update opcode with current instruction
	cpu.Opcode = uint16(cpu.Memory[cpu.PC])
	cpu.Opcode = cpu.Opcode << 8
	cpu.Opcode |= uint16(cpu.Memory[cpu.PC + 1])

	// after succesful fetch point to the next instruction
	cpu.PC += 2
}

func (cpu *CPU) Decode() {
	switch (cpu.Opcode & 0xF000) {

	case 0x0000:
		switch cpu.Opcode & 0xFF {
		case 0x00E0:
			// 00E0 - CLS
			// Clear the display
			cpu.Output = [64][32]byte{}

		// TODO: properly think about stack instructiions
		case 0x00EE:
			// 00EE - RET
            // Return from a subroutine.
			cpu.PC = cpu.Stack[cpu.SP]

			// TODO: clear the top of the stack?
			// then decrement the stack pointer
			cpu.SP -= 1
		}
	case 0x1000:
		// 1nnn - JP addr
		// Jump to location nnn.
		cpu.PC = cpu.Opcode & 0xFFF

	case 0x3000:
		// 3xkk - SE Vx, byte
		// Skip next instruction if Vx = kk.
		x := (cpu.Opcode & 0x0F00) >> 8
		kk := cpu.Opcode & 0xFF
		if cpu.V[x] == byte(kk) {
			cpu.PC += 2
		}

	case 0x4000:
		// 4xkk - SNE Vx, byte
		// Skip next instruction if Vx != kk.
		x := (cpu.Opcode & 0x0F00) >> 8
		kk := cpu.Opcode & 0xFF
		if cpu.V[x] != byte(kk) {
			cpu.PC += 2
		}

	case 0x5000:
		// 5xy0 - SE Vx, Vy
		// Skip next instruction if Vx = Vy.
		x := (cpu.Opcode & 0x0F00) >> 8
		y := (cpu.Opcode & 0x00F0) >> 4
		if cpu.V[x] == cpu.V[y] {
			cpu.PC += 2
		}

	case 0x6000:
		// 6xkk - LD Vx, byte
		// Set Vx = kk.
		x := (cpu.Opcode & 0x0F00) >> 8
		kk := cpu.Opcode & 0xFF
		cpu.V[x] = byte(kk)

	case 0x7000:
		// 7xkk - ADD Vx, byte
		// Set Vx = Vx + kk.
		x := (cpu.Opcode & 0x0F00) >> 8
		kk := cpu.Opcode & 0xFF
		cpu.V[x] += byte(kk)

	case 0x8000:
		switch (cpu.Opcode & 0x000F) {
		case 0x0:
			// 8xy0 - LD Vx, Vy
			// Set Vx = Vy.
			x := (cpu.Opcode & 0x0F00) >> 8
			y := (cpu.Opcode & 0x00F0) >> 4
			cpu.V[x] = cpu.V[y]

		case 0x1:
			// 8xy1 - OR Vx, Vy
			// Vx = Vx or Vy
			x := (cpu.Opcode & 0x0F00) >> 8
			y := (cpu.Opcode & 0x00F0) >> 4
			cpu.V[x] |= cpu.V[y]

		case 0x2:
			// 8xy2 - AND Vx, Vy
			// Vx = Vx and Vy
			x := (cpu.Opcode & 0x0F00) >> 8
			y := (cpu.Opcode & 0x00F0) >> 4
			cpu.V[x] &= cpu.V[y]

		case 0x3:
			// 8xy3 - XOR Vx, Vy
			// Vx = Vx xor Vy
			x := (cpu.Opcode & 0x0F00) >> 8
			y := (cpu.Opcode & 0x00F0) >> 4
			cpu.V[x] ^= cpu.V[y]

		case 0x4:
			// 8xy4 - ADD Vx, Vy
			// Set Vx = Vx + Vy, set VF = carry.
			x := (cpu.Opcode & 0x0F00) >> 8
			y := (cpu.Opcode & 0x00F0) >> 4

			result := uint16(cpu.V[x]) + uint16(cpu.V[y])

			// Go wraps silently, X and Y were not widened
			if result > 255 {
				// carry = 1
				cpu.V[0xF] = 1
			} else {
				cpu.V[0XF] = 0
			}
			cpu.V[x] = byte(result)

		case 0x5:
			// 8xy5 - SUB Vx, Vy
			// Set Vx = Vx - Vy, set VF = NOT borrow.
			x := (cpu.Opcode & 0x0F00) >> 8
			y := (cpu.Opcode & 0x00F0) >> 4

			result := uint16(cpu.V[x]) - uint16(cpu.V[y])
			if cpu.V[x] > cpu.V[y] {
				cpu.V[0xF] = 1
			} else {
				cpu.V[0xF] = 0
			}
			cpu.V[x] = byte(result)

		case 0x6:
			// 8xy6 - SHR Vx {, Vy}
			// Set Vx = Vx SHR 1.
			x := (cpu.Opcode & 0x0F00) >> 8

			cpu.V[0xF] = cpu.V[x] & 0x1
			cpu.V[x] >>= 1

		case 0x7:
			// 8xy7 - SUBN Vx, Vy
			// Set Vx = Vy - Vx, set VF = NOT borrow.
			x := (cpu.Opcode & 0x0F00) >> 8
			y := (cpu.Opcode & 0x00F0) >> 4

			result := uint16(cpu.V[y]) - uint16(cpu.V[x])
			if cpu.V[y] > cpu.V[x] {
				cpu.V[0xF] = 1
			} else {
				cpu.V[0xF] = 0
			}
			cpu.V[x] = byte(result)

		case 0xE:
			// 8xyE - SHL Vx {, Vy}
			// Set Vx = Vx SHL 1.
			x := (cpu.Opcode & 0x0F00) >> 8
			cpu.V[0xF] = cpu.V[x] & 0x80

			cpu.V[x] <<= 1
		}

	case 0x9000:
		// 9xy0 - SNE Vx, Vy
		// Skip next instruction if Vx != Vy.
		x := (cpu.Opcode & 0x0F00) >> 8
		y := (cpu.Opcode & 0x00F0) >> 4
		if cpu.V[x] != cpu.V[y] {
			cpu.PC += 2
		}

	case 0xA000:
		// Annn - LD I, addr
		// The value of register I is set to nnn
		cpu.IR = cpu.Opcode & 0xFFF

	case 0xB000:
		// Bnnn - JP V0, addr
		// Jump to location nnn + V0.
		cpu.PC = (cpu.Opcode & 0xFFF) + uint16(cpu.V[0])

	case 0xC000:
		// Cxkk - RND Vx, byte
		// Set Vx = random byte AND kk.
		x := (cpu.Opcode & 0x0F00) >> 8
		kk := cpu.Opcode & 0xFF
		n := byte(rand.Intn(256))

		cpu.V[x] = n & byte(kk)

	case 0xD000:
		// Dxyn - DRW Vx, Vy, nibble
		// Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.
		// Drawing is done by xor
		x := (cpu.Opcode & 0x0F00) >> 8
		y := (cpu.Opcode & 0x00F0) >> 4
		n := (cpu.Opcode & 0x00F)

		Vx := uint16(cpu.V[x])
		Vy := uint16(cpu.V[y])
		cpu.V[0xF] = 0

		for row := range n {
			spriteMem := cpu.Memory[cpu.IR + row]
			for col := range uint16(8) {
				bit := (spriteMem >> (7 - col)) & 1

				if bit == 1 {
					px := (Vx + col) % 64
					py := (Vy + row) % 32

					// XOR with this erases the pixel
					// Hence, V[F] = 1
					if cpu.Output[px][py] == 1 {
						cpu.V[0xF] = 1
					}
					cpu.Output[px][py] ^= bit
				}
			}
		}

	case 0xE000:
		switch (cpu.Opcode & 0xFF) {
		case 0x9E:
			// Ex9E - SKP Vx
			// Skip next instruction if key with the value of Vx is pressed.
			x := (cpu.Opcode & 0x0F00) >> 8
			Vx := cpu.V[x]

			// Pressed -> 1
			if cpu.Input[Vx] == 1 {
				cpu.PC += 2
			}

		case 0xA1:
			// ExA1 - SKNP Vx
			// Skip next instruction if key with the value of Vx is not pressed.
			x := (cpu.Opcode & 0x0F00) >> 8
			Vx := cpu.V[x]

			// Not Pressed == 0
			if cpu.Input[Vx] == 0 {
				cpu.PC += 2
			}
		}

	case 0xF000:
		switch (cpu.Opcode & 0xFF) {
		case 0x07:
			// Fx07 - LD Vx, DT
			// Set Vx = delay timer value.
			x := (cpu.Opcode & 0x0F00) >> 8
			cpu.V[x] = cpu.DelayTimer

		case 0x0A:
			// Fx0A - LD Vx, K
			// Wait for a key press, store the value of the key in Vx.
			x := (cpu.Opcode & 0x0F00) >> 8
			for key := range 16 {
				if cpu.Input[key] == 1 {
					cpu.V[x] = byte(key)
					return
				}
			}
			cpu.PC -= 2

		case 0x15:
			// Fx15 - LD DT, Vx
			// Set delay timer = Vx.
			x := (cpu.Opcode & 0x0F00) >> 8
			cpu.DelayTimer = cpu.V[x]

		case 0x18:
			// Fx18 - LD ST, Vx
			// Set sound timer = Vx.
			x := (cpu.Opcode & 0x0F00) >> 8
			cpu.SoundTimer = cpu.V[x]

		case 0x1E:
			// Fx1E - ADD I, Vx
			// Set I = I + Vx.
			x := (cpu.Opcode & 0x0F00) >> 8
			cpu.IR += uint16(cpu.V[x])

		case 0x29:
			// Fx29 - LD F, Vx
			// Set I = location of sprite for digit Vx.
			x := (cpu.Opcode & 0x0F00) >> 8
			num := cpu.V[x]

			// Sprites are 5 bytes each stored from 0x0000
			cpu.IR = uint16(num * 5) + 0x0000

		case 0x33:
			// Fx33 - LD B, Vx
			// Store BCD representation of Vx in memory locations I, I+1, and I+2.
			x := (cpu.Opcode & 0x0F00) >> 8
			num := cpu.V[x]

			for i := 2; i >= 0; i-- {
				cpu.Memory[cpu.IR + uint16(i)] = num % 10
				num /= 10
			}

		case 0x55:
			// Fx55 - LD [I], Vx
			// Store registers V0 through Vx in memory starting at location I.
			x := (cpu.Opcode & 0x0F00) >> 8

			for i := uint16(0); i <= x; i++ {
				cpu.Memory[cpu.IR + i] = cpu.V[i]
			}

		case 0x65:
			// Fx65 - LD Vx, [I]
			// Read registers V0 through Vx from memory starting at location I.
			x := (cpu.Opcode & 0x0F00) >> 8

			for i := uint16(0); i <= x; i++ {
				cpu.V[i] = cpu.Memory[cpu.IR + i]
			}
		}

	default:
		fmt.Printf("Invalid opcode: 0x%X\n", cpu.Opcode)
	}
}