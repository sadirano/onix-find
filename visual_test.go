package main

import "testing"

func TestFallback(t *testing.T) {
	tests := []struct {
		value, def, want string
	}{
		{"set", "default", "set"},
		{"", "default", "default"},
		{"   ", "default", "default"},
		{"set", "", "set"},
	}
	for _, tt := range tests {
		t.Run(tt.value+"|"+tt.def, func(t *testing.T) {
			if got := fallback(tt.value, tt.def); got != tt.want {
				t.Errorf("fallback(%q, %q) = %q, want %q", tt.value, tt.def, got, tt.want)
			}
		})
	}
}

func TestAppendLayoutArg(t *testing.T) {
	base := []string{"--multi"}
	tests := []struct {
		layout string
		extra  bool
	}{
		{"", false},
		{"default", false},
		{"DEFAULT", false},
		{"reverse-list", true},
		{"  reverse-list  ", true},
	}
	for _, tt := range tests {
		t.Run(tt.layout, func(t *testing.T) {
			got := appendLayoutArg(append([]string{}, base...), tt.layout)
			if tt.extra {
				if len(got) != len(base)+2 || got[len(got)-2] != "--layout" {
					t.Errorf("expected --layout to be appended, got %v", got)
				}
			} else {
				if len(got) != len(base) {
					t.Errorf("expected no change, got %v", got)
				}
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := defaultConfig()
	if cfg.FZF.FF.Prompt == "" {
		t.Error("FZF.FF.Prompt should not be empty")
	}
}
