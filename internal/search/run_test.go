package search

import (
	"slices"
	"strings"
	"testing"
)

func TestCollectMatches(t *testing.T) {
	t.Run("keeps all matches without limit", func(t *testing.T) {
		result := collectMatches(strings.NewReader("a.go\x1f1\x1ffoo\n\nb.go\x1f2\x1fbar\n"), 0, func() {
			t.Fatal("stop should not be called without limit")
		})

		if result.truncated {
			t.Fatal("result should not be truncated")
		}
		if got, want := len(result.matches), 2; got != want {
			t.Fatalf("len(matches) = %d, want %d", got, want)
		}
		if got, want := result.matches[0], (match{path: "a.go", line: "1", text: "foo"}); got != want {
			t.Fatalf("matches[0] = %#v, want %#v", got, want)
		}
	})

	t.Run("truncates after limit", func(t *testing.T) {
		stopped := false
		result := collectMatches(strings.NewReader("1\n2\n3\n"), 2, func() {
			stopped = true
		})

		if !stopped {
			t.Fatal("stop should be called when result is truncated")
		}
		if !result.truncated {
			t.Fatal("result should be truncated")
		}
		if got, want := len(result.matches), 2; got != want {
			t.Fatalf("len(matches) = %d, want %d", got, want)
		}
	})
}

func TestBuildRipgrepArgs(t *testing.T) {
	args := buildRipgrepArgs("needle")

	for _, want := range []string{
		"!.git/**",
		"!node_modules/**",
		"!.idea/**",
		"!.vscode/**",
		"!.pnpm-store/**",
		"!.swc/**",
		"!.temp/**",
		"!.rn_temp/**",
		"!.cache/**",
		"!.nuxt/**",
		"!.output/**",
		"!.data/**",
		"!.nitro/**",
		"!.fleet/**",
		"!.DS_Store",
	} {
		if !slices.Contains(args, want) {
			t.Fatalf("args should contain exclude glob %q", want)
		}
	}

	if !slices.Contains(args, "--field-match-separator") {
		t.Fatal("args should contain --field-match-separator")
	}
	if !slices.Contains(args, rgFieldSeparator) {
		t.Fatal("args should contain rg field separator")
	}
	if got, want := args[len(args)-2], "--"; got != want {
		t.Fatalf("args second last = %q, want %q", got, want)
	}
	if got, want := args[len(args)-1], "needle"; got != want {
		t.Fatalf("args last = %q, want %q", got, want)
	}
}

func TestParseMatchLine(t *testing.T) {
	t.Run("keeps colons in text", func(t *testing.T) {
		got := parseMatchLine("a.go\x1f12\x1ffoo:bar:baz")
		want := match{path: "a.go", line: "12", text: "foo:bar:baz"}
		if got != want {
			t.Fatalf("parseMatchLine() = %#v, want %#v", got, want)
		}
	})

	t.Run("keeps raw line when rg output is unexpected", func(t *testing.T) {
		got := parseMatchLine("raw")
		want := match{text: "raw"}
		if got != want {
			t.Fatalf("parseMatchLine() = %#v, want %#v", got, want)
		}
	})
}
