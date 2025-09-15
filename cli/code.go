package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/senither/zen-lang/compiler"
	"github.com/senither/zen-lang/lexer"
	"github.com/senither/zen-lang/parser"
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

		scanner := bufio.NewScanner(os.Stdin)

		fmt.Println("Welcome to the Zen Bytecode generator, type your code below to see the Bytecode output.")
		fmt.Println("Type 'exit' to exit the Bytecode generator or press Ctrl+C.")

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

			runAndEvalBytecode(line, nil)
		}
	},
}

func runAndEvalBytecode(input string, filePath interface{}) {
	lexer := lexer.New(input)
	parser := parser.New(lexer, filePath)

	program := parser.ParseProgram()

	if len(parser.Errors()) > 0 {
		printParseErrors(parser.Errors())
		return
	}

	compiler := compiler.New()
	err := compiler.Compile(program)
	if err != nil {
		fmt.Println("Compilation error:", err)
		return
	}

	fmt.Print(compiler.Bytecode().Instructions.String())
}
