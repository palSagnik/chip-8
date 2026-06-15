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

	workList := []int{0}
	visited := map[int]bool{}
	for len(workList) > 0 {
		address := workList[0]
		workList = workList[1:]

		// if visited then skip; else mark as true
		if visited[address] {
			continue
		}
		visited[address] = true
		if address+1 >= len(rom) {
			break
		}

		// decode inst
		inst := uint16(rom[address]) << 8
		inst |= uint16(rom[address + 1])
		disasmInst := disassembleInstruction(inst)
		disasmOutput = append(disasmOutput, disasmObject{
			Memory: 0x200 + uint16(address),
			Opcode: inst,
			Instruction: disasmInst,
		})


		// pop the old address put the next one
		memAddr := 0x200 + uint16(address)
		nextAddress := findNextAddress(inst, memAddr)
		// RET: no more address after here, but worklist might have address
		if len(nextAddress) == 0 {
			continue
		}
		for _, addr := range nextAddress {
			workList = append(workList, int(addr - 0x200))
		}
	}

	for _, output := range disasmOutput {
		fmt.Printf("0x%X: 0x%04X -> %s\n", output.Memory, output.Opcode, output.Instruction)
	}

	return nil
}

func findNextAddress(inst uint16, address uint16) []uint16 {
	switch inst & 0xF000 {
	case 0x0000:
		if inst & 0xFF == 0xEE {
			return []uint16{}
		}
	case 0x1000:
		return []uint16{inst & 0xFFF}

	// CALL has two branches, either it goes to the next after RET or jumps to addr
	case 0x2000:
		return []uint16{inst & 0xFFF, address + 2}
	case 0x3000:
		return []uint16{address + 2, address + 4}
	case 0x4000:
		return []uint16{address + 2, address + 4}
	case 0x5000:
		return []uint16{address + 2, address + 4}
	case 0x9000:
		return []uint16{address + 2, address + 4}
	case 0xE000:
		return []uint16{address + 2, address + 4}
	}

	return []uint16{address + 2}
}