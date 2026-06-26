package branch

import "testing"

func TestParseEntries(t *testing.T) {
	out := "*\x1frefs/heads/main\x1fmain\x1f51d331b\x1f2026-06-24 10:22\x1ffeat: main\n" +
		" \x1frefs/heads/refactor/demo\x1frefactor/demo\x1f342bd2a\x1f2026-06-18 17:57\x1ffeat: demo\n" +
		" \x1frefs/remotes/origin/HEAD\x1forigin/HEAD\x1f51d331b\x1f2026-06-24 10:22\x1ffeat: main\n" +
		" \x1frefs/remotes/origin/main\x1forigin/main\x1f51d331b\x1f2026-06-24 10:22\x1ffeat: main\n" +
		"bad line\n"

	entries := parseEntries(out)
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}

	if !entries[0].IsCurrent {
		t.Fatal("expected main to be current branch")
	}
	if entries[0].Kind != "本地" {
		t.Fatalf("expected main kind to be 本地, got %q", entries[0].Kind)
	}
	if entries[1].Kind != "本地" {
		t.Fatalf("expected refactor/demo kind to be 本地, got %q", entries[1].Kind)
	}
	if entries[2].Name != "origin/main" {
		t.Fatalf("expected remote branch origin/main, got %q", entries[2].Name)
	}
	if entries[2].Kind != "远程" {
		t.Fatalf("expected origin/main kind to be 远程, got %q", entries[2].Kind)
	}
}

func TestParseEntriesEmpty(t *testing.T) {
	entries := parseEntries("\n\n")
	if len(entries) != 0 {
		t.Fatalf("expected empty entries, got %d", len(entries))
	}
}
