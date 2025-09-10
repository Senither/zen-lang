package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/senither/zen-lang/compiler"
	"github.com/senither/zen-lang/lexer"
	"github.com/senither/zen-lang/parser"
	"github.com/senither/zen-lang/vm"
	"github.com/spf13/cobra"
)

func init() {
	rootCommand.AddCommand(replCommand)
}

var replCommand = &cobra.Command{
	Use:   "repl",
	Short: "Run code and get the JIT-compiled output",
	Long:  "Runs the code provided and outputs the JIT-compiled result.",
	Run: func(cmd *cobra.Command, args []string) {
		scanner := bufio.NewScanner(os.Stdin)

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
			parser := parser.New(lexer, nil)

			program := parser.ParseProgram()
			if len(parser.Errors()) > 0 {
				printParseErrors(parser.Errors())
				continue
			}

			compiler := compiler.New()
			err := compiler.Compile(program)
			if err != nil {
				fmt.Println(err)
				continue
			}

			vm := vm.New(compiler.Bytecode())
			err = vm.Run()
			if err != nil {
				fmt.Println(err)
				continue
			}

			stackTop := vm.StackTop()
			if stackTop != nil {
				fmt.Printf("%s\n", stackTop.Inspect())
			}
		}
	},
}
