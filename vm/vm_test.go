package vm

import (
	"testing"

	"github.com/senither/zen-lang/objects"
	"github.com/senither/zen-lang/objects/timer"
)

type vmTestCase struct {
	name     any
	input    string
	expected any
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{nil, "1", 1},
		{nil, "2", 2},
		{nil, "1 + 2", 3},
		{nil, "1.5 + 2.75", 4.25},
		{nil, "1 - 2", -1},
		{nil, "2 * 2", 4},
		{nil, "2 * 2.5", 5.0},
		{nil, "4 / 2", 2},
		{nil, "5 / 2", 2.5},
		{nil, "50 / 2 * 2 + 10 - 5", 55},
		{nil, "49 / 2 * 3 + 10 - 5", 78.5},
		{nil, "5 + 5 + 5 + 5 - 10", 10},
		{nil, "4.25 + 4.25 + 4.25 + 4.25 - 10", 7.0},
		{nil, "2 * 2 * 2 * 2 * 2", 32},
		{nil, "2.5 * 2.5 * 2.5 * 2.5 * 2.5", 97.65625},
		{nil, "5 * 2.125 + 10", 20.625},
		{nil, "5 + 2 * 10", 25},
		{nil, "5 + 2.125 * 10", 26.25},
		{nil, "5 * (2 + 10)", 60},
		{nil, "5 * (2.125 + 10)", 60.625},
		{nil, "-5", -5},
		{nil, "-5.5", -5.5},
		{nil, "-10 + 5", -5},
		{nil, "-10.5 + 5.5", -5.0},
		{nil, "-(5 + 5)", -10},
		{nil, "-(5.5 + 4.5)", -10.0},
		{nil, "-50 + 100 + -50", 0},
		{nil, "(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
		{nil, "5.5 ^ 0", 1.0},
		{nil, "5.5 ^ 2", 30.25},
		{nil, "2.5 + 5.5 ^ 2", 32.75},
		{nil, "(2.5 + 5.5) ^ 2 * 2", 128.0},
		{nil, "10.75 % 3", 1.75},
		{nil, "12.34 % 5", 2.34},
		{nil, "3.14 % 2", 1.14},
		{nil, "2 * 3 ^ 4", 162},
		{nil, "2 * 3 ^ 4 % 5", 2},
	}

	runVmTests(t, tests)
}

func BenchmarkIntegerArithmeticSimple(b *testing.B) {
	runVmBenchmark(b, `5 + 10 - 3 * 2 / 4 + 6 ^ 2 % 4`)
}

func BenchmarkIntegerArithmeticComplex(b *testing.B) {
	runVmBenchmark(b, `(5 + 10 - 3 * 2 / 4 + 6 ^ 2 % 4) * (12 - 4 + 3 * 7 / 2 ^ 3 % 5) + (8 ^ 2 % 5 + 14 - 3 * 6 / 2 + 9)`)
}

func TestBooleanExpressions(t *testing.T) {
	tests := []vmTestCase{
		{nil, "true", true},
		{nil, "false", false},
		{"int truthy lt", "1 < 2", true},
		{"int falsy gt", "1 > 2", false},
		{"int falsy lt", "1 < 1", false},
		{"int falsy gt", "1 > 1", false},
		{"int truthy eq", "1 == 1", true},
		{"int falsy neq", "1 != 1", false},
		{"int falsy eq", "1 == 2", false},
		{"int truthy neq", "1 != 2", true},
		{"int truthy lte", "1 <= 2", true},
		{"int falsy gte", "1 >= 2", false},
		{"int truthy lte", "1 <= 1", true},
		{"int truthy gte", "1 >= 1", true},
		{"int falsy lte", "2 <= 1", false},
		{"int truthy gte", "2 >= 1", true},
		{"bool truthy eq", "true == true", true},
		{"bool truthy eq", "false == false", true},
		{"bool falsy eq", "true == false", false},
		{"bool truthy neq", "true != false", true},
		{"bool truthy neq", "false != true", true},
		{"bool truthy eq", "(1 < 2) == true", true},
		{"bool falsy eq", "(1 < 2) == false", false},
		{"bool falsy eq", "(1 > 2) == true", false},
		{"bool truthy eq", "(1 > 2) == false", true},
		{"bool falsy not", "!true", false},
		{"bool truthy not", "!false", true},
		{"bool falsy not", "!5", false},
		{"bool truthy not", "!!true", true},
		{"bool falsy not", "!!false", false},
		{"bool truthy not", "!!5", true},
		{"bool truthy not", "!(if (false) { 5; })", true},
		{"bool truthy and", "true && true", true},
		{"bool falsy and", "true && false", false},
		{"bool falsy and", "false && true", false},
		{"bool truthy or", "true || true", true},
		{"bool truthy or", "true || false", true},
		{"bool truthy or", "false || true", true},
		{"bool falsy or", "false || false", false},
		{"bool truthy and", "true && (false || true)", true},
		{"bool truthy and", "(1 < 2) && (2 < 3)", true},
		{"bool falsy and", "(1 < 2) && (2 > 3)", false},
		{"bool truthy or", "(1 > 2) || (2 < 3)", true},
		{"bool falsy or", "(1 > 2) || (2 > 3)", false},
	}

	runVmTests(t, tests)
}

func BenchmarkBooleanExpressions(b *testing.B) {
	runVmBenchmark(b, `(1 < 2) && (2 < 3) || (3 < 4) && !(4 > 5)`)
}

func TestStringExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"string", `"hello world"`, "hello world"},
		{"string concatenation", `"hello" + " " + "world"`, "hello world"},
		{"string concatenation", `"foo" + "bar"`, "foobar"},
		{"string concatenation", `"foo" + " " + "bar"`, "foo bar"},
		{"string concatenation with number", `"The answer is: " + 42`, "The answer is: 42"},
		{"string concatenation with number", `"Pi is approximately " + 3.14`, "Pi is approximately 3.14"},
		{"string concatenation with boolean", `"Value: " + true`, "Value: true"},
		{"string concatenation with boolean", `"Value: " + false`, "Value: false"},
		{"string concatenation with number", `"Number: " + (10 + 5)`, "Number: 15"},
	}

	runVmTests(t, tests)
}

