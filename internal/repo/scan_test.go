package repo

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

// makeRepo 在 dir 下创建一个伪 Git 仓库（仅 .git 目录）。
func makeRepo(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(dir, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
}

func TestScan(t *testing.T) {
	root := t.TempDir()
	makeRepo(t, filepath.Join(root, "a"))                // 深度 1
	makeRepo(t, filepath.Join(root, "b"))                // 深度 1
	makeRepo(t, filepath.Join(root, "nest", "c"))        // 深度 2
	makeRepo(t, filepath.Join(root, "a", "inner"))       // 仓库内部，应被忽略
	_ = os.MkdirAll(filepath.Join(root, "plain"), 0o755) // 非仓库目录

	names := func(ts []Target) []string {
		var out []string
		for _, tgt := range ts {
			out = append(out, tgt.Name)
		}
		sort.Strings(out)
		return out
	}

	t.Run("depth 1 finds only top level", func(t *testing.T) {
		got := names(Scan(root, 1))
		want := []string{"a", "b"}
		if !equal(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("depth 2 finds nested but not repo-internal", func(t *testing.T) {
		got := names(Scan(root, 2))
		want := []string{"a", "b", filepath.Join("nest", "c")}
		if !equal(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("depth 0 checks root itself", func(t *testing.T) {
		makeRepo(t, root)
		got := Scan(root, 0)
		if len(got) != 1 || got[0].Path != root {
			t.Errorf("got %+v, want single target for root", got)
		}
	})
}

func TestScanSkipsDefaultExcludedDirectories(t *testing.T) {
	root := t.TempDir()
	makeRepo(t, filepath.Join(root, "node_modules"))
	makeRepo(t, filepath.Join(root, "vendor"))
	makeRepo(t, filepath.Join(root, "vender"))
	makeRepo(t, filepath.Join(root, "archive"))
	makeRepo(t, filepath.Join(root, ".hidden"))
	makeRepo(t, filepath.Join(root, "project"))

	got := Scan(root, 1)
	if len(got) != 1 || got[0].Name != "project" {
		t.Fatalf("got %+v, want only project", got)
	}
}

func equal(a, b []string) bool {
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
