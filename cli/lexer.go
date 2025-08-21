package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/senither/zen-lang/lexer"
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
		scanner := bufio.NewScanner(os.Stdin)

		fmt.Println("Welcome to the Zen lexer, type your code below to see the lexer output.")
		fmt.Println("Type 'exit' to exit the lexer or press Ctrl+C.")

		for {
			fmt.Printf(">>> ")

			scanned := scanner.Scan()
			if !scanned {
				return
			}

			line := scanner.Text()
			if line == "exit" {
				return
			}

			lexer := lexer.New(line)

			for tok := lexer.NextToken(); tok.Type != tokens.EOF; tok = lexer.NextToken() {
				fmt.Printf("%+v\n", tok)
			}
		}
	},
}