func BenchmarkStringExpressions(b *testing.B) {
	runVmBenchmark(b, `"hello" + " " + "world"`)
}

func TestArrayLiterals(t *testing.T) {
	tests := []vmTestCase{
		{"empty array", "[]", []int{}},
		{"array with integers", "[1, 2, 3]", []int{1, 2, 3}},
		{"array with expressions", "[1 + 2, 3 * 4, 5 + 6]", []int{3, 12, 11}},
	}

	runVmTests(t, tests)
}

func BenchmarkArrayLiterals(b *testing.B) {
	runVmBenchmark(b, "[1 + 2, 3 * 4, 5 + 6]")
}

func TestHashLiterals(t *testing.T) {
	tests := []vmTestCase{
		{
			"empty hash",
			"{}", map[objects.HashKey]int64{},
		},
		{
			"hash with integers",
			"{1: 2, 3: 4, 5: 6}",
			map[objects.HashKey]int64{
				(&objects.Integer{Value: 1}).HashKey(): 2,
				(&objects.Integer{Value: 3}).HashKey(): 4,
				(&objects.Integer{Value: 5}).HashKey(): 6,
			},
		},
		{
			"hash with expressions",
			"{1 + 1: 2 * 2, 3 + 3: 4 * 4}",
			map[objects.HashKey]int64{
				(&objects.Integer{Value: 2}).HashKey(): 4,
				(&objects.Integer{Value: 6}).HashKey(): 16,
			},
		},
	}

	runVmTests(t, tests)
}

func BenchmarkHashLiterals(b *testing.B) {
	runVmBenchmark(b, "{1 + 1: 2 * 2, 3 + 3: 4 * 4}")
}

func TestSuffixExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"suffix increment once", "var x = 5; x++", 6},
		{"suffix decrement once", "var x = 5; x--", 4},
		{"suffix increment and get value", "var x = 5; x++; x", 6},
		{"suffix decrement and get value", "var x = 5; x--; x", 4},
		{"suffix increment twice", "var x = 5; x++; x++; x", 7},
		{"suffix decrement twice", "var x = 5; x--; x--; x", 3},
		{"suffix increment and decrement", "var x = 5; x++; x--; x", 5},
	}

	runVmTests(t, tests)
}

func BenchmarkSuffixExpressions(b *testing.B) {
	runVmBenchmark(b, "var x = 5; x++; x--; x")
}

func TestIndexExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"array index of 1", "[1, 2, 3][1]", 2},
		{"array index expression", "[1, 2, 3][0 + 2]", 3},
		{"array index of index", "[[1, 1, 1]][0][0]", 1},
		{"array index of 0 out of bounds", "[][0]", nil},
		{"array index of 99 out of bounds", "[1, 2, 3][99]", nil},
		{"array index negative", "[1][-1]", 1},
		{"array index negative out of bounds", "[1][-2]", nil},
		{"hash index of 1", "{1: 1, 2: 2}[1]", 1},
		{"hash index of 2", "{1: 1, 2: 2}[2]", 2},
		{"hash index not exists", "{1: 1}[0]", nil},
		{"hash index of 0 not exists", "{}[0]", nil},
	}

	runVmTests(t, tests)
}

func BenchmarkIndexExpressions(b *testing.B) {
	runVmBenchmark(b, "[1, 2, 3][1]")
}

func TestChainIndexExpressions(t *testing.T) {
	tests := []vmTestCase{
		{
			name: "hash chain indexing",
			input: `
				var obj = {'a': 2};
				obj.a
			`,
			expected: 2,
		},
		{
			name: "nested hash chain indexing",
			input: `
				var obj = {'a': {'b': 3}};
				obj.a.b
			`,
			expected: 3,
		},
		{
			name: "array inside hash chain indexing",
			input: `
				var obj = {'a': [4]};
				obj.a[0]
			`,
			expected: 4,
		},
		{
			name: "hash inside array chain indexing",
			input: `
				var obj = {'a': {'b': func() { 5 }}};
				obj.a.b()
			`,
			expected: 5,
		},
		{
			name: "deeply nested chain indexing function call",
			input: `
				var obj = {
					'a': {
						'b': func(a, b) {
							return a + b
						}
					}
				}

				obj.a.b(2, 4)
			`,
			expected: 6,
		},
	}

	runVmTests(t, tests)
}

func BenchmarkChainIndexExpressions(b *testing.B) {
	runVmBenchmark(b, `
		var obj = {'a': {'b': 3}};
		obj.a.b
	`)
}

func TestChainIndexAssignment(t *testing.T) {
	tests := []vmTestCase{
		{
			name: "hash index assignment with block notation",
			input: `
				var obj = {'a': 1};
				obj['a'] = 10;
				obj['a']
			`,
			expected: 10,
		},
		{
			name: "hash index assignment with dot notation",
			input: `
				var obj = {'a': 1};
				obj.a = 10;
				obj.a
			`,
			expected: 10,
		},
		{
			name: "hash index assignment with new key",
			input: `
				var obj = {'a': 1};
				obj['b'] = 10;
				obj['b']
			`,
			expected: 10,
		},
		{
			name: "nested hash index assignment with block notation",
			input: `
				var obj = {'nested': {'key': 'value'}};
				obj['nested']['key'] = 42;
				obj['nested']['key']
			`,
			expected: 42,
		},
		{
			name: "nested hash index assignment with new key",
			input: `
				var obj = {'nested': {'key': 'value'}};
				obj['nested']['newKey'] = 42;
				obj['nested']['newKey']
			`,
			expected: 42,
		},
	}

	runVmTests(t, tests)
}

