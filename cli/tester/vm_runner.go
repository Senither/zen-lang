package tester

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/senither/zen-lang/ast"
	"github.com/senither/zen-lang/compiler"
	"github.com/senither/zen-lang/objects"
	"github.com/senither/zen-lang/vm"
)

func (tr *TestRunner) runVMTest(test *Test, program *ast.Program, fullPath, file string) {
	start := time.Now()

	c := compiler.New(file)
	compilerErr := c.Compile(program)
	tr.setTiming(CompilationTiming, tr.getTiming(CompilationTiming)+time.Since(start))

	if compilerErr != nil {
		tr.compareCompilerErrorWithExpected(test, fullPath, compilerErr.Error())
		return
	}

	vm.Stdout.Clear()
	start = time.Now()
	result := vm.Stdout.Mute(func() objects.Object {
		runner := vm.New(c.Bytecode())
		runner.EnableStdoutCapture()

		if runErr := runner.Run(); runErr != nil {
			return objects.NativeErrorToErrorObject(runErr)
		}

		return runner.LastPoppedStackElem()
	})
	tr.setTiming(VMExecutionTiming, tr.getTiming(VMExecutionTiming)+time.Since(start))

	if objects.IsError(result) {
		tr.compareCompilerErrorWithExpected(test, fullPath, result.Inspect())
	} else if result != nil && result.Type() != objects.NULL_OBJ {
		tr.compareResultWithExpected(test, fullPath, result)
	} else {
		tr.compareStandardOutputWithExpectedVM(test, fullPath)
	}
}

func (tr *TestRunner) compareCompilerErrorWithExpected(test *Test, fullPath string, errorMessage string) {
	err := tr.stripFileLocationsFromError(test, fullPath, errorMessage)

	if err != test.errors {
		tr.printErrorStatusMessage(
			test,
			fullPath,
			fmt.Sprintf(
				"%s\n     -----------------[ RESULT ]-----------------\n%s\n     ----------------[ EXPECTED ]-----------------\n%s",
				"Test expectation does not match the compiler error",
				err,
				test.errors,
			),
			VirtualMachineEngine,
		)
		return
	}

	tr.printSuccessStatusMessage(test, VirtualMachineEngine)
}

func (tr *TestRunner) compareResultWithExpected(test *Test, fullPath string, result objects.Object) {
	value := tr.stripPointerLocationsFromContent(strings.Trim(result.Inspect(), "\n"))

	if value != test.expect {
		tr.printErrorStatusMessage(
			test,
			fullPath,
			fmt.Sprintf(
				"%s\n     -----------------[ RESULT ]-----------------\n%s\n     ----------------[ EXPECTED ]-----------------\n%s",
				"Test expectation does not match the evaluated result",
				value,
				test.expect,
			),
			VirtualMachineEngine,
		)
		return
	}

	tr.printSuccessStatusMessage(test, VirtualMachineEngine)
}

func (tr *TestRunner) compareStandardOutputWithExpectedVM(test *Test, fullPath string) {
	messages := vm.Stdout.ReadAll()
	if len(messages) == 0 {
		tr.printErrorStatusMessage(test, fullPath, "No output captured from standard output", VirtualMachineEngine)
		return
	}

	comparison := test.expect
	if comparison == "" {
		comparison = test.errors
	}

	out := tr.stripPointerLocationsFromContent(strings.Trim(strings.Join(messages, ""), "\n"))
	if out != comparison {
		tr.printErrorStatusMessage(
			test,
			fullPath,
			fmt.Sprintf(
				"%s\n     -----------------[ RESULT ]-----------------\n%s\n     ----------------[ EXPECTED ]-----------------\n%s",
				"Test expectation does not match the standard output",
				out,
				comparison,
			),
			VirtualMachineEngine,
		)
		return
	}

	tr.printSuccessStatusMessage(test, VirtualMachineEngine)
}

func (tr *TestRunner) stripPointerLocationsFromContent(content string) string {
	r, _ := regexp.Compile(`Closure\[0x[a-zA-Z0-9]+\]`)

	return r.ReplaceAllString(content, "Closure[<pointer>]")
}
