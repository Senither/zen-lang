package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

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
		if len(args) > 0 {
			content, err := loadFileContents(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}

			path, _ := filepath.Abs(args[0])
			runAndEvalAST(content, path)

			return
		}

		fmt.Printf("Args: %v\n", args)

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

			runAndEvalAST(line, nil)
		}
	},
}

func runAndEvalAST(input string, filePath interface{}) {
	lexer := lexer.New(input)
	parser := parser.New(lexer, filePath)

	program := parser.ParseProgram()

	if len(parser.Errors()) > 0 {
		for _, err := range parser.Errors() {
			fmt.Println(err.String())
		}
		return
	}

	for _, statement := range program.Statements {
		fmt.Printf("%s\n", statement.String())
	}
}