func BenchmarkChainIndexAssignment(b *testing.B) {
	runVmBenchmark(b, `
		var obj = {'nested': {'key': 'value'}};
		obj['nested']['key'] = 42;
		obj['nested']['key']
	`)
}

func TestChainedHashAssignmentExpressions(t *testing.T) {
	tests := []vmTestCase{
		{
			"hash index assignment with dot notation",
			"var x = {'foo': 5}; x.foo = 10; x.foo;",
			10,
		},
		{
			"hash index assignment with block notation",
			"var x = {'foo': 5}; x['foo'] = 10; x.foo;",
			10,
		},
		{
			"nested hash index assignment with dot notation",
			"var x = {'foo': {'bar': 5}}; x.foo.bar = 10; x.foo.bar;",
			10,
		},
		{
			"nested hash index assignment with block notation",
			"var x = {'foo': {'bar': 5}}; x['foo']['bar'] = 10; x.foo.bar;",
			10,
		},
		{
			"hash index assignment with dot notation",
			"var x = {'foo': 5}; var y = {'bar': 10}; x.foo = y.bar; x.foo;",
			10,
		},
		{
			"hash index assignment with dot notation",
			"var x = {'foo': 5}; var y = {'bar': 10}; x.newKey = y.bar; x.newKey;",
			10,
		},
	}

	runVmTests(t, tests)
}

func BenchmarkChainedHashAssignmentExpressions(b *testing.B) {
	runVmBenchmark(b, `
		var obj = {'nested': {'key': 'value'}};
		obj.nested.key = 10;
		obj.nested.key;
	`)
}

func TestChainedArrayIndexAssignments(t *testing.T) {
	tests := []vmTestCase{
		{
			name: "simple array index assignment",
			input: `
				var obj = {'arr': [1, 2, 3]};
				obj.arr[1] = 42;
				obj.arr
			`,
			expected: []int{1, 42, 3},
		},
		{
			name: "array index assignment with expression",
			input: `
				var obj = {'arr': [1, 2, 3, 4]};
				obj.arr[2] = obj.arr[2] + 40;
				obj.arr
			`,
			expected: []int{1, 2, 43, 4},
		},
		{
			name: "chained array index assignment",
			input: `
				var obj = {'foo': {'bar': [5, 6]}};
				obj.foo.bar[0] = 99;
				obj.foo.bar
			`,
			expected: []int{99, 6},
		},
	}

	runVmTests(t, tests)
}

func BenchmarkChainedArrayIndexAssignments(b *testing.B) {
	runVmBenchmark(b, `
		var obj = {'foo': {'bar': [5, 6]}};
		obj.foo.bar[0] = 99;
		obj.foo.bar
	`)
}

func TestConditionals(t *testing.T) {
	tests := []vmTestCase{
		{
			"get value from truthy block",
			"if (true) { 10 }",
			10,
		},
		{
			"get value from truthy block with else",
			"if (true) { 10 } else { 20 }",
			10,
		},
		{
			"get value from falsy block with else",
			"if (false) { 10 } else { 20 }",
			20,
		},
		{
			"get value from truthy block with integer condition",
			"if (1) { 10 }",
			10,
		},
		{
			"get value from truthy block with less than condition",
			"if (1 < 2) { 10 }",
			10,
		},
		{
			"get value from truthy block with less than condition and else",
			"if (1 < 2) { 10 } else { 20 }",
			10,
		},
		{
			"get value from falsy block with greater than condition and else",
			"if (1 > 2) { 10 } else { 20 }",
			20,
		},
		{
			"get value from falsy block",
			"if (false) { 10 }",
			nil,
		},
		{
			"get value from falsy block",
			"if (1 > 2) { 10 }",
			nil,
		},
		{
			"get value from falsy block with nested if",
			"if ((if (false) { 10 })) { 10 } else { 20 }",
			20,
		},
		{
			"get value from truthy block with else if",
			"if (true) { 10 } else if (false) { 20 } else { 30 }",
			10,
		},
		{
			"get value from else if block",
			"if (false) { 10 } else if (true) { 20 } else { 30 }",
			20,
		},
		{
			"get value from else block",
			"if (false) { 10 } else if (false) { 20 } else { 30 }",
			30,
		},
		{
			"get value from else if block with nested if",
			"if (false) { 10 } else if (if (true) { 15}) { 20 } else { 30 }",
			20,
		},
	}

	runVmTests(t, tests)
}

func BenchmarkConditionals(b *testing.B) {
	runVmBenchmark(b, "if (1 < 2) { 10 } else { 20 }")
}

func TestGlobalVarStatements(t *testing.T) {
	tests := []vmTestCase{
		{"variable declaration and usage", "var a = 1; a;", 1},
		{"variable declaration with multiple variables", "var a = 1; var b = 2; a + b;", 3},
		{"variable declaration with expressions", "var a = 1; var b = a + 1; a + b;", 3},
		{"variable declaration with self", "var a = 1; var b = a + a; a + b;", 3},
	}

	runVmTests(t, tests)
}

func BenchmarkGlobalVarStatements(b *testing.B) {
	runVmBenchmark(b, "var a = 1; var b = a + 1; a + b;")
}

func TestCallingFunctionsWithoutArguments(t *testing.T) {
	tests := []vmTestCase{
		{
			name: "function call without arguments",
			input: `
				var fivePlusTen = func() { 5 + 10; };
				fivePlusTen();
			`,
			expected: 15,
		},
		{
			name: "function call without arguments added together",
			input: `
				var one = func() { 1; };
				var two = func() { 2; };
				one() + two()
			`,
			expected: 3,
		},
		{
			name: "nested function calls without arguments",
			input: `
				var a = func() { 1 };
				var b = func() { a() + 1 };
				var c = func() { b() + 1 };
				c();
			`,
			expected: 3,
		},
	}

	runVmTests(t, tests)
}

