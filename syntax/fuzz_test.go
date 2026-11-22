package syntax

import "testing"

func FuzzParseAndWrite(f *testing.F) {
	f.Add("")
	f.Add(`(?<word>\w+)\s+\k<word>`)
	f.Add(`(?<=a)b|[c-z]+`)
	f.Fuzz(func(t *testing.T, pattern string) {
		tree, err := Parse(pattern, ParseOptions{})
		if err != nil {
			return
		}

		if _, err := Write(tree); err != nil {
			t.Fatal(err)
		}
	})
}
