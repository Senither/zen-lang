package cli

import (
	"fmt"

	"github.com/senither/zen-lang/vm"
	"github.com/spf13/cobra"
)

func init() {
	rootCommand.AddCommand(replCommand)
}

var replCommand = &cobra.Command{
	Use:   "repl",
	Short: "Run code and get the JIT-compiled output",
	Long:  "Runs the code provided and outputs the JIT-compiled result.",
	Run: func(cmd *cobra.Command, args []string) {
		table, globals, constants := createCompilerParameters()

		createREPLRunner(args, []string{
			"Welcome to the Zen REPL, type your code below to see the evaluated output.",
			"Type 'exit' to exit the REPL or press Ctrl+C.",
		}, func(input string, path any) {
			lexer := inputToLexer(input)
			program := lexerToProgram(lexer, path)
			bytecode := programToBytecode(program, table, constants)
			if bytecode == nil {
				return
			}

			vm := vm.NewWithGlobalsStore(bytecode, globals)
			if err := vm.Run(); err != nil {
				fmt.Println(err)
				return
			}

			stackTop := vm.LastPoppedStackElem()
			if stackTop != nil {
				fmt.Printf("%s\n", stackTop.Inspect())
			}
		})
	},
}
