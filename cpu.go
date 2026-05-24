package chip8

import "fmt"

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
		switch cpu.Opcode & 0x00FF {
		case 0x00E0:
			// 00E0 - CLS
			// Clear the display
			cpu.Output = [64][32]byte{}

		// TODO: properly think about Stack instructiions
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
		kk := cpu.Opcode & 0x00FF
		if cpu.V[x] == byte(kk) {
			cpu.PC += 2
		}

	case 0x4000:
		// 4xkk - SNE Vx, byte
		// Skip next instruction if Vx != kk.
		x := (cpu.Opcode & 0x0F00) >> 8
		kk := cpu.Opcode & 0x00FF
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
		kk := cpu.Opcode & 0x00FF
		cpu.V[x] = byte(kk)

	case 0x7000:
		// 7xkk - ADD Vx, byte
		// Set Vx = Vx + kk.
		x := (cpu.Opcode & 0x0F00) >> 8
		kk := cpu.Opcode & 0x00FF
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
	
	case 0xA000:
		// Annn - LD I, addr
		// The value of register I is set to nnn
		cpu.IR = cpu.Opcode & 0xFFF

	default:
		fmt.Printf("Invalid opcode: 0x%X\n", cpu.Opcode)
	}
}