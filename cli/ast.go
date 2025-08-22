package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/senither/zen-lang/lexer"
	"github.com/senither/zen-lang/parser"
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
		scanner := bufio.NewScanner(os.Stdin)

		fmt.Println("Welcome to the Zen AST generator, type your code below to see the AST output.")
		fmt.Println("Type 'exit' to exit the AST generator or press Ctrl+C.")

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
			parser := parser.New(lexer)

			program := parser.ParseProgram()

			for _, statement := range program.Statements {
				fmt.Printf("%s\n", statement.String())
			}
		}
	},
}
