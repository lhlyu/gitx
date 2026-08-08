package reset

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunDryRunDoesNotModifyFiles(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(root, "marker.txt")
	if err := os.WriteFile(marker, []byte("keep"), 0o644); err != nil {
		t.Fatal(err)
	}
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(oldDir) }()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}

	if err := Run(1, true); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(marker); err != nil {
		t.Fatalf("dry-run modified marker: %v", err)
	}
}
