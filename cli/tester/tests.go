package tester

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/senither/zen-lang/cli/colors"
	"github.com/senither/zen-lang/lexer"
	"github.com/senither/zen-lang/parser"
)

type TestRunner struct {
	directory string
	path      string
	options   RunnerOptions
	timings   map[RunnerTimings]time.Duration

	passedTests int
}

type RunnerOptions struct {
	Engine  EngineType
	Filter  string
	Verbose bool
	Compact bool
}

type EngineType int

const (
	AllEngines EngineType = iota
	EvaluatorEngine
	VirtualMachineEngine
)

func (et EngineType) GetTag() string {
	switch et {
	case AllEngines:
		return colors.Green + "All"
	case EvaluatorEngine:
		return colors.Cyan + "Eval"
	case VirtualMachineEngine:
		return colors.Magenta + "VM"

	default:
		return "unknown"
	}
}

type RunnerTimings string

const (
	FileDiscoveryTiming      RunnerTimings = "File Discovery"
	ReadingFilesTiming       RunnerTimings = "Reading Files"
	LexingAndParsingTiming   RunnerTimings = "Lexing + Parsing"
	CompilationTiming        RunnerTimings = "Compilation"
	EvaluatorExecutionTiming RunnerTimings = "Evaluator Execution"
	VMExecutionTiming        RunnerTimings = "VM Execution"
)

type Test struct {
	message         string
	file            string
	expect          string
	errors          string
	supportedEngine EngineType

	metadata map[any]any
}

var (
	PASSED          = colors.BgGreen + colors.Gray + " PASS " + colors.Reset
	FAILED          = colors.BgRed + colors.Gray + " FAIL " + colors.Reset
	messages        = []string{}
	collectedErrors = []string{}
	exitStatusCode  = 0
)

func NewTestRunner(directory, path string, options RunnerOptions) *TestRunner {
	return &TestRunner{
		directory: directory,
		path:      path,
		options:   options,
		timings:   make(map[RunnerTimings]time.Duration),
	}
}

func (tr *TestRunner) addTiming(timingType RunnerTimings, duration time.Duration) {
	tr.timings[timingType] = tr.getTiming(timingType) + duration
}

func (tr *TestRunner) getTiming(timingType RunnerTimings) time.Duration {
	timing := tr.timings[timingType]
	if timing <= 0 {
		return 0
	}

	return timing
}

func (tr *TestRunner) RunTests() error {
	groupedFiles, err := tr.discoverAndGroupTestFiles()
	if err != nil {
		return err
	}

	var groups []string
	for dir := range groupedFiles {
		groups = append(groups, dir)
	}

	sort.Strings(groups)

	for _, dir := range groups {
		messages = []string{}

		fullPath := fmt.Sprintf("%s%s%s", tr.path, string(os.PathSeparator), dir[len(tr.path)+1:])
		errorsCount := len(collectedErrors)

		for _, file := range groupedFiles[dir] {
			tr.runTestFile(fullPath, file)
		}

		if !tr.options.Compact {
			if errorsCount != len(collectedErrors) {
				fmt.Printf("  %s %s\n", FAILED, fullPath)
			} else if len(messages) > 0 {
				fmt.Printf("  %s %s\n", PASSED, fullPath)
			}

			if len(messages) > 0 || errorsCount != len(collectedErrors) {
				fmt.Println(strings.Join(messages, ""))
			}
		}
	}

	if tr.options.Compact {
		fmt.Print("\n\n")
	}

	if len(collectedErrors) > 0 {
		fmt.Println("Test suite failed with the following errors:")

		for _, err := range collectedErrors {
			parts := strings.Split(err, "\n")

			fmt.Println()
			fmt.Printf(" %s- %s%s\n%s\n", colors.Red, parts[0], colors.Reset, strings.Join(parts[1:], "\n"))
		}

		fmt.Println()
	}

	tr.printFinishedTestSuiteSummary()

	os.Exit(exitStatusCode)
	return nil
}

