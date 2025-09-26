package cli

import (
	"fmt"

	"github.com/senither/zen-lang/tokens"
	"github.com/spf13/cobra"
)

func init() {
	rootCommand.AddCommand(lexerCommand)
}

var lexerCommand = &cobra.Command{
	Use:   "lexer",
	Short: "Run code and get the lexer output",
	Long:  "Runs the code provided to the lexer and outputs the tokens it generates.",
	Run: func(cmd *cobra.Command, args []string) {
		createREPLRunner(args, []string{
			"Welcome to the Zen lexer, type your code below to see the lexer output.",
			"Type 'exit' to exit the lexer or press Ctrl+C.",
		}, func(input string, _ any) {
			lexer := inputToLexer(input)

			for tok := lexer.NextToken(); tok.Type != tokens.EOF; tok = lexer.NextToken() {
				fmt.Printf("%+v\n", tok)
			}
		})
	},
}
