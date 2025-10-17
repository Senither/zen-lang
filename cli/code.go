package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCommand.AddCommand(codeCommand)
}

var codeCommand = &cobra.Command{
	Use:   "code",
	Short: "Run code and get the bytecode instructions as output",
	Long:  "Runs the code provided and outputs the bytecode instructions it generates.",
	Run: func(cmd *cobra.Command, args []string) {
		table, _, constants := createCompilerParameters()

		createREPLRunner(args, []string{
			"Welcome to the Zen Bytecode generator, type your code below to see the Bytecode output.",
			"Type 'exit' to exit the Bytecode generator or press Ctrl+C.",
		}, func(input string, path any) {
			lexer := inputToLexer(input)
			program := lexerToProgram(lexer, path)

			if bytecode := programToBytecode(program, table, constants); bytecode != nil {
				fmt.Print(bytecode.String())
			}
		})
	},
}
