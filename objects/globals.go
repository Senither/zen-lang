package objects

var Globals = []struct {
	Name     string
	Builtins []*BuiltinDefinition
}{
	{
		Name: "strings",
		Builtins: []*BuiltinDefinition{
			{
				Name: "contains",
				Schema: BuiltinSchema{
					NewRequiredArgument(STRING_OBJ),
					NewRequiredArgument(STRING_OBJ),
				},
				Builtin: &Builtin{Fn: globalStringsContains},
			},
			{
				Name: "split",
				Schema: BuiltinSchema{
					NewRequiredArgument(STRING_OBJ),
					NewRequiredArgument(STRING_OBJ),
				},
				Builtin: &Builtin{Fn: globalStringsSplit},
			},
			{
				Name: "join",
				Schema: BuiltinSchema{
					NewRequiredArgument(ARRAY_OBJ),
					NewRequiredArgument(STRING_OBJ),
				},
				Builtin: &Builtin{Fn: globalStringsJoin},
			},
			{
				Name: "format",
				Schema: BuiltinSchema{
					NewRequiredArgument(ARRAY_OBJ),
					NewRequiredArgument(),
				},
				Builtin: &Builtin{Fn: globalStringsFormat},
			},
			{
				Name: "startsWith",
				Schema: BuiltinSchema{
					NewRequiredArgument(STRING_OBJ),
					NewRequiredArgument(STRING_OBJ, ARRAY_OBJ),
				},
				Builtin: &Builtin{Fn: globalStringsStartsWith},
			},
			{
				Name: "endsWith",
				Schema: BuiltinSchema{
					NewRequiredArgument(STRING_OBJ),
					NewRequiredArgument(STRING_OBJ, ARRAY_OBJ),
				},
				Builtin: &Builtin{Fn: globalStringsEndsWith},
			},
			{
				Name: "toUpper",
				Schema: BuiltinSchema{
					NewRequiredArgument(STRING_OBJ),
				},
				Builtin: &Builtin{Fn: globalStringsToUpper},
			},
			{
				Name: "toLower",
				Schema: BuiltinSchema{
					NewRequiredArgument(STRING_OBJ),
				},
				Builtin: &Builtin{Fn: globalStringsToLower},
			},
			{
				Name: "trim",
				Schema: BuiltinSchema{
					NewRequiredArgument(STRING_OBJ),
					NewOptionalArgument(STRING_OBJ),
				},
				Builtin: &Builtin{Fn: globalStringsTrim},
			},
		},
	},
	{
		Name: "arrays",
		Builtins: []*BuiltinDefinition{
			{
				Name: "push",
				Schema: BuiltinSchema{
					NewRequiredArgument(ARRAY_OBJ),
					NewRequiredArgument(),
				},
				Builtin: &Builtin{Fn: globalArraysPush},
			},
			{
				Name: "shift",
				Schema: BuiltinSchema{
					NewRequiredArgument(ARRAY_OBJ),
				},
				Builtin: &Builtin{Fn: globalArraysShift},
			},
			{
				Name: "pop",
				Schema: BuiltinSchema{
					NewRequiredArgument(ARRAY_OBJ),
				},
				Builtin: &Builtin{Fn: globalArraysPop},
			},
			{
				Name: "filter",
				Schema: BuiltinSchema{
					NewRequiredArgument(ARRAY_OBJ),
					NewRequiredArgument(FUNCTION_OBJ, CLOSURE_OBJ),
				},
				Builtin: &Builtin{Fn: globalArraysFilter},
			},
			{
				Name: "concat",
				Schema: BuiltinSchema{
					NewRequiredArgument(ARRAY_OBJ),
					NewRequiredArgument(ARRAY_OBJ),
					NewOptionalArgument(ARRAY_OBJ),
				},
				Builtin: &Builtin{Fn: globalArraysConcat},
			},
			{
				Name: "flatten",
				Schema: BuiltinSchema{
					NewRequiredArgument(ARRAY_OBJ),
				},
				Builtin: &Builtin{Fn: globalArraysFlatten},
			},
			{
				Name: "first",
				Schema: BuiltinSchema{
					NewRequiredArgument(ARRAY_OBJ),
					NewRequiredArgument(FUNCTION_OBJ, CLOSURE_OBJ),
				},
				Builtin: &Builtin{Fn: globalArraysFirst},
			},
			{
				Name: "sort",
				Schema: BuiltinSchema{
					NewRequiredArgument(ARRAY_OBJ),
					NewOptionalArgument(FUNCTION_OBJ, CLOSURE_OBJ),
				},
				Builtin: &Builtin{Fn: globalArraysSort},
			},
		},
	},
	{
		Name: "maps",
		Builtins: []*BuiltinDefinition{
			{
				Name: "keys",
				Schema: BuiltinSchema{
					NewRequiredArgument(HASH_OBJ),
				},
				Builtin: &Builtin{Fn: globalMapsKeys},
			},
			{
				Name: "values",
				Schema: BuiltinSchema{
					NewRequiredArgument(HASH_OBJ),
				},
				Builtin: &Builtin{Fn: globalMapsValues},
			},
			{
				Name: "has",
				Schema: BuiltinSchema{
					NewRequiredArgument(HASH_OBJ),
					NewRequiredArgument(STRING_OBJ, INTEGER_OBJ, FLOAT_OBJ, BOOLEAN_OBJ),
				},
				Builtin: &Builtin{Fn: globalMapsHas},
			},
			{
				Name: "each",
				Schema: BuiltinSchema{
					NewRequiredArgument(HASH_OBJ),
					NewRequiredArgument(FUNCTION_OBJ, CLOSURE_OBJ),
				},
				Builtin: &Builtin{Fn: globalMapsEach},
			},
			{
				Name: "merge",
				Schema: BuiltinSchema{
					NewRequiredArgument(HASH_OBJ),
					NewRequiredArgument(HASH_OBJ),
					NewOptionalArgument(HASH_OBJ),
				},
				Builtin: &Builtin{Fn: globalMapsMerge},
			},
		},
	},
	{
		Name: "math",
		Builtins: []*BuiltinDefinition{
			{
				Name: "min",
				Schema: BuiltinSchema{
					NewRequiredArgument(GetNumberTypes()...),
					NewRequiredArgument(GetNumberTypes()...),
				},
				Builtin: &Builtin{Fn: globalMathMin},
			},
			{
				Name: "max",
				Schema: BuiltinSchema{
					NewRequiredArgument(GetNumberTypes()...),
					NewRequiredArgument(GetNumberTypes()...),
				},
				Builtin: &Builtin{Fn: globalMathMax},
			},
			{
				Name: "ceil",
				Schema: BuiltinSchema{
					NewRequiredArgument(GetNumberTypes()...),
				},
				Builtin: &Builtin{Fn: globalMathCeil},
			},
			{
				Name: "floor",
				Schema: BuiltinSchema{
					NewRequiredArgument(GetNumberTypes()...),
				},
				Builtin: &Builtin{Fn: globalMathFloor},
			},
			{
				Name: "round",
				Schema: BuiltinSchema{
					NewRequiredArgument(GetNumberTypes()...),
				},
				Builtin: &Builtin{Fn: globalMathRound},
			},
			{
				Name: "log",
				Schema: BuiltinSchema{
					NewRequiredArgument(GetNumberTypes()...),
				},
				Builtin: &Builtin{Fn: globalMathLog},
			},
			{
				Name: "sqrt",
				Schema: BuiltinSchema{
					NewRequiredArgument(GetNumberTypes()...),
				},
				Builtin: &Builtin{Fn: globalMathSqrt},
			},
		},
	},
	{
		Name: "time",
		Builtins: []*BuiltinDefinition{
			{
				Name:    "now",
				Builtin: &Builtin{Fn: globalTimeNow},
			},
			{
				Name: "sleep",
				Schema: BuiltinSchema{
					NewRequiredArgument(INTEGER_OBJ, FLOAT_OBJ),
				},
				Builtin: &Builtin{Fn: globalTimeSleep},
			},
			{
				Name: "parse",
				Schema: BuiltinSchema{
					NewRequiredArgument(STRING_OBJ),
					NewRequiredArgument(STRING_OBJ),
				},
				Builtin: &Builtin{Fn: globalTimeParse},
			},
			{
				Name: "format",
				Schema: BuiltinSchema{
					NewRequiredArgument(INTEGER_OBJ, FLOAT_OBJ),
					NewRequiredArgument(STRING_OBJ),
				},
				Builtin: &Builtin{Fn: globalTimeFormat},
			},
			{
				Name: "timezone",
				Schema: BuiltinSchema{
					NewRequiredArgument(STRING_OBJ),
				},
				Builtin: &Builtin{Fn: globalTimeTimezone},
			},
			{
				Name: "delayTimer",
				Schema: BuiltinSchema{
					NewRequiredArgument(FUNCTION_OBJ, CLOSURE_OBJ),
					NewRequiredArgument(INTEGER_OBJ),
				},
				Builtin: &Builtin{Fn: globalTimeDelayTimer},
			},
			{
				Name: "scheduleTimer",
				Schema: BuiltinSchema{
					NewRequiredArgument(FUNCTION_OBJ, CLOSURE_OBJ),
					NewRequiredArgument(INTEGER_OBJ),
				},
				Builtin: &Builtin{Fn: globalTimeScheduleTimer},
			},
		},
	},
	{
		Name: "process",
		Builtins: []*BuiltinDefinition{
			{
				Name:    "exit",
				Schema:  BuiltinSchema{NewRequiredArgument(INTEGER_OBJ)},
				Builtin: &Builtin{Fn: globalProcessExit},
			},
			{
				Name:    "argv",
				Builtin: &Builtin{Fn: globalProcessArgv},
			},
			{
				Name:    "env",
				Schema:  BuiltinSchema{NewRequiredArgument(STRING_OBJ)},
				Builtin: &Builtin{Fn: globalProcessEnv},
			},
		},
	},
	{
		Name: "json",
		Builtins: []*BuiltinDefinition{
			{
				Name:    "parse",
				Schema:  BuiltinSchema{NewRequiredArgument(STRING_OBJ)},
				Builtin: &Builtin{Fn: globalJSONParse},
			},
			{
				Name:    "stringify",
				Schema:  BuiltinSchema{NewRequiredArgument(HASH_OBJ, ARRAY_OBJ)},
				Builtin: &Builtin{Fn: globalJSONStringify},
			},
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
