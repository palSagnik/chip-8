package main

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "chip8",
	Short: "CHIP-8 emulator and disassembler",
}

func init() {
	rootCmd.AddCommand(playCmd)
	rootCmd.AddCommand(disasmCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
