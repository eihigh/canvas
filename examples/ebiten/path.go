package main

import (
	"image/color"

	"github.com/eihigh/canvas"
	renderer "github.com/eihigh/canvas/renderers/ebiten"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	white = ebiten.NewImage(1, 1)
)

func init() {
	white.Fill(color.White)
}

func (g *game) drawPath(screen *ebiten.Image) {
	p := &canvas.Path{}
	p.MoveTo(150, 150)
	p.LineTo(250, 350)
	p.LineTo(50, 350)
	p.Close()
	vs, is := renderer.AppendVerticesAndIncides(p, nil, nil)
	op := &ebiten.DrawTrianglesOptions{FillRule: ebiten.EvenOdd}
	screen.DrawTriangles(vs, is, white, op)

	p = canvas.Rectangle(800, 100)
	vs, is = renderer.AppendVerticesAndIncides(p, nil, nil)
	for i, v := range vs {
		if v.DstX < 400 {
			vs[i].ColorG = 0
		} else {
			vs[i].ColorB = 0
		}
	}
	screen.DrawTriangles(vs, is, white, op)

	polyline := &canvas.Polyline{}
	polyline.Add(0.0, 0.0)
	polyline.Add(120.0, 0.0)
	polyline.Add(120.0, 45.0)
	polyline.Add(0.0, 120.0)
	polyline.Add(0.0, 0.0)
	p = polyline.ToPath().Translate(500, 200)
	vs, is = renderer.AppendVerticesAndIncides(p, nil, nil)
	op.ColorM.ScaleWithColor(canvas.Seagreen)
	screen.DrawTriangles(vs, is, white, op)
}
