//go:build !js

package canvas

import "github.com/flopp/go-findfont"

// FindLocalFont finds the path to a font from the system's fonts.
func FindLocalFont(name string, style FontStyle) (string, error) {
	// TODO: use style to match font
	return findfont.Find(name)
}
