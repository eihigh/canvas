// +build gofuzz

package fuzz

import "github.com/eihigh/canvas/font"

// Fuzz is a fuzz test.
func Fuzz(data []byte) int {
	_, _ = font.ParseWOFF2(data)
	return 1
}