func BenchmarkCallingFunctionsWithoutArguments(b *testing.B) {
	runVmBenchmark(b, `
		var fivePlusTen = func() { 5 + 10; };
		fivePlusTen();
	`)
}

func TestFunctionsWithReturnStatement(t *testing.T) {
	tests := []vmTestCase{
		{
			name: "function with early return",
			input: `
				var earlyExit = func() { return 99; 100; };
				earlyExit();
			`,
			expected: 99,
		},
		{
			name: "function with multiple return statements",
			input: `
				var earlyExit = func() { return 99; return 100; };
				earlyExit();
			`,
			expected: 99,
		},
	}

	runVmTests(t, tests)
}

func BenchmarkFunctionsWithReturnStatement(b *testing.B) {
	runVmBenchmark(b, `
		var earlyExit = func() { return 99; 100; };
		earlyExit();
	`)
}

func TestFunctionsWithoutReturnValue(t *testing.T) {
	tests := []vmTestCase{
		{
			name: "function with implicitly return value",
			input: `
				var noReturn = func() { };
				noReturn();
			`,
			expected: nil,
		},
		{
			name: "nested function with implicitly return value",
			input: `
				var noReturn = func() { };
				var noReturnTwo = func() { noReturn(); };
				noReturn();
				noReturnTwo();
			`,
			expected: nil,
		},
	}

	runVmTests(t, tests)
}

func BenchmarkFunctionsWithoutReturnValue(b *testing.B) {
	runVmBenchmark(b, `
		var noReturn = func() { };
		noReturn();
	`)
}

func TestCallingFunctionsWithBindings(t *testing.T) {
	tests := []vmTestCase{
		{
			name: "function with local binding",
			input: `
				var one = func() {
					var one = 1;
					one
				}
				one();
			`,
			expected: 1,
		},
		{
			name: "function with multiple local bindings",
			input: `
				var oneAndTwo = func() {
					var one = 1;
					var two = 2;
					one + two;
				}
				oneAndTwo();
			`,
			expected: 3,
		},
		{
			name: "functions with multiple local bindings",
			input: `
				var oneAndTwo = func() {
					var one = 1;
					var two = 2;
					one + two;
				}
				var threeAndFour = func() {
					var three = 3;
					var four = 4;
					three + four;
				}
				oneAndTwo() + threeAndFour();
			`,
			expected: 10,
		},
		{
			name: "functions with same local binding names",
			input: `
				var firstFoobar = func() {
					var foobar = 50;
					foobar;
				}
				var secondFoobar = func() {
					var foobar = 100;
					foobar;
				}
				firstFoobar() + secondFoobar();
			`,
			expected: 150,
		},
		{
			name: "functions with local and global bindings",
			input: `
				var globalSeed = 50;
				var minusOne = func() {
					var num = 1;
					globalSeed - num;
				}
				var minusTwo = func() {
					var num = 2;
					globalSeed - num;
				}
				minusOne() + minusTwo();
			`,
			expected: 97,
		},
	}

	runVmTests(t, tests)
}

func BenchmarkCallingFunctionsWithBindings(b *testing.B) {
	runVmBenchmark(b, `
		var oneAndTwo = func() {
			var one = 1;
			var two = 2;
			one + two;
		}
		oneAndTwo();
	`)
}

func TestCallingFunctionsWithArgumentsAndBindings(t *testing.T) {
	tests := []vmTestCase{
		{
			name: "function with one argument",
			input: `
				var value = func(a) { a; };
				value(5);
			`,
			expected: 5,
		},
		{
			name: "function with two arguments",
			input: `
				var sum = func(a, b) { a + b; };
				sum(1, 2);
			`,
			expected: 3,
		},
		{
			name: "function with local binding and arguments",
			input: `
				var sum = func(a, b) {
					var c = a + b;
					c;
				};
				sum(1, 2);
			`,
			expected: 3,
		},
		{
			name: "function with local binding and arguments added together",
			input: `
				var sum = func(a, b) {
					var c = a + b;
					c;
				};
				sum(1, 2) + sum(3, 4);
			`,
			expected: 10,
		},
		{
			name: "nested function calls without arguments but with arguments and bindings",
			input: `
				var sum = func(a, b) {
					var c = a + b;
					c;
				}
				var outer = func() {
					sum(1, 2) + sum(3, 4);
				}
				outer();
			`,
			expected: 10,
		},
		{
			name: "function with global variable and arguments and bindings",
			input: `
				var globalNum = 10;

				var sum = func(a, b) {
					var x = a + b;
					x + globalNum;
				}

				var outer = func() {
					sum(1, 2) + sum(3, 4) + globalNum;
				}

				outer() + globalNum;
			`,
			expected: 50,
		},
	}

	runVmTests(t, tests)
}

func BenchmarkCallingFunctionsWithArgumentsAndBindings(b *testing.B) {
	runVmBenchmark(b, `
		var sum = func(a, b) {
			var c = a + b;
			c;
		};
		sum(1, 2) + sum(3, 4);
	`)
}

func TestCallingFunctionsWithWrongArguments(t *testing.T) {
	tests := []vmTestCase{
		{
			name:     "function call without arguments",
			input:    `func() { 1; }(1);`,
			expected: "wrong number of arguments to `<anonymous>`: got 1, want 0",
		},
		{
			name:     "function call with one argument",
			input:    `func(a) { a; }();`,
			expected: "wrong number of arguments to `<anonymous>`: got 0, want 1",
		},
		{
			name:     "function call with two arguments",
			input:    `func(a, b) { a + b; }(1);`,
			expected: "wrong number of arguments to `<anonymous>`: got 1, want 2",
		},
	}

	for _, tt := range tests {
		compiler, err := compile(tt.input)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		vm := New(compiler.Bytecode())
		err = vm.Run()
		if err == nil {
			t.Fatalf("expected VM error but resulted in none")
		}

		if err.Error() != tt.expected {
			t.Fatalf("wrong VM error:\nwant:\n\t%q\ngot:\n\t%q", tt.expected, err)
		}
	}
}

