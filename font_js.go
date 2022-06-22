//go:build js

package canvas

// FindLocalFont finds the path to a font from the system's fonts.
func FindLocalFont(name string, style FontStyle) (string, error) {
	panic("FindLocalFont is not supported in js")
}
