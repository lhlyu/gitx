package pull

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveDepth(t *testing.T) {
	root := t.TempDir()

	if got := resolveDepth(root, AutoDepth); got != 1 {
		t.Fatalf("non-repository directory: got depth %d, want 1", got)
	}

	if err := os.Mkdir(filepath.Join(root, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	if got := resolveDepth(root, AutoDepth); got != 0 {
		t.Fatalf("repository directory: got depth %d, want 0", got)
	}
}

func TestResolveDepthKeepsExplicitDepth(t *testing.T) {
	root := t.TempDir()
	for _, depth := range []int{0, 1, 2, 10} {
		if got := resolveDepth(root, depth); got != depth {
			t.Errorf("explicit depth %d: got %d", depth, got)
		}
	}
}

func TestPullArgs(t *testing.T) {
	if got := pullArgs(Options{}); len(got) != 1 || got[0] != "pull" {
		t.Fatalf("default args = %v, want [pull]", got)
	}
	got := pullArgs(Options{Rebase: true, Prune: true})
	want := []string{"pull", "--rebase", "--prune"}
	if !equalStrings(got, want) {
		t.Fatalf("args = %v, want %v", got, want)
	}
	got = pullArgs(Options{FFOnly: true, Prune: true})
	want = []string{"pull", "--ff-only", "--prune"}
	if !equalStrings(got, want) {
		t.Fatalf("args = %v, want %v", got, want)
	}
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
