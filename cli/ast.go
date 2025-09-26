package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCommand.AddCommand(astCommand)
}

var astCommand = &cobra.Command{
	Use:   "ast",
	Short: "Run code and get the AST output",
	Long:  "Runs the code provided and outputs the AST it generates.",
	Run: func(cmd *cobra.Command, args []string) {
		createREPLRunner(args, []string{
			"Welcome to the Zen AST generator, type your code below to see the AST output.",
			"Type 'exit' to exit the AST generator or press Ctrl+C.",
		}, func(input string, path any) {
			lexer := inputToLexer(input)

			if program := lexerToProgram(lexer, path); program != nil {
				for _, statement := range program.Statements {
					fmt.Printf("%s\n", statement.String())
				}
			}
		})
	},
}
