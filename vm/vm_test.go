package vm

import "testing"

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
	}

	runVmTests(t, tests)
}
