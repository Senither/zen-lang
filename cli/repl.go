package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/senither/zen-lang/evaluator"
	"github.com/senither/zen-lang/lexer"
	"github.com/senither/zen-lang/objects"
	"github.com/senither/zen-lang/parser"
	"github.com/spf13/cobra"
)

func init() {
	rootCommand.AddCommand(replCommand)
}

var replCommand = &cobra.Command{
	Use:   "repl",
	Short: "Start a Read, Eval, Print, Loop",
	Long:  "Creates an environment for evaluating Zen code.",
	Run: func(cmd *cobra.Command, args []string) {
		scanner := bufio.NewScanner(os.Stdin)
		env := objects.NewEnvironment(nil)

		fmt.Println("Welcome to the Zen REPL, type your code below to see the output.")
		fmt.Println("Type 'exit' to exit the REPL or press Ctrl+C.")

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
			if len(parser.Errors()) > 0 {
				for _, err := range parser.Errors() {
					fmt.Println(err.String())
				}
				continue
			}

			evaluated := evaluator.Eval(program, env)
			if evaluated != nil {
				fmt.Printf("%s\n", evaluated.Inspect())
			}
		}
	},
}