func BenchmarkCallingFunctionsWithWrongArguments(b *testing.B) {
	runVmBenchmark(b, "func() { 1; }(1);")
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []vmTestCase{
		{nil, `len("")`, 0},
		{nil, `len("four")`, 4},
		{nil, `len("hello world")`, 11},
		{nil, `len([])`, 0},
		{nil, `len([1, 2, 3])`, 3},
		{nil, `len(1)`, &objects.Error{Message: "argument 1 to `len` has invalid type: got INTEGER, want STRING|ARRAY|NULL"}},
		{nil, `print("Hello, World")`, nil},
		{nil, `print("Hello", "World")`, nil},
		{nil, `println("Hello, World")`, nil},
		{nil, `println("Hello", " ", "World")`, nil},
		{nil, `string(123)`, "123"},
		{nil, `string(45.67)`, "45.670000"},
		{nil, `string(true)`, "true"},
		{nil, `string(false)`, "false"},
		{nil, `int(123)`, 123},
		{nil, `int(45.67)`, 45},
		{nil, `int("789")`, 789},
		{nil, `int(0)`, 0},
		{nil, `int(-42)`, -42},
		{nil, `int(-3.99)`, -3},
		{nil, `int("0")`, 0},
		{nil, `int("   456   ")`, 456},
		{nil, `int(true)`, 1},
		{nil, `int(false)`, 0},
		// {nil, `int(null)`, 0},f
		{nil, `int("")`, &objects.Error{Message: "error in `int`: failed to convert `` to INTEGER"}},
		{nil, `int("not a number")`, &objects.Error{Message: "error in `int`: failed to convert `not a number` to INTEGER"}},
		{nil, `int(1, 2)`, &objects.Error{Message: "wrong number of arguments to `int`: got 2, want 1"}},
		{nil, `float(123)`, 123.0},
		{nil, `float(45.67)`, 45.67},
		{nil, `float("789.01")`, 789.01},
		{nil, `float(0)`, 0.0},
		{nil, `float(-42)`, -42.0},
		{nil, `float(-3.99)`, -3.99},
		{nil, `float("0")`, 0.0},
		{nil, `float("   456.78   ")`, 456.78},
		{nil, `float(true)`, 1.0},
		{nil, `float(false)`, 0.0},
		{nil, `float(null)`, 0.0},
		{nil, `float("")`, &objects.Error{Message: "error in `float`: failed to convert `` to FLOAT"}},
		{nil, `float("not a number")`, &objects.Error{Message: "error in `float`: failed to convert `not a number` to FLOAT"}},
		{nil, `float(1, 2)`, &objects.Error{Message: "wrong number of arguments to `float`: got 2, want 1"}},
		{nil, `type(123)`, "INTEGER"},
		{nil, `type(45.67)`, "FLOAT"},
		{nil, `type("hello")`, "STRING"},
		{nil, `type(true)`, "BOOLEAN"},
		{nil, `type(false)`, "BOOLEAN"},
		{nil, `type([1, 2, 3])`, "ARRAY"},
		{nil, `type({'test': 'value'})`, "HASH"},
		{nil, `type(func() { })`, "FUNCTION"},
		{"isNaN with NaN input", `isNaN(0/0)`, true},
		{"isNaN with 1 over 0", `isNaN(1/0)`, false},
		{"isNaN with -1 over 0", `isNaN(-1/0)`, false},
		{"isNaN with 5", `isNaN(5)`, false},
		{"isNaN with 3.14", `isNaN(3.14)`, false},
		{"isNaN with string", `isNaN("not a number")`, false},
		{"isNaN with true", `isNaN(true)`, false},
		{"isNaN with null", `isNaN(null)`, false},
	}

	runVmTests(t, tests)
}

