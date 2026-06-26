package search

import (
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
