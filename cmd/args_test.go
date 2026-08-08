package cmd

import "testing"

func TestParseOptionalInt(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		defaultVal int
		minValue   int
		want       int
		wantError  bool
	}{
		{name: "default", defaultVal: 3, minValue: 0, want: 3},
		{name: "zero", args: []string{"0"}, minValue: 0, want: 0},
		{name: "positive", args: []string{"2"}, minValue: 1, want: 2},
		{name: "negative", args: []string{"-1"}, minValue: 0, wantError: true},
		{name: "not a number", args: []string{"x"}, minValue: 0, wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseOptionalInt(tt.args, tt.defaultVal, tt.minValue, "depth")
			if (err != nil) != tt.wantError {
				t.Fatalf("error = %v, wantError %v", err, tt.wantError)
			}
			if err == nil && got != tt.want {
				t.Fatalf("got %d, want %d", got, tt.want)
			}
		})
	}
}
