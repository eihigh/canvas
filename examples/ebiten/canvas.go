package main

import (
	"fmt"
	"image/color"

	"github.com/eihigh/canvas"
	"github.com/eihigh/canvas/text"
)

func (g *game) drawCanvas(c *canvas.Context) {
	// Draw a comprehensive text box
	pt := 56.0 * canvas.PtToMM
	face := g.fontLatin.Face(pt, canvas.Black, canvas.FontRegular, canvas.FontNormal)

	rt := canvas.NewRichText(face)
	rt.Add(face, "Lorem dolor ipsum ")
	rt.Add(g.fontLatin.Face(pt, canvas.White, canvas.FontBold, canvas.FontNormal), "confiscator")
	rt.Add(face, " cur\u200babitur ")
	rt.Add(g.fontLatin.Face(pt, canvas.Black, canvas.FontItalic, canvas.FontNormal), "mattis")
	rt.Add(face, " dui ")
	rt.Add(g.fontLatin.Face(pt, canvas.Black, canvas.FontBold|canvas.FontItalic, canvas.FontNormal), "tellus")
	rt.Add(face, " vel. Proin ")
	rt.Add(g.fontLatin.Face(pt, canvas.Black, canvas.FontRegular, canvas.FontNormal, canvas.FontUnderline), "sodales")
	rt.Add(face, " eros vel ")
	rt.Add(g.fontLatin.Face(pt, canvas.Black, canvas.FontRegular, canvas.FontNormal, canvas.FontSineUnderline), "nibh")
	rt.Add(face, " fringilla pellen\u200btesque eu cillum. ")

	face = g.fontLatin.Face(pt, canvas.Black, canvas.FontRegular, canvas.FontNormal)
	face.Language = "ru"
	face.Script = text.Cyrillic
	rt.Add(face, "дёжжэнтиюнт холст ")

	face = g.fontDevanagari.Face(pt, canvas.Black, canvas.FontRegular, canvas.FontNormal)
	face.Language = "hi"
	face.Script = text.Devanagari
	rt.Add(face, "हालाँकि प्र ")

	g.drawText(c, 20, 20, face, rt)

	// Draw the word Stroke being stroked
	face = g.fontLatin.Face(320.0*canvas.PtToMM, canvas.Black, canvas.FontRegular, canvas.FontNormal)
	p, _, _ := face.ToPath("Stroke")
	p = p.Scale(1, -1)
	c.DrawPath(400, 395, p.Stroke(3, canvas.RoundCap, canvas.RoundJoin))

	// Draw an elliptic arc being dashed
	ellipse, err := canvas.ParseSVG(fmt.Sprintf("A40 120 30 1 0 120 0z"))
	if err != nil {
		panic(err)
	}
	ellipse = ellipse.Scale(1, -1)
	c.SetFillColor(canvas.Whitesmoke)
	c.DrawPath(440, 240, ellipse)

	c.SetFillColor(canvas.Transparent)
	c.SetStrokeColor(canvas.Black)
	c.SetStrokeWidth(3)
	c.SetStrokeCapper(canvas.RoundCap)
	c.SetStrokeJoiner(canvas.RoundJoin)
	c.SetDashes(0.0, 6.0, 12.0, 6.0, 6.0, 12.0, 6.0)
	//ellipse = ellipse.Dash(0.0, 2.0, 4.0, 2.0).Stroke(0.5, canvas.RoundCap, canvas.RoundJoin)
	c.DrawPath(440, 240, ellipse)
	c.SetStrokeColor(canvas.Transparent)
	c.SetDashes(0.0)

	// Draw a LaTeX formula
	// latex, err := canvas.ParseLaTeX(`$y = \sin(\frac{x}{180}\pi)$`)
	// if err != nil {
	// 	panic(err)
	// }
	// latex = latex.Transform(canvas.Identity.Rotate(-30))
	// c.SetFillColor(canvas.Black)
	// c.DrawPath(135, 85, latex)
	//
	// // Draw a raster image
	// lenna, err := resources.FS.ReadFile("lenna.png")
	// if err != nil {
	// 	panic(err)
	// }
	// img, err := canvas.NewPNGImage(bytes.NewBuffer(lenna))
	// if err != nil {
	// 	panic(err)
	// }
	// c.Rotate(5)
	// c.DrawImage(50.0, 0.0, img, 15)
	// c.SetView(canvas.Identity)

	// Draw an closed set of points being smoothed
	polyline := &canvas.Polyline{}
	polyline.Add(0.0, 0.0)
	polyline.Add(120.0, 0.0)
	polyline.Add(120.0, 45.0)
	polyline.Add(0.0, 120.0)
	polyline.Add(0.0, 0.0)
	c.SetFillColor(canvas.Seagreen)
	c.FillColor.R = byte(float64(c.FillColor.R) * 0.25)
	c.FillColor.G = byte(float64(c.FillColor.G) * 0.25)
	c.FillColor.B = byte(float64(c.FillColor.B) * 0.25)
	c.FillColor.A = byte(float64(c.FillColor.A) * 0.25)
	c.SetStrokeColor(canvas.Seagreen)
	c.DrawPath(620, 260, polyline.Smoothen().Scale(1, -1))

	c.SetFillColor(canvas.Transparent)
	c.SetStrokeColor(canvas.Black)
	c.SetStrokeWidth(2)
	c.DrawPath(620, 260, polyline.ToPath().Scale(1, -1))
	c.SetStrokeWidth(3)

	// Draw a open set of points being smoothed
	polyline = &canvas.Polyline{}
	polyline.Add(0.0, 0.0)
	polyline.Add(80.0, 40.0)
	polyline.Add(160.0, 120.0)
	polyline.Add(240.0, 160.0)
	polyline.Add(320.0, 80.0)
	c.SetStrokeColor(canvas.Dodgerblue)
	c.DrawPath(40, 340, polyline.Smoothen().Scale(1, -1))
	c.SetStrokeColor(canvas.Black)
}

func (g *game) drawText(c *canvas.Context, x, y float64, face *canvas.FontFace, rich *canvas.RichText) {
	metrics := face.Metrics()
	width, height := 360.0, 128.0

	text := rich.ToText(width, height, canvas.Justify, canvas.Top, 0.0, 0.0)

	c.SetFillColor(color.RGBA{192, 0, 64, 255})
	c.DrawPath(x, y, text.Bounds().ToPath().Scale(1, -1))
	c.SetFillColor(color.RGBA{51, 51, 51, 51})
	c.DrawPath(x, y, canvas.Rectangle(width, -metrics.LineHeight).Scale(1, -1))
	c.SetFillColor(color.RGBA{0, 0, 0, 51})
	c.DrawPath(x, y+metrics.CapHeight-metrics.Ascent, canvas.Rectangle(width, -metrics.CapHeight-metrics.Descent).Scale(1, -1))
	c.DrawPath(x, y+metrics.XHeight-metrics.Ascent, canvas.Rectangle(width, -metrics.XHeight).Scale(1, -1))

	c.SetFillColor(canvas.Black)
	c.DrawPath(x, y, canvas.Rectangle(width, -height).Stroke(0.8, canvas.RoundCap, canvas.RoundJoin).Scale(1, -1))
	c.DrawText(x, y, text)
}
