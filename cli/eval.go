package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/senither/zen-lang/evaluator"
	"github.com/senither/zen-lang/lexer"
	"github.com/senither/zen-lang/objects"
	"github.com/senither/zen-lang/parser"
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
		if len(args) > 0 {
			content, err := loadFileContents(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}

			path, _ := filepath.Abs(args[0])
			runAndEval(content, path, nil)

			return
		}

		env := objects.NewEnvironment(nil)
		scanner := bufio.NewScanner(os.Stdin)

		fmt.Println("Welcome to the Zen REPL, type your code below to see the evaluated output.")
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

			runAndEval(line, nil, env)
		}
	},
}

func runAndEval(input string, filePath interface{}, env *objects.Environment) {
	lexer := lexer.New(input)
	parser := parser.New(lexer, filePath)

	program := parser.ParseProgram()
	if len(parser.Errors()) > 0 {
		printParseErrors(parser.Errors())
		return
	}

	if env == nil {
		env = objects.NewEnvironment(filePath)
	}

	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		fmt.Printf("%s\n", evaluated.Inspect())
	}
}
