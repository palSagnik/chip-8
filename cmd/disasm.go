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
	for i := 0; i < len(rom) - 1; i += 2 {
		inst := uint16(rom[i]) << 8
		inst |= uint16(rom[i + 1])

		disasmOutput = append(disasmOutput, disasmObject{
			Memory: uint16(0x200 + i),
			Opcode: inst,
			Instruction: disassembleInstruction(inst),
		})
	}

	for i := range disasmOutput {
		fmt.Printf("0x%X: 0x%X -> %s\n", disasmOutput[i].Memory, disasmOutput[i].Opcode, disasmOutput[i].Instruction)
	}
	return nil
}
