package font

import (
	"testing"

	"github.com/eihigh/canvas/internal/resources"
	"github.com/tdewolff/test"
	"golang.org/x/image/font/sfnt"
)

func TestParseTTF(t *testing.T) {
	b, err := resources.FS.ReadFile("DejaVuSerif.ttf")
	test.Error(t, err)

	sfnt, err := ParseFont(b, 0)
	test.Error(t, err)
	test.T(t, sfnt.Head.UnitsPerEm, uint16(2048))
}

func TestParseOTF(t *testing.T) {
	b, err := resources.FS.ReadFile("EBGaramond12-Regular.otf")
	test.Error(t, err)

	sfnt, err := ParseFont(b, 0)
	test.Error(t, err)
	test.T(t, sfnt.Head.UnitsPerEm, uint16(1000))
}

//func TestParseOTF_CFF2(t *testing.T) {
//	b, err := resources.FS.ReadFile("AdobeVFPrototype.otf") // TODO: CFF2
//	test.Error(t, err)
//
//	sfnt, err := ParseFont(b, 0)
//	test.Error(t, err)
//	test.T(t, sfnt.Head.UnitsPerEm, uint16(1000))
//}

func TestParseWOFF(t *testing.T) {
	b, err := resources.FS.ReadFile("DejaVuSerif.woff")
	test.Error(t, err)

	sfnt, err := ParseFont(b, 0)
	test.Error(t, err)
	test.T(t, sfnt.Head.UnitsPerEm, uint16(2048))
}

func TestParseWOFF2(t *testing.T) {
	b, err := resources.FS.ReadFile("DejaVuSerif.woff2")
	test.Error(t, err)

	sfnt, err := ParseFont(b, 0)
	test.Error(t, err)
	test.T(t, sfnt.Head.UnitsPerEm, uint16(2048))
}

func TestParseEOT(t *testing.T) {
	b, err := resources.FS.ReadFile("DejaVuSerif.eot")
	test.Error(t, err)

	sfnt, err := ParseFont(b, 0)
	test.Error(t, err)
	test.T(t, sfnt.Head.UnitsPerEm, uint16(2048))
}

func TestFromGoSFNT(t *testing.T) {
	b, err := resources.FS.ReadFile("DejaVuSerif.ttf")
	test.Error(t, err)

	font, err := sfnt.Parse(b)
	test.Error(t, err)

	buf := FromGoSFNT(font)
	sfnt, err := ParseSFNT(buf, 0)
	test.Error(t, err)
	test.T(t, sfnt.Head.UnitsPerEm, uint16(2048))
}
