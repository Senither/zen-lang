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
	debugCommand.Flags().BoolP("verbose", "v", false, "Disables print capture and panic recoveries so failures show full stack traces.")

	rootCommand.AddCommand(debugCommand)
}

var debugCommand = &cobra.Command{
	Use:    "debug",
	Short:  "Takes a file as input and produces the bytecode, evaluated, and VM results.",
	Long:   "Runs the provided file and outputs the compiled bytecode, the result from the evaluator, and the result from the virtual machine.",
	Hidden: true,
	Args:   cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		table, globals, constants := createCompilerParameters()

		createREPLRunner(args, []string{}, func(input string, path any) {
			if !verbose {
				defer recoverFromPanic()
			}

			lexer := inputToLexer(input)
			program := lexerToProgram(lexer, path)
			if program == nil {
				return
			}

			compile := compiler.NewWithState(path, table, constants)
			compilerErr := compile.Compile(program)
			bytecode := compile.Bytecode()

			fmt.Println("=====[ Compiled Bytecode ]=====")
			if compilerErr != nil {
				fmt.Printf(BgRed+"\nCOMPILATION ERROR%s\n\n%s\n", Reset, compilerErr.Error())
			} else {
				fmt.Println(strings.TrimRight(bytecode.String(), "\n"))
			}

			evalStart := time.Now()
			evalRes := runAndReturnEvaluated(program, path)
			evalDuration := time.Since(evalStart)

			fmt.Printf("=====[ Evaluator Result (Time: %s) ]=====\n", evalDuration)
			fmt.Println(evalRes)

			var vmRes = " ~ Not executed due to compiler errors ~ "
			vmStart := time.Now()
			if compilerErr == nil {
				vmRes = runAndReturnVirtualMachineResult(verbose, bytecode, globals)
			}
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
	evaluator.Stdout.Clear()
	env := objects.NewEnvironment(path)

	evaluated := evaluator.Stdout.Mute(func() objects.Object {
		return evaluator.Eval(program, env)
	})

	var output strings.Builder

	if len(evaluator.Stdout.ReadAll()) > 0 {
		output.WriteString(strings.Join(evaluator.Stdout.ReadAll(), ""))
	}

	if evaluated == nil {
		return output.String()
	}

	switch evaluated := evaluated.(type) {
	case *objects.Error:
		output.WriteString(evaluated.Inspect())

	default:
		if output.Len() == 0 {
			output.WriteString(evaluated.Inspect())
		}
	}

	return output.String()
}

func runAndReturnVirtualMachineResult(verbose bool, bytecode *compiler.Bytecode, globals []objects.Object) string {
	vm.Stdout.Clear()

	machine := vm.NewWithGlobalsStore(bytecode, globals)
	if !verbose {
		machine.EnableStdoutCapture()
	}

	result := vm.Stdout.Mute(func() objects.Object {
		if err := machine.Run(); err != nil {
			return objects.NativeErrorToErrorObject(err)
		}

		return machine.LastPoppedStackElem()
	})

	var output strings.Builder

	if len(vm.Stdout.ReadAll()) > 0 {
		output.WriteString(strings.Join(vm.Stdout.ReadAll(), ""))
	}

	if result == nil {
		return output.String()
	}

	switch result := result.(type) {
	case *objects.Error:
		output.WriteString(result.Inspect())

	default:
		if output.Len() == 0 {
			output.WriteString(result.Inspect())
		}
	}

	return output.String()
}

func recoverFromPanic() {
	if r := recover(); r != nil {
		fmt.Println("")
		fmt.Println("RECOVERED FROM PANIC DURING RUNTIME")
		fmt.Println("ERROR: " + fmt.Sprint(r))
	}
}
