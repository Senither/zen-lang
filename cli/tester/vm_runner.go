package tester

import (
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
	timeTaken := time.Since(start)

	tr.setTiming(CompilationTiming, tr.getTiming(CompilationTiming)+timeTaken)
	test.metadata[CompilationTiming] = timeTaken

	if compilerErr != nil {
		tr.compareCompliedVMWithError(test, fullPath, compilerErr.Error())
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
	timeTaken = time.Since(start)

	tr.setTiming(VMExecutionTiming, tr.getTiming(VMExecutionTiming)+timeTaken)
	test.metadata[VMExecutionTiming] = timeTaken

	if objects.IsError(result) {
		tr.compareCompliedVMWithError(test, fullPath, result.Inspect())
	} else if result != nil && result.Type() != objects.NULL_OBJ {
		tr.compareCompliedVMWithExpected(test, fullPath, result)
	} else {
		tr.compareCompliedVMWithStandardOutput(test, fullPath)
	}
}

func (tr *TestRunner) compareCompliedVMWithError(test *Test, fullPath string, errorMessage string) {
	err := tr.normalizeFileLocations(errorMessage)

	if err != test.errors {
		tr.printErrorDoesNotMatchExpectation(
			test, fullPath,
			"Test expectation does not match the compiler error",
			err, test.errors,
			VirtualMachineEngine,
		)
		return
	}

	tr.printSuccessStatusMessage(test, VirtualMachineEngine)
}

func (tr *TestRunner) compareCompliedVMWithExpected(test *Test, fullPath string, result objects.Object) {
	value := tr.normalizeClosurePointers(strings.Trim(result.Inspect(), "\n"))

	if value != test.expect {
		tr.printErrorDoesNotMatchExpectation(
			test, fullPath,
			"Test expectation does not match the evaluated result",
			value, test.expect,
			VirtualMachineEngine,
		)
		return
	}

	tr.printSuccessStatusMessage(test, VirtualMachineEngine)
}

func (tr *TestRunner) compareCompliedVMWithStandardOutput(test *Test, fullPath string) {
	messages := vm.Stdout.ReadAll()
	if len(messages) == 0 {
		tr.printErrorStatusMessage(test, fullPath, "No output captured from standard output", VirtualMachineEngine)
		return
	}

	comparison := test.expect
	if comparison == "" {
		comparison = test.errors
	}

	out := tr.normalizeClosurePointers(strings.Trim(strings.Join(messages, ""), "\n"))
	if out != comparison {
		tr.printErrorDoesNotMatchExpectation(
			test, fullPath,
			"Test expectation does not match the standard output",
			out, comparison,
			VirtualMachineEngine,
		)
		return
	}

	tr.printSuccessStatusMessage(test, VirtualMachineEngine)
}
