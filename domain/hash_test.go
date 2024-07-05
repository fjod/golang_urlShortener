package domain

import (
	"testing"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		name string
		id   int
		want string
	}{
		{"zero", 0, ""},
		{"one", 1, "b"},
		{"base", 62, "ba"},
		{"max", 3844, "baa"},
		{"max", 19158, "e9a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := encode(tt.id); got != tt.want {
				t.Errorf("encode(%d) = %q, want %q", tt.id, got, tt.want)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int
	}{
		{"empty", "", 0},
		{"single", "a", 0},
		{"multiple", "ba", 62},
		{"max", "e9a", 19158},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := decode(tt.s); got != tt.want {
				t.Errorf("decode(%q) = %d, want %d", tt.s, got, tt.want)
			}
		})
	}
}
