package cli

import (
	"fmt"

	"github.com/senither/zen-lang/evaluator"
	"github.com/senither/zen-lang/objects"
	"github.com/spf13/cobra"
)

func init() {
	rootCommand.AddCommand(evalCommand)
}

var evalCommand = &cobra.Command{
	Use:   "eval",
	Short: "Run code and get the evaluated output",
	Long:  "Runs the code provided and outputs the evaluated result.",
	Run: func(cmd *cobra.Command, args []string) {
		env := objects.NewEnvironment(nil)

		createREPLRunner(args, []string{
			"Welcome to the Zen REPL, type your code below to see the evaluated output.",
			"Type 'exit' to exit the REPL or press Ctrl+C.",
		}, func(input string, path any) {
			lexer := inputToLexer(input)
			program := lexerToProgram(lexer, path)

			if path != nil {
				env = objects.NewEnvironment(path)
			}

			evaluated := evaluator.Eval(program, env)
			if evaluated != nil && evaluated.Type() != objects.NULL_OBJ {
				fmt.Printf("%s\n", evaluated.Inspect())
			}
		})
	},
}
