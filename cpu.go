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
}

func (cpu *CPU) Decode() {
	switch (cpu.Opcode & 0xF000) {
	
	case 0x0000:
		switch cpu.Opcode & 0x00FF {
		case 0x00E0:
			// 00E0 - CLS
			// Clear the display
			cpu.Output = [64][32]byte{}
			cpu.PC += 2
		case 0x00EE:
			// 00EE - RET
            // Return from a subroutine.
			cpu.PC = cpu.Stack[cpu.SP]

			// TODO: clear the top of the stack?
			// then decrement the stack pointer
			cpu.SP -= 1
			cpu.PC += 2
		}
	case 0x1000:
		// 1nnn - JP addr
		// Jump to location nnn.
		cpu.PC = cpu.Opcode & 0xFFF
		
	case 0xA000:
		// Annn - LD I, addr
		// The value of register I is set to nnn
		cpu.IR = cpu.Opcode & 0xFFF
		cpu.PC += 2
	
	default:
		fmt.Printf("Invalid opcode: 0x%X\n", cpu.Opcode)
	}
}