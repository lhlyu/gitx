package tag

import "testing"

func TestParseEntries(t *testing.T) {
	out := "v1.2.0\x1ftag\x1f2026-06-26 10:20\x1fa1b2c3d\x1f1111111\x1frelease: v1.2.0\x1ftag message\n" +
		"v1.1.0\x1fcommit\x1f2026-06-20 09:12\x1f\x1fe5f6a7b\x1f\x1ffeat: add log command\n" +
		"发布-候选\x1ftag\x1f2026-06-18 08:30\x1f1234567\x1f7654321\x1f修复: 中文标题\x1f标注消息\n" +
		"bad line\n"

	entries := parseEntries(out)
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}

	if entries[0].Name != "v1.2.0" {
		t.Fatalf("expected first tag v1.2.0, got %q", entries[0].Name)
	}
	if entries[0].Kind != "标注" {
		t.Fatalf("expected annotated tag kind 标注, got %q", entries[0].Kind)
	}
	if entries[0].Hash != "a1b2c3d" {
		t.Fatalf("expected annotated tag target hash, got %q", entries[0].Hash)
	}
	if entries[0].Subject != "release: v1.2.0" {
		t.Fatalf("expected annotated tag target subject, got %q", entries[0].Subject)
	}

	if entries[1].Kind != "轻量" {
		t.Fatalf("expected lightweight tag kind 轻量, got %q", entries[1].Kind)
	}
	if entries[1].Hash != "e5f6a7b" {
		t.Fatalf("expected lightweight tag commit hash, got %q", entries[1].Hash)
	}
	if entries[1].Subject != "feat: add log command" {
		t.Fatalf("expected lightweight tag commit subject, got %q", entries[1].Subject)
	}

	if entries[2].Name != "发布-候选" {
		t.Fatalf("expected Chinese tag name, got %q", entries[2].Name)
	}
	if entries[2].Subject != "修复: 中文标题" {
		t.Fatalf("expected Chinese subject, got %q", entries[2].Subject)
	}
}

func TestParseEntriesEmpty(t *testing.T) {
	entries := parseEntries("\n\n")
	if len(entries) != 0 {
		t.Fatalf("expected empty entries, got %d", len(entries))
	}
}