func TestGlobalBuiltinFunctions(t *testing.T) {
	timer.SetTimezone("UTC")

	tests := []vmTestCase{
		{"strings.contains with same case", `strings.contains("Hello, World", "World")`, true},
		{"strings.contains with different case", `strings.contains("Hello, World", "world")`, false},
		{"strings.split", `strings.split("a,b,c", ",")`, []string{"a", "b", "c"}},
		{"strings.join", `strings.join(["a", "b", "c"], ",")`, "a,b,c"},
		{"strings.format", `strings.format("Hello, %s! You have %d new messages.", "Alice", 5)`, "Hello, Alice! You have 5 new messages."},
		{"strings.startsWith with prefix", `strings.startsWith("hello", "he")`, true},
		{"strings.startsWith with full string", `strings.startsWith("hello", "hello")`, true},
		{"strings.startsWith with different prefix", `strings.startsWith("hello", "world")`, false},
		{"strings.startsWith with array of strings", `strings.startsWith("hello", ["Hello", "World"])`, false},
		{"strings.startsWith with array containing prefix", `strings.startsWith("hello", ["he", "world"])`, true},
		{"strings.endsWith with suffix", `strings.endsWith("hello", "lo")`, true},
		{"strings.endsWith with full string", `strings.endsWith("hello", "hello")`, true},
		{"strings.endsWith with different suffix", `strings.endsWith("hello", "world")`, false},
		{"strings.endsWith with array of strings", `strings.endsWith("hello", ["Hello", "World"])`, false},
		{"strings.endsWith with array containing suffix", `strings.endsWith("hello", ["lo", "world"])`, true},
		{"strings.toUpper", `strings.toUpper("Hello World!")`, "HELLO WORLD!"},
		{"strings.toUpper with empty string", `strings.toUpper("")`, ""},
		{"strings.toLower", `strings.toLower("Hello WORLD!")`, "hello world!"},
		{"strings.toLower with empty string", `strings.toLower("")`, ""},
		{"strings.trim", `strings.trim("  spaced  ")`, "spaced"},
		{"strings.trim with whitespace", `strings.trim("\n\tHello World\t\n")`, "Hello World"},
		{"strings.trim with custom characters", `strings.trim("xxxhelloxxx", "x")`, "hello"},
		{"arrays.push", `arrays.push([1, 2, 3], 4)`, []int{1, 2, 3, 4}},
		{"arrays.push with strings", `arrays.push(["a", "b", "c"], "d")`, []string{"a", "b", "c", "d"}},
		{"arrays.shift", `arrays.shift([1, 2, 3])`, 1},
		{"arrays.shift with variable", `var arr = [1, 2, 3]; arrays.shift(arr); arr`, []int{2, 3}},
		{"arrays.pop", `arrays.pop([1, 2, 3])`, 3},
		{"arrays.pop with variable", `var arr = [1, 2, 3]; arrays.pop(arr); arr`, []int{1, 2}},
		{"arrays.filter with variable", `var arr = [1, 2, 3, 4, 5]; arrays.filter(arr, func(x) { x % 2 == 0}); arr`, []int{1, 2, 3, 4, 5}},
		{"arrays.filter", `arrays.filter([1, 2, 3, 4, 5, 6], func(x) { x % 2 == 0 });`, []int{2, 4, 6}},
		{"arrays.concat empty", `arrays.concat([], [])`, []int{}},
		{"arrays.concat two arrays", `arrays.concat([1, 2], [3, 4])`, []int{1, 2, 3, 4}},
		{"arrays.concat two string arrays", `arrays.concat(["a", "b"], ["c", "d"])`, []string{"a", "b", "c", "d"}},
		{"arrays.concat three arrays", `arrays.concat([1, 2], [3, 4], [5, 6])`, []int{1, 2, 3, 4, 5, 6}},
		{"arrays.first", "arrays.first([100, 200, 300], func (x) { x >= 100 })", 100},
		{"arrays.first", "arrays.first([100, 200, 300], func (x) { x > 100 })", 200},
		{"arrays.first", "arrays.first([100, 200, 300], func (x, i) { i == 2 })", 300},
		{"arrays.first", "arrays.first([100, 200, 300], func (x) { x > 500 })", nil},
		{
			"arrays.first with invalid first argument",
			"arrays.first(5, func (a) { })",
			&objects.Error{Message: "argument 1 to `first` has invalid type: got INTEGER, want ARRAY"},
		},
		{
			"arrays.first with wrong number of arguments",
			"arrays.first([100, 200, 300], func () { })",
			&objects.Error{Message: "error in `first`: function passed to `first` must take at least one argument"},
		},
		{
			"arrays.first with too many arguments",
			"arrays.first([100, 200, 300], func (a, b, c) { })",
			&objects.Error{Message: "error in `first`: function passed to `first` must take at most two arguments"},
		},
		{
			"arrays.sort positive ints",
			"arrays.sort([5, 3, 1, 4, 2])",
			[]int{1, 2, 3, 4, 5},
		},
		{
			"arrays.sort mixed ints",
			"arrays.sort([-1, 0, 3, -5, 2, 1])",
			[]int{-5, -1, 0, 1, 2, 3},
		},
		{
			"arrays.sort empty",
			"arrays.sort([])",
			[]int{},
		},
		{
			"arrays.sort variable of positive ints",
			"var x = [3, 2, 1]; arrays.sort(x);",
			[]int{1, 2, 3},
		},
		{
			"arrays.sort variable of mixed ints",
			"var x = [10, 5, -5, -10, 0]; arrays.sort(x);",
			[]int{-10, -5, 0, 5, 10},
		},
		{
			"arrays.sort positive floats",
			"arrays.sort([1.5, 2.2, 0.3, -1.1])",
			[]float64{-1.1, 0.3, 1.5, 2.2},
		},
		{
			"arrays.sort strings",
			"arrays.sort(['banana', 'apple', 'cherry'])",
			[]string{"apple", "banana", "cherry"},
		},
		{
			"arrays.sort variable of strings",
			"var x = ['zen', 'lang', 'is', 'awesome']; arrays.sort(x);",
			[]string{"awesome", "is", "lang", "zen"},
		},
		{
			"arrays.sort booleans",
			"arrays.sort([true, false, true, false])",
			[]bool{false, false, true, true},
		},
		{
			"arrays.sort integers with custom comparator",
			"arrays.sort([5, 3, 1, 4, 2], func (a, b) { a < b })",
			[]int{1, 2, 3, 4, 5},
		},
		{
			"arrays.sort integers with custom comparator",
			"arrays.sort([5, 3, 1, 4, 2], func (a, b) { a > b })",
			[]int{5, 4, 3, 2, 1},
		},
		{
			"arrays.sort strings with custom comparator",
			"arrays.sort(['bb', 'a', 'ccc'], func (a, b) { len(a) < len(b) })",
			[]string{"a", "bb", "ccc"},
		},
		{
			"arrays.sort strings with custom comparator",
			"arrays.sort(['bb', 'a', 'ccc'], func (a, b) { len(a) > len(b) })",
			[]string{"ccc", "bb", "a"},
		},
		{nil, `math.min(19, 42)`, 19},
		{nil, `math.max(19, 42)`, 42},
		{nil, `math.ceil(12.34)`, 13.0},
		{nil, `math.floor(98.76)`, 98.0},
		{nil, `math.round(12.49)`, 12.0},
		{nil, `math.round(12.50)`, 13.0},
		{nil, `math.log(10 ^ 42)`, 42.0},
		{nil, `math.sqrt(1764)`, 42.0},
		{
			"time.parse with 'd-m-Y h:i:s' format",
			`time.parse("02-01-2026 16:23:48", "%d-%m-%Y %h:%i:%s")`,
			1767371028000,
		},
		{
			"time.parse with 'Y/m/d H:i:s A' format",
			`time.parse("2026/03/25 08:45:25 PM", "%Y/%m/%d %H:%i:%s %A")`,
			1774471525000,
		},
		{
			"time.parse with 'm-d-Y H:i:s a' format",
			`time.parse("03-25-2026 08:45:25 am", "%m-%d-%Y %H:%i:%s %a")`,
			1774428325000,
		},
		{
			"time.parse with 'D, M, y' format",
			`time.parse("Fri, Jan, 26", "%D, %M, %y")`,
			1767225600000,
		},
		{
			"time.parse with 'D, M, y' format",
			`time.parse("Fri, Jan, 26", "%D, %M, %y")`,
			1767225600000,
		},
		{
			"time.parse with 'd F Y' format",
			`time.parse("27 February 1993","%d %F %Y")`,
			730771200000,
		},
		{
			"time.format with 'Y-m-d H:i:s' format",
			`time.format(1767371028000, "%Y-%m-%d %H:%i:%s")`, "2026-01-02 04:23:48",
		},
		{
			"time.format with 'Y/m/d H-i-s' format",
			`time.format(1774471525000, "%Y/%m/%d %H-%i-%s")`, "2026/03/25 08-45-25",
		},
		{
			"time.format with 'd-m-Y s:i-H' format",
			`time.format(1774428325000, "%d-%m-%Y %s:%i:%H")`, "25-03-2026 25:45:08",
		},
		{
			"time.format with 'Y-m-d' format",
			`time.format(1767225600000, "%Y-%m-%d")`,
			"2026-01-01",
		},
		{
			"time.format with 'Y/m/d' format",
			`time.format(1767225600000, "%Y/%m/%d")`,
			"2026/01/01",
		},
		{
			"time.format with 'Y' format",
			`time.format(730771200000, "%Y")`,
			"1993",
		},
		{
			"time.format with 'D d, F, y' format",
			`time.format(730771200000, "%D %d, %F, %y")`,
			"Sat 27, February, 93",
		},
		{
			"time.format with 'H:i:s a' format",
			`time.format(1767606155000, "%H:%i:%s %a")`,
			"09:42:35 am",
		},
		{
			"time.format with 'H:i:s A' format",
			`time.format(1767606155000, "%H:%i:%s %A")`,
			"09:42:35 AM",
		},
		{
			"time.format with timezone 'America/New_York' and 'H:i:s A' format",
			`time.timezone("America/New_York"); time.format(1767606155000, "%H:%i:%s %A")`,
			"04:42:35 AM",
		},
		{
			"time.format with timezone 'Asia/Tokyo' and 'H:i:s A' format",
			`time.timezone("Asia/Tokyo"); time.format(1767606155000, "%H:%i:%s %A")`,
			"06:42:35 PM",
		},
		{
			"time.format with timezone 'Europe/Copenhagen' and 'H:i:s A' format",
			`time.timezone("Europe/Copenhagen"); time.format(1767606155000, "%H:%i:%s %A")`,
			"10:42:35 AM",
		},
		{
			"time.format with timezone 'UTC' and 'H:i:s A' format",
			`time.timezone("UTC"); time.format(1767606155000, "%H:%i:%s %A")`,
			"09:42:35 AM",
		},
	}

	runVmTests(t, tests)

	timer.ResetTimezone()
}

