package repo

import "testing"

func TestParseStatus(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want Status
	}{
		{
			name: "clean tracked, synced",
			in:   "## main...origin/main\n",
			want: Status{Branch: "main", Upstream: "origin/main"},
		},
		{
			name: "ahead and behind with changes",
			in:   "## dev...origin/dev [ahead 2, behind 3]\n M f.txt\n?? new.txt\n",
			want: Status{Branch: "dev", Upstream: "origin/dev", Ahead: 2, Behind: 3, ChangedFiles: 2},
		},
		{
			name: "ahead only",
			in:   "## main...origin/main [ahead 1]\n",
			want: Status{Branch: "main", Upstream: "origin/main", Ahead: 1},
		},
		{
			name: "no upstream",
			in:   "## feature-x\n",
			want: Status{Branch: "feature-x", NoUpstream: true},
		},
		{
			name: "no commits yet",
			in:   "## No commits yet on main\n?? a.txt\n",
			want: Status{Branch: "main", NoUpstream: true, ChangedFiles: 1},
		},
		{
			name: "detached head",
			in:   "## HEAD (no branch)\n",
			want: Status{Branch: "(detached)", NoUpstream: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseStatus(tt.in)
			if got != tt.want {
				t.Errorf("ParseStatus(%q)\n got = %+v\nwant = %+v", tt.in, got, tt.want)
			}
		})
	}
}

func TestStatusIsClean(t *testing.T) {
	if !(Status{}).IsClean() {
		t.Error("zero Status should be clean")
	}
	if (Status{ChangedFiles: 1}).IsClean() {
		t.Error("Status with changes should not be clean")
	}
}

func TestFirstLine(t *testing.T) {
	tests := map[string]string{
		"":                             "",
		"\n\n":                         "",
		"error: pathspec":              "error: pathspec",
		"\n  fatal: not a git repo \n": "fatal: not a git repo",
	}
	for in, want := range tests {
		if got := FirstLine([]byte(in)); got != want {
			t.Errorf("FirstLine(%q) = %q, want %q", in, got, want)
		}
	}
}
