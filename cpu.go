package chip8

type CPU struct {
	V          [16]byte   // 16 8-bit registers (V0 - VF)
	PC         uint16     // Program Counter: stores the address of the current instruction
	SP         byte       // Stack Pointer: refers to the top of the stack
	IR         uint16     // I register: used to store memory address, last 12 bits used in general
	Stack      [16]uint16 // Stack
	DelayTimer byte       // Delay Register
	SoundTimer byte       // Sound Register

	Memory [4096]byte   // Memory
	Output [64][32]byte // Output buffer: working as display
	Input  [16]byte     // Input: Keypad mapping

	Instruction uint16 // current instruction
}

func (cpu *CPU) Init() {
	// clear registers
	for i := range 16 {
		cpu.V[i] = 0
	}
	cpu.DelayTimer = 0
	cpu.SoundTimer = 0
	cpu.IR = 0

	// clear memory, stack
	cpu.Memory = [4096]byte{}
	cpu.Stack = [16]uint16{}

	// clear input, output
	cpu.Output = [64][32]byte{}
	cpu.Input = [16]byte{}

	// loading fontset
	copy(cpu.Memory[:], fontset[:])

	// setting PC
	cpu.PC = 0x200
}

func (cpu *CPU) Fetch() {
	// update instruction
	cpu.Instruction = uint16(cpu.Memory[cpu.PC])
	cpu.Instruction = cpu.Instruction << 8
	cpu.Instruction |= uint16(cpu.Memory[cpu.PC+1])

	// after succesful fetch point to the next instruction
	cpu.PC += 2
}

func (cpu *CPU) Decode() {
	switch cpu.Instruction & 0xF000 {
	case 0x0000:
		switch cpu.Instruction & 0xFF {
		case 0x00E0:
			cpu.__00E0()
		case 0x00EE:
			cpu.__00EE()
		}

	case 0x1000:
		cpu.__1nnn()
	case 0x2000:
		cpu.__2nnn()
	case 0x3000:
		cpu.__3xkk()
	case 0x4000:
		cpu.__4xkk()
	case 0x5000:
		cpu.__5xy0()
	case 0x6000:
		cpu.__6xkk()
	case 0x7000:
		cpu.__7xkk()

	case 0x8000:
		switch cpu.Instruction & 0x000F {
		case 0x0:
			cpu.__8xy0()
		case 0x1:
			cpu.__8xy1()
		case 0x2:
			cpu.__8xy2()
		case 0x3:
			cpu.__8xy3()
		case 0x4:
			cpu.__8xy4()
		case 0x5:
			cpu.__8xy5()
		case 0x6:
			cpu.__8xy6()
		case 0x7:
			cpu.__8xy7()
		case 0xE:
			cpu.__8xyE()
		}

	case 0x9000:
		cpu.__9xy0()
	case 0xA000:
		cpu.__Annn()
	case 0xB000:
		cpu.__Bnnn()
	case 0xC000:
		cpu.__Cxkk()
	case 0xD000:
		cpu.__Dxyn()

	case 0xE000:
		switch cpu.Instruction & 0xFF {
		case 0x9E:
			cpu.__Ex9E()
		case 0xA1:
			cpu.__ExA1()
		}

	case 0xF000:
		switch cpu.Instruction & 0xFF {
		case 0x07:
			cpu.__Fx07()
		case 0x0A:
			cpu.__Fx0A()
		case 0x15:
			cpu.__Fx15()
		case 0x18:
			cpu.__Fx18()
		case 0x1E:
			cpu.__Fx1E()
		case 0x29:
			cpu.__Fx29()
		case 0x33:
			cpu.__Fx33()
		case 0x55:
			cpu.__Fx55()
		case 0x65:
			cpu.__Fx65()
		}

	default:
		cpu.__invalid()
	}
}