func (tr *TestRunner) discoverAndGroupTestFiles() (map[string][]string, error) {
	start := time.Now()
	var relativeTestFiles []string

	err := filepath.Walk(tr.path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".zent") {
			absolutePath, _ := filepath.Abs(path)
			relativeTestFiles = append(relativeTestFiles, absolutePath)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	groupedTestFiles := make(map[string][]string)

	for _, relativePath := range relativeTestFiles {
		dir := filepath.Dir(relativePath)
		groupedTestFiles[dir] = append(groupedTestFiles[dir], relativePath)
	}

	tr.addTiming(FileDiscoveryTiming, time.Since(start))

	return groupedTestFiles, nil
}

func (tr *TestRunner) runTestFile(fullPath, file string) {
	test, err := tr.readTestFileFromDisk(file)
	if err != nil {
		collectedErrors = append(collectedErrors, fmt.Sprintf("Error parsing test file %s: %s", file, err.Error()))
		return
	}

	if tr.options.Filter != "" && !strings.Contains(test.message, tr.options.Filter) {
		return
	}

	startLexingAndParsing := time.Now()
	l := lexer.New(test.file)
	p := parser.New(l, file)

	program := p.ParseProgram()
	tr.addTiming(LexingAndParsingTiming, time.Since(startLexingAndParsing))

	if len(p.Errors()) > 0 {
		msg := []string{"Parser errors found"}
		for _, err := range p.Errors() {
			msg = append(msg, fmt.Sprintf("     %s", err.String()))
		}

		tr.printErrorStatusMessage(test, fullPath, strings.Join(msg, "\n"), AllEngines)
		return
	}

	if tr.shouldRunTest(test, EvaluatorEngine) {
		tr.runEvaluatorTest(test, program, fullPath, file)
	}

	if tr.shouldRunTest(test, VirtualMachineEngine) {
		tr.runVMTest(test, program, fullPath, file)
	}
}

func (tr *TestRunner) readTestFileFromDisk(file string) (*Test, error) {
	start := time.Now()

	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	key := ""
	test := &Test{
		message:         "",
		file:            "",
		expect:          "",
		errors:          "",
		supportedEngine: AllEngines,

		metadata: make(map[any]any),
	}

	for line := range strings.SplitSeq(string(content), "\n") {
		if strings.HasPrefix(line, "--TEST--") || strings.HasPrefix(line, "---TEST---") {
			key = "message"
		} else if strings.HasPrefix(line, "--FILE--") || strings.HasPrefix(line, "---FILE---") {
			key = "file"
		} else if strings.HasPrefix(line, "--EXPECT--") || strings.HasPrefix(line, "---EXPECT---") {
			key = "expect"
		} else if strings.HasPrefix(line, "--ERROR--") || strings.HasPrefix(line, "---ERROR---") {
			key = "errors"
		}

		if strings.HasPrefix(line, "--") {
			continue
		}

		switch key {
		case "message":
			test.message += line + "\n"
		case "file":
			test.file += line + "\n"
		case "expect":
			test.expect += line + "\n"
		case "errors":
			test.errors += line + "\n"
		}
	}

	test.file = tr.normalizeLineEndings(test.file)
	test.expect = tr.normalizeLineEndings(test.expect)
	test.errors = tr.normalizeLineEndings(test.errors)

	if strings.HasSuffix(file, ".vm.zent") {
		test.supportedEngine = VirtualMachineEngine
	} else if strings.HasSuffix(file, ".eval.zent") {
		test.supportedEngine = EvaluatorEngine
	}

	tr.addTiming(ReadingFilesTiming, time.Since(start))

	return test, nil
}

func (tr *TestRunner) getStatusSummary() string {
	var parts []string

	if len(collectedErrors) > 0 {
		parts = append(parts, fmt.Sprintf("%s%d failed%s", colors.Red, len(collectedErrors), colors.Reset))
	}

	parts = append(parts, fmt.Sprintf("%s%d passed%s", colors.Green, tr.passedTests, colors.Reset))

	return strings.Join(parts, ", ")
}

func (tr *TestRunner) printFinishedTestSuiteSummary() {
	var totalTimeTake time.Duration
	for _, timing := range tr.timings {
		totalTimeTake += timing
	}

	fmt.Printf("  Finished running the test suite in %s\n", tr.directory)

	if !tr.options.Verbose {
		fmt.Println()
		fmt.Printf("  Tests:    %s\n", tr.getStatusSummary())
		fmt.Printf("  Duration: %s\n\n", totalTimeTake)

		return
	}

	fmt.Printf("  Tests: %s\n\n", tr.getStatusSummary())
	fmt.Printf("     Tests discovery: %s\n", tr.getTiming(FileDiscoveryTiming))
	fmt.Printf("       Reading files: %s\n", tr.getTiming(ReadingFilesTiming))
	fmt.Printf("      Lexer + Parser: %s\n", tr.getTiming(LexingAndParsingTiming))

	if tr.options.Engine == AllEngines || tr.options.Engine == EvaluatorEngine {
		fmt.Printf(" -----------------------------------\n")
		fmt.Printf("          Evaluation: %s\n", tr.getTiming(EvaluatorExecutionTiming))
	}

	if tr.options.Engine == AllEngines || tr.options.Engine == VirtualMachineEngine {
		fmt.Printf(" -----------------------------------\n")
		fmt.Printf("  Compile + Optimize: %s\n", tr.getTiming(CompilationTiming))
		fmt.Printf("          VM Runtime: %s\n", tr.getTiming(VMExecutionTiming))
	}

	fmt.Printf(" -----------------------------------\n")
	fmt.Printf("               Total: %s\n", totalTimeTake)
	fmt.Printf("\n")
}

func (tr *TestRunner) shouldRunTest(test *Test, engine EngineType) bool {
	if tr.options.Engine != AllEngines {
		if tr.options.Engine != engine {
			return false
		}

		if test.supportedEngine != AllEngines && test.supportedEngine != tr.options.Engine {
			return false
		}

		return true
	}

	if test.supportedEngine != AllEngines && test.supportedEngine != engine {
		return false
	}

	return true
}
