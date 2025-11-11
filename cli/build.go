package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCommand.AddCommand(buildCommand)
	buildCommand.Flags().StringP("output", "o", "", "Output file name (default: same as input with .zenb extension)")
}

var buildCommand = &cobra.Command{
	Use:   "build <file>",
	Short: "Compile a Zen source file to bytecode",
	Long:  "Compiles the specified Zen source file and outputs the bytecode to a .zenb file.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		inputFile := args[0]

		content, err := os.ReadFile(inputFile)
		if err != nil {
			fmt.Printf("Error reading file '%s': %v\n", inputFile, err)
			os.Exit(1)
		}

		outputFile, _ := cmd.Flags().GetString("output")
		if outputFile == "" {
			ext := filepath.Ext(inputFile)
			outputFile = strings.TrimSuffix(inputFile, ext) + ".zenb"
		}

		lexer := inputToLexer(string(content))
		program := lexerToProgram(lexer, inputFile)
		if program == nil {
			fmt.Printf("\nFailed to parse file '%s'\n", inputFile)
			os.Exit(1)
		}

		table, _, constants := createCompilerParameters()
		bytecode := programToBytecode(inputFile, program, table, constants)
		if bytecode == nil {
			fmt.Printf("\nFailed to compile file '%s'\n", inputFile)
			os.Exit(1)
		}

		err = os.WriteFile(outputFile, bytecode.Serialize(), 0644)
		if err != nil {
			fmt.Printf("\nError writing to file '%s': %v\n", outputFile, err)
			os.Exit(1)
		}

		fmt.Printf("Successfully compiled '%s' to '%s'\n", inputFile, outputFile)
	},
}
