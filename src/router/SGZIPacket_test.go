package main

import (
	"testing"
)

func Test_sgzipChecksum(t *testing.T) {
	tests := []struct {
		name string
		args []byte
		want int16
	}{
		{"basic", []byte{1, 1, 1, 1, 1}, 5},
		{"basic2", []byte{1, 1, 1, 2, 1}, 6},
		{"zero", []byte{0, 0, 0, 0, 0, 0, 0, 0}, 0},
		{"big", []byte{255, 255, 255, 255, 255}, 1275},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sgzipChecksum(tt.args); got != tt.want {
				t.Errorf("sgzipChecksum() = %v, want %v", got, tt.want)
			}
		})
	}
}
