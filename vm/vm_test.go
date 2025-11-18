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
		{"2 * 3 ^ 4", 162},
		{"2 * 3 ^ 4 % 5", 2},
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

func TestSuffixExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"var x = 5; x++", 6},
		{"var x = 5; x--", 4},
		{"var x = 5; x++; x", 6},
		{"var x = 5; x--; x", 4},
		{"var x = 5; x++; x++; x", 7},
		{"var x = 5; x--; x--; x", 3},
		{"var x = 5; x++; x--; x", 5},
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
		{"[1][-1]", 1},
		{"[1][-2]", nil},
		{"{1: 1, 2: 2}[1]", 1},
		{"{1: 1, 2: 2}[2]", 2},
		{"{1: 1}[0]", nil},
		{"{}[0]", nil},
	}

	runVmTests(t, tests)
}

func TestChainIndexExpressions(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
				var obj = {'a': 2};
				obj.a
			`,
			expected: 2,
		},
		{
			input: `
				var obj = {'a': {'b': 3}};
				obj.a.b
			`,
			expected: 3,
		},
		{
			input: `
				var obj = {'a': [4]};
				obj.a[0]
			`,
			expected: 4,
		},
		{
			input: `
				var obj = {'a': {'b': func() { 5 }}};
				obj.a.b()
			`,
			expected: 5,
		},
		{
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

func TestCallingFunctionsWithoutArguments(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
				var fivePlusTen = func() { 5 + 10; };
				fivePlusTen();
			`,
			expected: 15,
		},
		{
			input: `
				var one = func() { 1; };
				var two = func() { 2; };
				one() + two()
			`,
			expected: 3,
		},
		{
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

func TestFunctionsWithReturnStatement(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
				var earlyExit = func() { return 99; 100; };
				earlyExit();
			`,
			expected: 99,
		},
		{
			input: `
				var earlyExit = func() { return 99; return 100; };
				earlyExit();
			`,
			expected: 99,
		},
	}

	runVmTests(t, tests)
}

func TestFunctionsWithoutReturnValue(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
				var noReturn = func() { };
				noReturn();
			`,
			expected: nil,
		},
		{
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

func TestCallingFunctionsWithBindings(t *testing.T) {
	tests := []vmTestCase{
		{
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

func TestCallingFunctionsWithArgumentsAndBindings(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
				var value = func(a) { a; };
				value(5);
			`,
			expected: 5,
		},
		{
			input: `
				var sum = func(a, b) { a + b; };
				sum(1, 2);
			`,
			expected: 3,
		},
		{
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

func TestCallingFunctionsWithWrongArguments(t *testing.T) {
	tests := []vmTestCase{
		{
			input:    `func() { 1; }(1);`,
			expected: `wrong number of arguments. got 1, want 0`,
		},
		{
			input:    `func(a) { a; }();`,
			expected: `wrong number of arguments. got 0, want 1`,
		},
		{
			input:    `func(a, b) { a + b; }(1);`,
			expected: `wrong number of arguments. got 1, want 2`,
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
			t.Fatalf("wrong VM error: want %q, got %q", tt.expected, err)
		}
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []vmTestCase{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len([])`, 0},
		{`len([1, 2, 3])`, 3},
		{`len(1)`, &objects.Error{Message: "argument to `len` not supported, got INTEGER"}},
		{`print("Hello, World")`, nil},
		{`print("Hello", "World")`, nil},
		{`println("Hello, World")`, nil},
		{`println("Hello", " ", "World")`, nil},
		{`string(123)`, "123"},
		{`string(45.67)`, "45.670000"},
		{`string(true)`, "true"},
		{`string(false)`, "false"},
	}

	runVmTests(t, tests)
}

func TestGlobalBuiltinFunctions(t *testing.T) {
	tests := []vmTestCase{
		{`strings.contains("Hello, World", "World")`, true},
		{`strings.contains("Hello, World", "world")`, false},
		{`strings.split("a,b,c", ",")`, []string{"a", "b", "c"}},
		{`strings.join(["a", "b", "c"], ",")`, "a,b,c"},
		{`strings.format("Hello, %s! You have %d new messages.", "Alice", 5)`, "Hello, Alice! You have 5 new messages."},
		{`arrays.push([1, 2, 3], 4)`, []int{1, 2, 3, 4}},
		{`arrays.push(["a", "b", "c"], "d")`, []string{"a", "b", "c", "d"}},
		{`arrays.shift([1, 2, 3])`, 1},
		{`var arr = [1, 2, 3]; arrays.shift(arr); arr`, []int{2, 3}},
		{`arrays.pop([1, 2, 3])`, 3},
		{`var arr = [1, 2, 3]; arrays.pop(arr); arr`, []int{1, 2}},
		{`var arr = [1, 2, 3, 4, 5]; arrays.filter(arr, func(x) { x % 2 == 0}); arr`, []int{1, 2, 3, 4, 5}},
		{`arrays.filter([1, 2, 3, 4, 5, 6], func(x) { x % 2 == 0 });`, []int{2, 4, 6}},
		{`arrays.concat([], [])`, []int{}},
		{`arrays.concat([1, 2], [3, 4])`, []int{1, 2, 3, 4}},
		{`arrays.concat(["a", "b"], ["c", "d"])`, []string{"a", "b", "c", "d"}},
		{`arrays.concat([1, 2], [3, 4], [5, 6])`, []int{1, 2, 3, 4, 5, 6}},
		{`math.min(19, 42)`, 19},
		{`math.max(19, 42)`, 42},
		{`math.ceil(12.34)`, 13.0},
		{`math.floor(98.76)`, 98.0},
		{`math.round(12.49)`, 12.0},
		{`math.round(12.50)`, 13.0},
		{`math.log(10 ^ 42)`, 42.0},
		{`math.sqrt(1764)`, 42.0},
	}

	runVmTests(t, tests)
}

func TestClosures(t *testing.T) {
	tests := []vmTestCase{
		{
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

func TestRecursiveFunctions(t *testing.T) {
	tests := []vmTestCase{
		{
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

func TestNamedFunctions(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
				func foo() {
					return 42;
				}

				foo();
			`,
			expected: 42,
		},
		{
			input: `
				func sum(a, b) {
					return a + b;
				}

				sum(1, 2);
			`,
			expected: 3,
		},
		{
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

func TestWhileLoops(t *testing.T) {
	tests := []vmTestCase{
		{
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

func TestAssignmentExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"var mut a = 5; a = 10; a;", 10},
		{"var mut a = 5; a = a + 5; a;", 10},
		{"var mut a = 5; var mut b = 10; a = b; a;", 10},
		{"var mut a = 5; var mut b = 10; b = a; b;", 5},
		{"var mut a = 5; var mut b = 10; a = b + a; a;", 15},
		{"var mut a = [1, 2, 3]; a[1] = 20; a;", []int{1, 20, 3}},
		{"var mut a = [1, 2, 3]; a[1] = a[0]; a;", []int{1, 1, 3}},
		{"var mut a = [1, 2, 3]; a[1] = a[0] + a[2]; a;", []int{1, 4, 3}},
		{"var mut a = [1, 2, 3]; a[1] = a[0] + a[1] + a[2]; a;", []int{1, 6, 3}},
		{"var mut a = {'x': 1, 'y': 2}; a['y'] = 20; a;", map[string]int{"x": 1, "y": 20}},
		{"var mut a = {'x': 1, 'y': 2}; a['z'] = 20; a;", map[string]int{"x": 1, "y": 2, "z": 20}},
		{"var mut a = {'x': 1, 'y': 2}; a['x'] = a['y']; a;", map[string]int{"x": 2, "y": 2}},
		{"var mut a = {'x': 1, 'y': 2}; a['y'] = a['x'] + a['y']; a;", map[string]int{"x": 1, "y": 3}},
	}

	runVmTests(t, tests)
}
