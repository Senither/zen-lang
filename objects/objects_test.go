package objects

import "testing"

func TestGetBuiltinByName(t *testing.T) {
	for _, def := range Builtins {
		builtin := GetBuiltinByName(def.Name)
		if builtin == nil {
			t.Errorf("expected to find builtin with name '%s', but got nil", def.Name)
		}

		if builtin != def.Builtin {
			t.Errorf("expected to get the same builtin instance for name '%s'", def.Name)
		}
	}
}

func TestGetGlobalBuiltinByName(t *testing.T) {
	for _, grp := range Globals {
		for _, def := range grp.Builtins {
			builtin := GetGlobalBuiltinByName(grp.Name, def.Name)
			if builtin == nil {
				t.Errorf("expected to find global builtin with name '%s.%s', but got nil", grp.Name, def.Name)
			}

			if builtin != def.Builtin {
				t.Errorf("expected to get the same global builtin instance for name '%s.%s'", grp.Name, def.Name)
			}
		}
	}
}
