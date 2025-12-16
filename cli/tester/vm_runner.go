package tester

import (
	"fmt"
	"strings"
	"time"

	"github.com/senither/zen-lang/ast"
	"github.com/senither/zen-lang/compiler"
	"github.com/senither/zen-lang/objects"
	"github.com/senither/zen-lang/optimizer"
	"github.com/senither/zen-lang/vm"
)

func (tr *TestRunner) runVMTest(test *Test, program *ast.Program, fullPath, file string) {
	if !tr.options.Verbose {
		defer func() {
			if r := recover(); r != nil {
				message := "Virtual machine panicked during execution:\n\t%s"

				switch e := r.(type) {
				case error:
					message = fmt.Sprintf(message, e.Error())

				default:
					message = fmt.Sprintf(message, r)
				}

				tr.printErrorStatusMessage(test, fullPath, message, VirtualMachineEngine)
			}
		}()
	}

	tr.incrementTestsFound(VirtualMachineEngine)

	start := time.Now()

	c := compiler.New(file)
	compilerErr := c.Compile(program)
	timeTaken := time.Since(start)

	tr.addTiming(CompilationTiming, timeTaken)
	test.metadata[CompilationTiming] = timeTaken

	if compilerErr != nil {
		tr.compareCompliedVMWithError(test, fullPath, compilerErr.Error(), VirtualMachineEngine)
		return
	}

	tr.runCompiledVMTest(test, c.Bytecode(), fullPath, _VirtualMachineEngineUnprocessed)

	start = time.Now()
	bytes := c.Bytecode().Serialize()
	timeTaken = time.Since(start)

	tr.addTiming(SerializationTiming, timeTaken)
	test.metadata[SerializationTiming] = timeTaken

	start = time.Now()
	deserializedBytecode, err := compiler.Deserialize(bytes)
	timeTaken = time.Since(start)

	tr.addTiming(DeserializationTiming, timeTaken)
	test.metadata[DeserializationTiming] = timeTaken

	if err != nil {
		tr.printErrorStatusMessage(
			test, fullPath,
			"Failed to deserialize bytecode: "+err.Error(),
			_VirtualMachineEngineSerialized,
		)
		return
	}

	tr.runCompiledVMTest(test, deserializedBytecode, fullPath, _VirtualMachineEngineSerialized)

	start = time.Now()
	optimized, err := optimizer.Optimize(deserializedBytecode)
	timeTaken = time.Since(start)

	tr.addTiming(OptimizationTiming, timeTaken)
	test.metadata[OptimizationTiming] = timeTaken

	if err != nil {
		tr.printErrorStatusMessage(
			test, fullPath,
			"Failed to optimize bytecode: "+err.Error(),
			_VirtualMachineEngineOptimized,
		)
		return
	}

	test.metadata[MetaOriginalSize] = deserializedBytecode.OperationsCount()
	test.metadata[MetaOptimizationSize] = optimized.OperationsCount()

	tr.runCompiledVMTest(test, optimized, fullPath, _VirtualMachineEngineOptimized)
}

func (tr *TestRunner) runCompiledVMTest(test *Test, bytecode *compiler.Bytecode, fullPath string, engineType EngineType) {
	tr.applyTestEnvVariables(test)

	start := time.Now()
	result := vm.Stdout.Mute(func() objects.Object {
		runner := vm.New(bytecode)
		runner.EnableStdoutCapture()

		if runErr := runner.Run(); runErr != nil {
			return objects.NativeErrorToErrorObject(runErr)
		}

		return runner.LastPoppedStackElem()
	})
	timeTaken := time.Since(start)

	objects.RestoreObjectsState()
	tr.clearTestEnvVariables()

	tr.addTiming(VMExecutionTiming, timeTaken)
	test.metadata[VMExecutionTiming] = timeTaken

	if objects.IsError(result) {
		tr.compareCompliedVMWithError(test, fullPath, result.Inspect(), engineType)
	} else if result != nil && len(vm.Stdout.ReadAll()) == 0 {
		tr.compareCompliedVMWithExpected(test, fullPath, result, engineType)
	} else {
		tr.compareCompliedVMWithStandardOutput(test, fullPath, engineType)
	}
}

func (tr *TestRunner) compareCompliedVMWithError(
	test *Test,
	fullPath string,
	errorMessage string,
	engineType EngineType,
) {
	err := tr.normalizeFileLocations(errorMessage)

	if err != test.errors {
		tr.printErrorDoesNotMatchExpectation(
			test, fullPath,
			"Test expectation does not match the compiler error",
			err, test.errors,
			engineType,
		)
		return
	}

	tr.printSuccessStatusMessage(test, engineType)
}

func (tr *TestRunner) compareCompliedVMWithExpected(
	test *Test,
	fullPath string,
	result objects.Object,
	engineType EngineType,
) {
	value := tr.normalizeClosurePointers(strings.Trim(result.Inspect(), "\n"))

	if value != test.expect {
		tr.printErrorDoesNotMatchExpectation(
			test, fullPath,
			"Test expectation does not match the evaluated result",
			value, test.expect,
			engineType,
		)
		return
	}

	tr.printSuccessStatusMessage(test, engineType)
}

func (tr *TestRunner) compareCompliedVMWithStandardOutput(
	test *Test,
	fullPath string,
	engineType EngineType,
) {
	messages := vm.Stdout.ReadAll()
	if len(messages) == 0 {
		tr.printErrorStatusMessage(test, fullPath, "No output captured from standard output", engineType)
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
			engineType,
		)
		return
	}

	tr.printSuccessStatusMessage(test, engineType)
}
