package main

import "testing"

func TestResolveEditor(t *testing.T) {
	tests := []struct {
		name       string
		onixEditor string
		editor     string
		want       string
	}{
		{"ONIX_EDITOR takes priority", "code", "vim", "code"},
		{"EDITOR fallback", "", "vim", "vim"},
		{"default nvim", "", "", "nvim"},
		{"ONIX_EDITOR beats EDITOR", "code", "vim", "code"},
		{"whitespace-only ONIX_EDITOR falls through", "  ", "vim", "vim"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("ONIX_EDITOR", tt.onixEditor)
			t.Setenv("EDITOR", tt.editor)
			if got := resolveEditor(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
