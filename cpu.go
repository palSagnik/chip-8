package chip8

type CPU struct {
	R           		[16]byte            // 16 8-bit registers (R0 - RF)
	PC          		uint16              // Program Counter: stores the address of the current instruction
	SP          		byte                // Stack Pointer: refers to the top of the stack
	IR          		uint16              // I register: used to store memory address, last 12 bits used in general
	Stack       		[16]uint16          // Stack
	DelayTimer  		byte                // Delay Register
	SoundTimer       	byte                // Sound Register

	Memory 				[4096]byte          // Memory
	Output     		    [64][32]byte        // Output buffer: working as display
	Input 				[16]byte			// Input: Keypad mapping
}
