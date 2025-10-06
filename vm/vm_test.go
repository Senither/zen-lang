package vm

import (
	"testing"

	"github.com/senither/zen-lang/objects"
)

type vmTestCase struct {
	input    string
	expected interface{}
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"1.5 + 2.75", 4.25},
		{"1 - 2", -1},
		{"2 * 2", 4},
		{"2 * 2.5", 5.0},
		{"4 / 2", 2},
		{"5 / 2", 2.5},
		{"50 / 2 * 2 + 10 - 5", 55},
		{"49 / 2 * 3 + 10 - 5", 78.5},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"4.25 + 4.25 + 4.25 + 4.25 - 10", 7.0},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"2.5 * 2.5 * 2.5 * 2.5 * 2.5", 97.65625},
		{"5 * 2.125 + 10", 20.625},
		{"5 + 2 * 10", 25},
		{"5 + 2.125 * 10", 26.25},
		{"5 * (2 + 10)", 60},
		{"5 * (2.125 + 10)", 60.625},
		{"-5", -5},
		{"-5.5", -5.5},
		{"-10 + 5", -5},
		{"-10.5 + 5.5", -5.0},
		{"-(5 + 5)", -10},
		{"-(5.5 + 4.5)", -10.0},
		{"-50 + 100 + -50", 0},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
		{"5.5 ^ 0", 1.0},
		{"5.5 ^ 2", 30.25},
		{"2.5 + 5.5 ^ 2", 32.75},
		{"(2.5 + 5.5) ^ 2 * 2", 128.0},
		{"10.75 % 3", 1.75},
		{"12.34 % 5", 2.34},
		{"3.14 % 2", 1.14},
	}

	runVmTests(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"1 <= 2", true},
		{"1 >= 2", false},
		{"1 <= 1", true},
		{"1 >= 1", true},
		{"2 <= 1", false},
		{"2 >= 1", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{"!(if (false) { 5; })", true},
	}

	runVmTests(t, tests)
}

func TestStringExpressions(t *testing.T) {
	tests := []vmTestCase{
		{`"hello world"`, "hello world"},
		{`"hello" + " " + "world"`, "hello world"},
		{`"foo" + "bar"`, "foobar"},
	}

	runVmTests(t, tests)
}

func TestArrayLiterals(t *testing.T) {
	tests := []vmTestCase{
		{"[]", []int{}},
		{"[1, 2, 3]", []int{1, 2, 3}},
		{"[1 + 2, 3 * 4, 5 + 6]", []int{3, 12, 11}},
	}

	runVmTests(t, tests)
}

func TestHashLiterals(t *testing.T) {
	tests := []vmTestCase{
		{
			"{}", map[objects.HashKey]int64{},
		},
		{
			"{1: 2, 3: 4, 5: 6}",
			map[objects.HashKey]int64{
				(&objects.Integer{Value: 1}).HashKey(): 2,
				(&objects.Integer{Value: 3}).HashKey(): 4,
				(&objects.Integer{Value: 5}).HashKey(): 6,
			},
		},
		{
			"{1 + 1: 2 * 2, 3 + 3: 4 * 4}",
			map[objects.HashKey]int64{
				(&objects.Integer{Value: 2}).HashKey(): 4,
				(&objects.Integer{Value: 6}).HashKey(): 16,
			},
		},
	}

	runVmTests(t, tests)
}

func TestIndexExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][0 + 2]", 3},
		{"[[1, 1, 1]][0][0]", 1},
		{"[][0]", nil},
		{"[1, 2, 3][99]", nil},
		{"[1][-1]", nil},
		{"{1: 1, 2: 2}[1]", 1},
		{"{1: 1, 2: 2}[2]", 2},
		{"{1: 1}[0]", nil},
		{"{}[0]", nil},
	}

	runVmTests(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []vmTestCase{
		{"if (true) { 10 }", 10},
		{"if (true) { 10 } else { 20 }", 10},
		{"if (false) { 10 } else { 20 }", 20},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (false) { 10 }", nil},
		{"if (1 > 2) { 10 }", nil},
		{"if ((if (false) { 10 })) { 10 } else { 20 }", 20},
		{"if (true) { 10 } else if (false) { 20 } else { 30 }", 10},
		{"if (false) { 10 } else if (true) { 20 } else { 30 }", 20},
		{"if (false) { 10 } else if (false) { 20 } else { 30 }", 30},
		{"if (false) { 10 } else if (if (true) { 15}) { 20 } else { 30 }", 20},
	}

	runVmTests(t, tests)
}

func TestGlobalVarStatements(t *testing.T) {
	tests := []vmTestCase{
		{"var a = 1; a;", 1},
		{"var a = 1; var b = 2; a + b;", 3},
		{"var a = 1; var b = a + 1; a + b;", 3},
		{"var a = 1; var b = a + a; a + b;", 3},
	}

	runVmTests(t, tests)
}
