package tester

import (
	"strings"
	"time"

	"github.com/senither/zen-lang/ast"
	"github.com/senither/zen-lang/evaluator"
	"github.com/senither/zen-lang/objects"
)

func (tr *TestRunner) runEvaluatorTest(test *Test, program *ast.Program, fullPath, file string) {
	tr.incrementTestsFound(EvaluatorEngine)
	tr.applyTestEnvVariables(test)

	start := time.Now()
	evaluated := evaluator.Stdout.Mute(func() objects.Object {
		env := objects.NewEnvironment(file)
		return evaluator.Eval(program, env)
	})
	timeTaken := time.Since(start)

	objects.RestoreObjectsState()
	tr.clearTestEnvVariables()

	tr.addTiming(EvaluatorExecutionTiming, timeTaken)
	test.metadata[EvaluatorExecutionTiming] = timeTaken

	if evaluated == nil {
		tr.printErrorStatusMessage(test, fullPath, "Evaluator returned nil, failed to evaluate the test input", EvaluatorEngine)
		return
	}

	if objects.IsError(evaluated) {
		tr.compareEvaluatedWithError(test, fullPath, evaluated)
	} else if evaluated.Type() != objects.NULL_OBJ {
		tr.compareEvaluatedWithExpected(test, fullPath, evaluated)
	} else {
		tr.compareEvaluatedWithStandardOutput(test, fullPath)
	}
}

func (tr *TestRunner) compareEvaluatedWithError(test *Test, fullPath string, result objects.Object) {
	err := tr.normalizeFileLocations(result.Inspect())

	if err != test.errors {
		var message = "Test expectation does not match the evaluated result"
		if len(test.errors) == 0 {
			message = "No error expectation were provided, despite the result being *objects.Error"
		}

		tr.printErrorDoesNotMatchExpectation(
			test, fullPath,
			message, err, test.errors,
			EvaluatorEngine,
		)
		return
	}

	tr.printSuccessStatusMessage(test, EvaluatorEngine)
}

func (tr *TestRunner) compareEvaluatedWithExpected(test *Test, fullPath string, evaluated objects.Object) {
	if strings.Trim(evaluated.Inspect(), "\n") != test.expect {
		tr.printErrorDoesNotMatchExpectation(
			test, fullPath,
			"Test expectation does not match the evaluated result",
			strings.Trim(evaluated.Inspect(), "\n"), test.expect,
			EvaluatorEngine,
		)
		return
	}

	tr.printSuccessStatusMessage(test, EvaluatorEngine)
}

func (tr *TestRunner) compareEvaluatedWithStandardOutput(test *Test, fullPath string) {
	messages := evaluator.Stdout.ReadAll()
	if len(messages) == 0 {
		tr.printErrorStatusMessage(test, fullPath, "No output captured from standard output", EvaluatorEngine)
		return
	}

	out := strings.Trim(strings.Join(messages, ""), "\n")
	if out != test.expect {
		tr.printErrorDoesNotMatchExpectation(
			test, fullPath,
			"Test expectation does not match the standard output",
			out, test.expect,
			EvaluatorEngine,
		)
		return
	}

	tr.printSuccessStatusMessage(test, EvaluatorEngine)
}