func TestClosures(t *testing.T) {
	tests := []vmTestCase{
		{
			name: "creating and calling a closure",
			input: `
				var newClosure = func(x) {
					return func() {
						return x
					}
				}

				var closure = newClosure(99)
				closure()
			`,
			expected: 99,
		},
		{
			name: "calling closure with free variables",
			input: `
				var newAdder = func(a, b) {
					return func(c) {
						return a + b + c
					}
				}

				var adder = newAdder(1, 2)
				adder(8)
			`,
			expected: 11,
		},
		{
			name: "calling closure that remembers binding from outer function",
			input: `
				var newAdder = func(a, b) {
					var c = a + b
					return func(d) { c + d }
				}

				var adder = newAdder(1, 2)
				adder(8)
			`,
			expected: 11,
		},
		{
			name: "calling nested closures that remember bindings from outer functions",
			input: `
				var newAdderOuter = func(a, b) {
					var c = a + b
					func(d) {
						var e = d + c
						func(f) { e + f }
					}
				}

				var newAdderInner = newAdderOuter(1, 2)
				var adder = newAdderInner(3)
				adder(8)
			`,
			expected: 14,
		},
		{
			name: "calling closure with free variables from multiple scopes",
			input: `
				var a = 1
				var newAdderOuter = func(b) {
					func(c) {
						func(d) { a + b + c + d }
					}
				}

				var newAdderInner = newAdderOuter(2)
				var adder = newAdderInner(3)
				adder(8)
			`,
			expected: 14,
		},
		{
			name: "calling closure that defines inner closures",
			input: `
				var newClosure = func(a, b) {
					var one = func() { a }
					var two = func() { b }
					func() { one() + two() }
				}

				var closure = newClosure(9, 90)
				closure()
			`,
			expected: 99,
		},
	}

	runVmTests(t, tests)
}

func BenchmarkClosures(b *testing.B) {
	runVmBenchmark(b, `
		var newAdder = func(a, b) {
			return func(c) {
				return a + b + c
			}
		}

		var adder = newAdder(1, 2)
		adder(8)
	`)
}

func TestRecursiveFunctions(t *testing.T) {
	tests := []vmTestCase{
		{
			name: "countDown function with recursion",
			input: `
				var countDown = func(x) {
					if (x == 0) {
						return 0
					} else {
						countDown(x - 1)
					}
				}

				countDown(1)
			`,
			expected: 0,
		},
		{
			name: "countDown function with recursion in wrapper",
			input: `
				var countDown = func(x) {
					if (x == 0) {
						return 0
					} else {
						countDown(x - 1)
					}
				}

				var wrapper = func() {
					countDown(1)
				}

				wrapper()
			`,
			expected: 0,
		},
		{
			name: "countDown function with recursion in wrapper and nested definition",
			input: `
				var wrapper = func() {
					var countDown = func(x) {
						if (x == 0) {
							return 0
						} else {
							countDown(x - 1)
						}
					}

					countDown(1)
				}

				wrapper()
			`,
			expected: 0,
		},
	}

	runVmTests(t, tests)
}

