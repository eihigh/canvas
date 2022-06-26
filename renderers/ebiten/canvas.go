package ebiten

import (
	"image"
	"image/color"

	"github.com/eihigh/canvas"
	"github.com/eihigh/canvas/renderers/rasterizer"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	white = ebiten.NewImage(1, 1)
	iOp   = &ebiten.DrawImageOptions{}
	tOp   = &ebiten.DrawTrianglesOptions{FillRule: ebiten.EvenOdd}
)

// RenderType is a render type. It affets the speed and quality of rendering.
// See the documents of RenderType constants for details.
type RenderType int

const (
	// The default render type is EvenOddTriangles for Paths and
	// CPURasterizer for Texts.
	_ RenderType = iota

	// EvenOddTriangles converts a path into triangles for passing to DrawTriangles with EvenOdd fill rule.
	// Currently anti-aliasing is not supported with this type.
	EvenOddTriangles

	// CPURasterizer draws with golang.org/x/image/vector.Rasterizer equivalent to 256xAA quality
	// and relatively expensive. Instead of canvas drawing every frame,
	// consider drawing only when necessary and storing the result in a temporary *ebiten.Image.
	CPURasterizer

	// Pathfinder draws with GPU shaders equivalent to 256xAA.
	// See github.com/servo/pathfinder for more details.
	// Pathfinder TODO: implement
)

type RenderOption struct {
	// Type specifies the render type. See the document of RenderType for more details.
	Type RenderType

	// Samples specifies the sample count for anti-aliasing.
	// It is ignored if the Type is CPURasterizer equivalent to 256xAA.
	Samples int
}

// Renderer is a renderer for canvas. Dst image is required.
type Renderer struct {
	Dst  *ebiten.Image
	Path RenderOption
	Text RenderOption
}

func init() {
	white.Fill(color.White)
}

func concatMatrix(g *ebiten.GeoM, m canvas.Matrix) {
	geom := ebiten.GeoM{}
	geom.SetElement(0, 0, m[0][0])
	geom.SetElement(0, 1, m[0][1])
	geom.SetElement(1, 0, m[1][0])
	geom.SetElement(1, 1, m[1][1])
	geom.SetElement(0, 2, m[0][2])
	geom.SetElement(1, 2, m[1][2])
	g.Concat(geom)
}

func New(dst *ebiten.Image) *Renderer {
	return &Renderer{Dst: dst}
}

func (r *Renderer) Size() (float64, float64) {
	w, h := r.Dst.Size()
	return float64(w), float64(h)
}

func (r *Renderer) RenderPath(p *canvas.Path, s canvas.Style, m canvas.Matrix) {
	switch r.Path.Type {
	case CPURasterizer:
		z := rasterizer.FromImage(r.Dst, canvas.DPMM(1), nil)
		z.RenderPath(p, s, m)

	default:
		if s.HasFill() {
			r.renderEvenOddTriangles(p.Transform(m), s.FillColor)
		}
		if s.HasStroke() {
			if s.IsDashed() {
				p = p.Dash(s.DashOffset, s.Dashes...)
			}
			p = p.Stroke(s.StrokeWidth, s.StrokeCapper, s.StrokeJoiner)
			r.renderEvenOddTriangles(p.Transform(m), s.StrokeColor)
		}
	}
}

func (r *Renderer) RenderText(t *canvas.Text, m canvas.Matrix) {
	switch r.Text.Type {
	case EvenOddTriangles:
		rp := *r
		rp.Path.Type = EvenOddTriangles
		t.RenderAsPath(&rp, m, canvas.DPMM(1))

	default:
		rp := *r
		rp.Path.Type = CPURasterizer
		t.WalkDecorations(func(col color.RGBA, p *canvas.Path) {
			style := canvas.DefaultStyle
			style.FillColor = col
			rp.RenderPath(p, style, m)
		})

		_, h := r.Dst.Size()
		rev := canvas.Identity.ReflectYAbout(float64(h) / 2)

		t.WalkSpans(func(x, y float64, span canvas.TextSpan) {
			ppu := span.Face.Size / float64(span.Face.Font.Head.UnitsPerEm) // pixels per unit
			for _, g := range span.Glyphs {
				e, ok := cachedGlyph(span.Face, g.ID)
				if !ok {
					e = makeGlyphCacheEntry(span.Face, g)
					cacheGlyph(e, span.Face, g.ID)
				}
				iOp.GeoM.Reset()
				x0 := x + e.xOffset + float64(span.Face.XOffset+g.XOffset)*ppu
				y0 := y + e.yOffset + float64(span.Face.YOffset+g.YOffset)*ppu
				iOp.GeoM.Translate(x0, y0)
				concatMatrix(&iOp.GeoM, m)
				concatMatrix(&iOp.GeoM, rev)
				iOp.ColorM.Reset()
				iOp.ColorM.ScaleWithColor(span.Face.Color)
				r.Dst.DrawImage(e.image, iOp)
				x += float64(g.XAdvance) * ppu
				y += float64(g.YAdvance) * ppu
			}
		})
	}
}

func (r *Renderer) RenderImage(img image.Image, m canvas.Matrix) {
	e, ok := img.(*ebiten.Image)
	if !ok {
		e = ebiten.NewImageFromImage(img)
	}
	iOp.GeoM.Reset()
	concatMatrix(&iOp.GeoM, m)
	r.Dst.DrawImage(e, iOp)
}

func (r *Renderer) renderEvenOddTriangles(p *canvas.Path, col color.RGBA) {
	_, h := r.Dst.Size()
	p = p.Transform(canvas.Identity.ReflectYAbout(float64(h) / 2))
	vs, is := AppendVerticesAndIncides(p, nil, nil)
	tOp.ColorM.Reset()
	tOp.ColorM.Scale(float64(col.R)/255, float64(col.G)/255, float64(col.B)/255, float64(col.A)/255)
	r.Dst.DrawTriangles(vs, is, white, tOp)
}
