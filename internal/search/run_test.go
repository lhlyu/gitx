package search

import (
	"slices"
	"strings"
	"testing"
)

func TestCollectMatches(t *testing.T) {
	t.Run("keeps all matches without limit", func(t *testing.T) {
		result := collectMatches(strings.NewReader("a.go:1:foo\n\nb.go:2:bar\n"), 0, func() {
			t.Fatal("stop should not be called without limit")
		})

		if result.truncated {
			t.Fatal("result should not be truncated")
		}
		if got, want := len(result.matches), 2; got != want {
			t.Fatalf("len(matches) = %d, want %d", got, want)
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

	if got, want := args[len(args)-2], "--"; got != want {
		t.Fatalf("args second last = %q, want %q", got, want)
	}
	if got, want := args[len(args)-1], "needle"; got != want {
		t.Fatalf("args last = %q, want %q", got, want)
	}
}