func BenchmarkRecursiveFunctions(b *testing.B) {
	runVmBenchmark(b, `
		var countDown = func(x) {
			if (x == 0) {
				return 0
			} else {
				countDown(x - 1)
			}
		}

		countDown(10)
	`)
}

func TestNamedFunctions(t *testing.T) {
	tests := []vmTestCase{
		{
			name: "simple named function",
			input: `
				func foo() {
					return 42;
				}

				foo();
			`,
			expected: 42,
		},
		{
			name: "simple sum function",
			input: `
				func sum(a, b) {
					return a + b;
				}

				sum(1, 2);
			`,
			expected: 3,
		},
		{
			name: "fibonacci function",
			input: `
				func fibonacci(x) {
					if (x <= 1) {
						return x
					}

					return fibonacci(x - 1) + fibonacci(x - 2);
				}

				fibonacci(10);
			`,
			expected: 55,
		},
		{
			name: "outer function with inner function",
			input: `
				func outer() {
					func inner() {
						return 99;
					}

					inner();
				}

				outer();
			`,
			expected: 99,
		},
		{
			name: "outer function with inner function returning inner function",
			input: `
				func outer() {
					func inner() {
						return 99;
					}
					return inner;
				}

				var innerFunc = outer();
				innerFunc();
			`,
			expected: 99,
		},
	}

	runVmTests(t, tests)
}

func BenchmarkNamedFunctions(b *testing.B) {
	runVmBenchmark(b, `
		func fibonacci(x) {
			if (x <= 1) {
				return x
			}

			return fibonacci(x - 1) + fibonacci(x - 2);
		}

		fibonacci(10);
	`)
}

func TestWhileLoops(t *testing.T) {
	tests := []vmTestCase{
		{
			name: "simple while loop",
			input: `
				var i = 0;
				var result = [];

				while (i < 10) {
					arrays.push(result, i);
					i++;
				}

				result;
			`,
			expected: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			name: "while loop with condition",
			input: `
				var i = 0;
				var result = [];

				while (i < 10) {
					if (i % 2 == 0) {
						arrays.push(result, i);
					}

					i++;
				}

				result;
			`,
			expected: []int{0, 2, 4, 6, 8},
		},
		{
			name: "while loop with continue",
			input: `
				var i = 0;
				var result = [];

				while (i < 10) {
					i++;

					if (i % 2 == 0) {
						continue;
					}

					arrays.push(result, i);
				}

				result;
			`,
			expected: []int{1, 3, 5, 7, 9},
		},
		{
			name: "while loop with break",
			input: `
				var i = 0;
				var result = [];

				while (i < 1000) {
					i++;

					if (i % 5 == 0) {
						break;
					}

					arrays.push(result, i);
				}

				result;
			`,
			expected: []int{1, 2, 3, 4},
		},
	}

	runVmTests(t, tests)
}

func BenchmarkWhileLoops(b *testing.B) {
	runVmBenchmark(b, `
		var i = 0;
		var result = [];

		while (i < 10) {
			arrays.push(result, i);
			i++;
		}

		result;
	`)
}

func TestAssignmentExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"simple assignment", "var mut a = 5; a = 10; a;", 10},
		{"assignment with addition", "var mut a = 5; a = a + 5; a;", 10},
		{"assignment with variable copy", "var mut a = 5; var mut b = 10; a = b; a;", 10},
		{"assignment with variable copy reversed", "var mut a = 5; var mut b = 10; b = a; b;", 5},
		{"assignment with addition reversed", "var mut a = 5; var mut b = 10; a = b + a; a;", 15},
		{"assignment with array element update", "var mut a = [1, 2, 3]; a[1] = 20; a;", []int{1, 20, 3}},
		{"assignment with array element copy", "var mut a = [1, 2, 3]; a[1] = a[0]; a;", []int{1, 1, 3}},
		{"assignment with array element addition", "var mut a = [1, 2, 3]; a[1] = a[0] + a[2]; a;", []int{1, 4, 3}},
		{"assignment with array element sum", "var mut a = [1, 2, 3]; a[1] = a[0] + a[1] + a[2]; a;", []int{1, 6, 3}},
		{"assignment with map element update", "var mut a = {'x': 1, 'y': 2}; a['y'] = 20; a;", map[string]int{"x": 1, "y": 20}},
		{"assignment with map element addition", "var mut a = {'x': 1, 'y': 2}; a['z'] = 20; a;", map[string]int{"x": 1, "y": 2, "z": 20}},
		{"assignment with map element copy", "var mut a = {'x': 1, 'y': 2}; a['x'] = a['y']; a;", map[string]int{"x": 2, "y": 2}},
		{"assignment with map element sum", "var mut a = {'x': 1, 'y': 2}; a['y'] = a['x'] + a['y']; a;", map[string]int{"x": 1, "y": 3}},
	}

	runVmTests(t, tests)
}

func BenchmarkAssignmentExpressions(b *testing.B) {
	runVmBenchmark(b, "var mut a = 5; var mut b = 10; a = b + a; a;")
}

func TestCompoundAssignments(t *testing.T) {
	tests := []vmTestCase{
		{"increment by 5", "var mut x = 5; x += 5;", 10},
		{"decrement by 5", "var mut x = 5; x -= 5;", 0},
		{"multiplication by 5", "var mut x = 5; x *= 5;", 25},
		{"division by 5", "var mut x = 5; x /= 5;", 1},
		{"modulus by 5", "var mut x = 5; x %= 5;", 0},
		{"exponentiation by 5", "var mut x = 5; x ^= 5;", 3125},
	}

	runVmTests(t, tests)
}

func BenchmarkCompoundAssignments(b *testing.B) {
	runVmBenchmark(b, "var mut x = 5; x ^= 5;")
}
