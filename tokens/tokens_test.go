package tokens

import "testing"

func TestLookupIdentForKeywords(t *testing.T) {
	for ident, expectedType := range keywords {
		if tokType := LookupIdent(ident); tokType != expectedType {
			t.Errorf("LookupIdent(%q) = %q; want %q", ident, tokType, expectedType)
		}
	}
}

func BenchmarkLookupIdentForKeywords(b *testing.B) {
	for b.Loop() {
		for ident, expectedType := range keywords {
			if tokType := LookupIdent(ident); tokType != expectedType {
				b.Errorf("LookupIdent(%q) = %q; want %q", ident, tokType, expectedType)
			}
		}
	}
}

func TestLookupIdentForNonKeywords(t *testing.T) {
	nonKeywords := []string{
		"x",
		"myVar",
		"functionName",
		"trueValue",
		"nullValue",
	}

	for _, ident := range nonKeywords {
		t.Run("lookup ident for keyword: "+ident, func(t *testing.T) {
			if tokType := LookupIdent(ident); tokType != IDENT {
				t.Errorf("LookupIdent(%q) = %q; want %q", ident, tokType, IDENT)
			}
		})
	}
}

func BenchmarkLookupIdentForNonKeywords(b *testing.B) {
	nonKeywords := []string{
		"x",
		"myVar",
		"functionName",
		"trueValue",
		"nullValue",
	}

	for b.Loop() {
		for _, ident := range nonKeywords {
			if tokType := LookupIdent(ident); tokType != IDENT {
				b.Errorf("LookupIdent(%q) = %q; want %q", ident, tokType, IDENT)
			}
		}
	}
}
