package cli

import (
	"fmt"
	"os"

	"github.com/senither/zen-lang/evaluator"
	"github.com/senither/zen-lang/lexer"
	"github.com/senither/zen-lang/objects"
	"github.com/senither/zen-lang/parser"
	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:        "zen",
	Short:      "Run Zen code",
	Long:       "Runs the Zen interpreter",
	Args:       cobra.MinimumNArgs(1),
	ArgAliases: []string{"file"},
	ValidArgs:  []string{"file"},
	Run: func(cmd *cobra.Command, args []string) {
		content, err := loadFileContents(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}

		lexer := lexer.New(string(content))
		parser := parser.New(lexer)

		program := parser.ParseProgram()
		if len(parser.Errors()) > 0 {
			for _, err := range parser.Errors() {
				fmt.Println(err.String())
			}
			return
		}

		env := objects.NewEnvironment()
		evaluated := evaluator.Eval(program, env)
		if evaluated == nil {
			fmt.Println("Failed to evaluate program, evaluation returned nil")
			return
		}

		if evaluated.Type() == objects.ERROR_OBJ {
			fmt.Println(evaluated.Inspect())
		}
	},
}

func loadFileContents(file string) (string, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file %q: %w", file, err)
	}
	return string(content), nil
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing command: %s\n", err)
		os.Exit(1)
	}
}
