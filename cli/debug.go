package cli

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/senither/zen-lang/ast"
	"github.com/senither/zen-lang/compiler"
	"github.com/senither/zen-lang/evaluator"
	"github.com/senither/zen-lang/objects"
	"github.com/senither/zen-lang/vm"
	"github.com/spf13/cobra"
)

func init() {
	rootCommand.AddCommand(debugCommand)
}

var debugCommand = &cobra.Command{
	Use:    "debug",
	Short:  "Takes a file as input and produces the bytecode, evaluated, and VM results.",
	Long:   "Runs the provided file and outputs the compiled bytecode, the result from the evaluator, and the result from the virtual machine.",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide a file that should be run within the debugger.")
			return
		}

		table, globals, constants := createCompilerParameters()

		createREPLRunner(args, []string{}, func(input string, path any) {
			defer recoverFromPanic()

			lexer := inputToLexer(input)
			program := lexerToProgram(lexer, path)
			bytecode := programToBytecode(program, table, constants)
			if bytecode == nil {
				return
			}

			fmt.Println("=====[ Compiled Bytecode ]=====")
			fmt.Println(strings.TrimRight(bytecode.String(), "\n"))

			evalStart := time.Now()
			evalRes := runAndReturnEvaluated(program, path)
			evalDuration := time.Since(evalStart)

			fmt.Printf("=====[ Evaluator Result (Time: %s) ]=====\n", evalDuration)
			fmt.Println(evalRes)

			vmStart := time.Now()
			vmRes := runAndReturnVirtualMachineResult(bytecode, globals)
			vmDuration := time.Since(vmStart)

			fmt.Printf("=====[ Virtual Machine Result (Time: %s) ]=====\n", vmDuration)
			fmt.Println(vmRes)

			if evalRes != vmRes {
				fmt.Println("")
				fmt.Println(Red + "MISMATCH BETWEEN EVALUATOR AND VM RESULTS" + Reset)
				fmt.Println("")

				os.Exit(1)
			}
		})
	},
}

func runAndReturnEvaluated(program *ast.Program, path any) string {
	env := objects.NewEnvironment(path)

	evaluated := evaluator.Eval(program, env)
	if evaluated == nil {
		return ""
	}

	return evaluated.Inspect()
}

func runAndReturnVirtualMachineResult(bytecode *compiler.Bytecode, globals []objects.Object) string {
	vm := vm.NewWithGlobalsStore(bytecode, globals)
	if err := vm.Run(); err != nil {
		return err.Error()
	}

	stackTop := vm.LastPoppedStackElem()
	if stackTop == nil {
		return ""
	}

	return stackTop.Inspect()
}

func recoverFromPanic() {
	if r := recover(); r != nil {
		fmt.Println("")
		fmt.Println("RECOVERED FROM PANIC DURING RUNTIME")
		fmt.Println("ERROR: " + fmt.Sprint(r))
	}
}
