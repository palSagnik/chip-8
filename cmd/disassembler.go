package main

import "fmt"

func disassembleInstruction(instruction uint16) string {
	nnn := instruction & 0xFFF
	x := (instruction & 0x0F00) >> 8
	y := (instruction & 0x00F0) >> 4
	kk := instruction & 0xFF
	n := instruction & 0xF

	switch instruction & 0xF000 {
	case 0x0000:
		switch instruction & 0xFF {
		case 0x00E0:
			return "CLS"
		case 0x00EE:
			return "RET"
		}

	case 0x1000:
		return fmt.Sprintf("JP 0x%03X", nnn)
	case 0x2000:
		return fmt.Sprintf("CALL 0x%03X", nnn)
	case 0x3000:
		return fmt.Sprintf("SE V%X, 0x%02X", x, kk)
	case 0x4000:
		return fmt.Sprintf("SNE V%X, 0x%02X", x, kk)
	case 0x5000:
		return fmt.Sprintf("SE V%X, V%X", x, y)
	case 0x6000:
		return fmt.Sprintf("LD V%X, 0x%02X", x, kk)
	case 0x7000:
		return fmt.Sprintf("ADD V%X, 0x%02X", x, kk)

	case 0x8000:
		switch instruction & 0x000F {
		case 0x0:
			return fmt.Sprintf("LD V%X, V%X", x, y)
		case 0x1:
			return fmt.Sprintf("OR V%X, V%X", x, y)
		case 0x2:
			return fmt.Sprintf("AND V%X, V%X", x, y)
		case 0x3:
			return fmt.Sprintf("XOR V%X, V%X", x, y)
		case 0x4:
			return fmt.Sprintf("ADD V%X, V%X", x, y)
		case 0x5:
			return fmt.Sprintf("SUB V%X, V%X", x, y)
		case 0x6:
			return fmt.Sprintf("SHR V%X {, V%X}", x, y)
		case 0x7:
			return fmt.Sprintf("SUBN V%X, V%X", x, y)
		case 0xE:
			return fmt.Sprintf("SHL V%X {, V%X}", x, y)
		}

	case 0x9000:
		return fmt.Sprintf("SNE V%X, V%X", x, y)
	case 0xA000:
		return fmt.Sprintf("LD I, 0x%03X", nnn)
	case 0xB000:
		return fmt.Sprintf("JP V0, 0x%03X", nnn)
	case 0xC000:
		return fmt.Sprintf("RND V%X, 0x%02X", x, kk)
	case 0xD000:
		return fmt.Sprintf("DRW V%X, V%X, %d", x, y, n)

	case 0xE000:
		switch instruction & 0xFF {
		case 0x9E:
			return fmt.Sprintf("SKP V%X", x)
		case 0xA1:
			return fmt.Sprintf("SKNP V%X", x)
		}

	case 0xF000:
		switch instruction & 0xFF {
		case 0x07:
			return fmt.Sprintf("LD V%X, DT", x)
		case 0x0A:
			return fmt.Sprintf("LD V%X, K", x)
		case 0x15:
			return fmt.Sprintf("LD DT, V%X", x)
		case 0x18:
			return fmt.Sprintf("LD ST, V%X", x)
		case 0x1E:
			return fmt.Sprintf("ADD I, V%X", x)
		case 0x29:
			return fmt.Sprintf("LD F, V%X", x)
		case 0x33:
			return fmt.Sprintf("LD B, V%X", x)
		case 0x55:
			return fmt.Sprintf("LD [I], V%X", x)
		case 0x65:
			return fmt.Sprintf("LD V%X, [I]", x)
		}
	}

	return fmt.Sprintf("UNKNOWN 0x%04X", instruction)
}
