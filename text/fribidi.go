//go:build !fribidi || js
// +build !fribidi js

package text

// Bidi is not supported.
func Bidi(text string) (string, []int) {
	runes := []rune(text)
	mapV2L := make([]int, len(runes))
	for i := range mapV2L {
		mapV2L[i] = i
	}
	return text, mapV2L
}
