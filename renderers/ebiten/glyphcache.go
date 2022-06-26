package ebiten

import (
	"image"
	"image/color"
	"math"

	"github.com/eihigh/canvas"
	"github.com/eihigh/canvas/font"
	"github.com/eihigh/canvas/text"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/vector"
)

var glyphCache = map[fontFace]map[uint16]*glyphCacheEntry{} // face => glyph ID => image

// fontFace is a key to cache glyph images.
type fontFace struct {
	font                 *canvas.Font
	size                 float64 // in millimeter per em
	style                canvas.FontStyle
	variant              canvas.FontVariant
	fauxBold, fauxItalic float64
	language             string
	script               text.Script
	direction            text.Direction // TODO: really needed here?
}

type glyphCacheEntry struct {
	image            *ebiten.Image
	xOffset, yOffset float64
}

func cachedGlyph(f *canvas.FontFace, glyphID uint16) (e *glyphCacheEntry, ok bool) {
	ff := fontFace{
		f.Font, f.Size, f.Style, f.Variant, f.FauxBold, f.FauxItalic, f.Language, f.Script, f.Direction,
	}
	fc, ok := glyphCache[ff]
	if !ok {
		glyphCache[ff] = map[uint16]*glyphCacheEntry{}
	}
	e, ok = fc[glyphID]
	return
}

func cacheGlyph(e *glyphCacheEntry, f *canvas.FontFace, glyphID uint16) {
	ff := fontFace{
		f.Font, f.Size, f.Style, f.Variant, f.FauxBold, f.FauxItalic, f.Language, f.Script, f.Direction,
	}
	fc, ok := glyphCache[ff]
	if !ok {
		glyphCache[ff] = map[uint16]*glyphCacheEntry{}
	}
	fc[glyphID] = e
}

func makeGlyphCacheEntry(f *canvas.FontFace, g text.Glyph) *glyphCacheEntry {
	// Get glyph path
	ppu := f.Size / float64(f.Font.Head.UnitsPerEm) // pixels per unit
	p := &canvas.Path{}
	f.Font.GlyphPath(p, g.ID, f.PPEM(canvas.DPMM(1)), 0, 0, ppu, font.NoHinting)
	if f.FauxBold != 0.0 {
		p = p.Offset(f.FauxBold*f.Size, canvas.NonZero)
	}
	if f.FauxItalic != 0.0 {
		p = p.Transform(canvas.Identity.Shear(f.FauxItalic, 0.0))
	}

	// Get glyph bounds
	var xmin, xmax, ymin, ymax float64
	s := p.Scanner()
	for s.Scan() {
		d := s.Values()
		switch s.Cmd() {
		case canvas.MoveToCmd:
			xmin = math.Min(xmin, d[0])
			xmax = math.Max(xmax, d[0])
			ymin = math.Min(ymin, d[1])
			ymax = math.Max(ymax, d[1])
		case canvas.LineToCmd:
			xmin = math.Min(xmin, d[0])
			xmax = math.Max(xmax, d[0])
			ymin = math.Min(ymin, d[1])
			ymax = math.Max(ymax, d[1])
		case canvas.QuadToCmd:
			xmin = math.Min(xmin, d[0])
			xmax = math.Max(xmax, d[0])
			ymin = math.Min(ymin, d[1])
			ymax = math.Max(ymax, d[1])
			xmin = math.Min(xmin, d[2])
			xmax = math.Max(xmax, d[2])
			ymin = math.Min(ymin, d[3])
			ymax = math.Max(ymax, d[3])
		case canvas.CubeToCmd:
			xmin = math.Min(xmin, d[0])
			xmax = math.Max(xmax, d[0])
			ymin = math.Min(ymin, d[1])
			ymax = math.Max(ymax, d[1])
			xmin = math.Min(xmin, d[2])
			xmax = math.Max(xmax, d[2])
			ymin = math.Min(ymin, d[3])
			ymax = math.Max(ymax, d[3])
			xmin = math.Min(xmin, d[4])
			xmax = math.Max(xmax, d[4])
			ymin = math.Min(ymin, d[5])
			ymax = math.Max(ymax, d[5])
		case canvas.CloseCmd:
		}
	}
	xmin = math.Floor(xmin)
	ymin = math.Floor(ymin)
	w := int(math.Ceil(xmax - xmin))
	h := int(math.Ceil(ymax - ymin))
	img := ebiten.NewImage(w, h)

	// Rasterize glyph image
	p = p.Translate(-xmin, -ymin)
	z := vector.NewRasterizer(w, h)
	p.ToInverseRasterizer(z, canvas.DPMM(1))
	z.Draw(img, image.Rect(0, 0, w, h), image.NewUniform(color.White), image.Point{})
	return &glyphCacheEntry{img, xmin, ymin}
}
