package objects

var Globals = []struct {
	Name     string
	Builtins []*BuiltinDefinition
}{
	{
		Name: "arrays",
		Builtins: []*BuiltinDefinition{
			{Name: "push", Builtin: &Builtin{Fn: globalArraysPush}},
			{Name: "shift", Builtin: &Builtin{Fn: globalArraysShift}},
			{Name: "pop", Builtin: &Builtin{Fn: globalArraysPop}},
			{Name: "filter", Builtin: &Builtin{Fn: globalArraysFilter}},
			{Name: "concat", Builtin: &Builtin{Fn: globalArraysConcat}},
			{Name: "first", Builtin: &Builtin{Fn: globalArraysFirst}},
			{Name: "sort", Builtin: &Builtin{Fn: globalArraysSort}},
		},
	},
	{
		Name: "strings",
		Builtins: []*BuiltinDefinition{
			{Name: "contains", Builtin: &Builtin{Fn: globalStringsContains}},
			{Name: "split", Builtin: &Builtin{Fn: globalStringsSplit}},
			{Name: "join", Builtin: &Builtin{Fn: globalStringsJoin}},
			{Name: "format", Builtin: &Builtin{Fn: globalStringsFormat}},
			{Name: "startsWith", Builtin: &Builtin{Fn: globalStringsStartsWith}},
			{Name: "endsWith", Builtin: &Builtin{Fn: globalStringsEndsWith}},
			{Name: "toUpper", Builtin: &Builtin{Fn: globalStringsToUpper}},
			{Name: "toLower", Builtin: &Builtin{Fn: globalStringsToLower}},
			{Name: "trim", Builtin: &Builtin{Fn: globalStringsTrim}},
		},
	},
	{
		Name: "math",
		Builtins: []*BuiltinDefinition{
			{Name: "min", Builtin: &Builtin{Fn: globalMathMin}},
			{Name: "max", Builtin: &Builtin{Fn: globalMathMax}},
			{Name: "ceil", Builtin: &Builtin{Fn: globalMathCeil}},
			{Name: "floor", Builtin: &Builtin{Fn: globalMathFloor}},
			{Name: "round", Builtin: &Builtin{Fn: globalMathRound}},
			{Name: "log", Builtin: &Builtin{Fn: globalMathLog}},
			{Name: "sqrt", Builtin: &Builtin{Fn: globalMathSqrt}},
		},
	},
	{
		Name: "time",
		Builtins: []*BuiltinDefinition{
			{Name: "now", Builtin: &Builtin{Fn: globalTimeNow}},
			{Name: "sleep", Builtin: &Builtin{Fn: globalTimeSleep}},
			{Name: "parse", Builtin: &Builtin{Fn: globalTimeParse}},
			{Name: "format", Builtin: &Builtin{Fn: globalTimeFormat}},
		},
	},
}

func GetGlobalBuiltinByName(scope, name string) *Builtin {
	for _, grp := range Globals {
		if grp.Name == scope {
			for _, def := range grp.Builtins {
				if def.Name == name {
					return def.Builtin
				}
			}
		}
	}

	return nil
}
