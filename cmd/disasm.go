package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var disasmFile string

var disasmCmd = &cobra.Command{
	Use:   "disasm",
	Short: "Disassemble a CHIP-8 ROM",
	RunE:  runDisasm,
}

func init() {
	disasmCmd.Flags().StringVarP(&disasmFile, "file", "f", "", "path to ROM file (required)")
	disasmCmd.MarkFlagRequired("file")
}

type disasmObject struct {
	Memory uint16
	Opcode uint16
	Instruction string
}

func runDisasm(cmd *cobra.Command, args []string) error {
	rom, err := os.ReadFile(disasmFile)
	if err != nil {
		log.Fatalf("could not read ROM: %v", err)
	}

	disasmOutput := make([]disasmObject, 0)

	// recursive descent: start at ROM entry point (offset 0 = address 0x200)
	// and follow control flow rather than reading linearly
	workList := []int{0}
	visited := map[int]bool{}
	for len(workList) > 0 {
		address := workList[0]
		workList = workList[1:]

		if visited[address] {
			continue
		}
		visited[address] = true
		if address+1 >= len(rom) {
			break
		}

		inst := uint16(rom[address])<<8 | uint16(rom[address+1])
		disasmOutput = append(disasmOutput, disasmObject{
			Memory:      0x200 + uint16(address),
			Opcode:      inst,
			Instruction: disassembleInstruction(inst),
		})

		// determine which addresses are reachable from this instruction
		// and add them to the worklist for future processing
		memAddr := 0x200 + uint16(address)
		nextAddresses := findNextAddress(inst, memAddr)
		for _, addr := range nextAddresses {
			workList = append(workList, int(addr-0x200))
		}
	}

	for _, output := range disasmOutput {
		fmt.Printf("0x%X: 0x%04X -> %s\n", output.Memory, output.Opcode, output.Instruction)
	}

	return nil
}

// findNextAddress returns all addresses reachable from the current instruction.
// This is the core of recursive descent: each instruction type produces different
// successor addresses, and we must follow all of them to avoid missing code.
func findNextAddress(inst uint16, address uint16) []uint16 {
	switch inst & 0xF000 {
	case 0x0000:
		// RET has no static successor — return address is on the stack at runtime.
		// The caller (CALL) already pushed address+2 to the worklist.
		if inst&0xFF == 0xEE {
			return nil
		}
	case 0x1000:
		// JP nnn — unconditional jump, only successor is the target
		return []uint16{inst & 0xFFF}
	case 0x2000:
		// CALL nnn — jumps to subroutine AND continues at address+2 after RET
		return []uint16{inst & 0xFFF, address + 2}
	case 0x3000, 0x4000, 0x5000, 0x9000:
		// conditional skip — both branches are possible at runtime
		return []uint16{address + 2, address + 4}
	case 0xE000:
		// SKP/SKNP — key-based conditional skip, same as above
		return []uint16{address + 2, address + 4}
	}

	// all other instructions execute sequentially
	return []uint16{address + 2}
}