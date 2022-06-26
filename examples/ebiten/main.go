package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/eihigh/canvas"
	"github.com/eihigh/canvas/internal/resources"
	renderer "github.com/eihigh/canvas/renderers/ebiten"
)

type drawMode int

const (
	pathMode drawMode = iota
	canvasMode
	canvas2Mode
)

const (
	vw, vh = 800, 400
)

type game struct {
	mode           drawMode
	fontLatin      *canvas.FontFamily
	fontArabic     *canvas.FontFamily
	fontDevanagari *canvas.FontFamily
}

func newGame() (*game, error) {
	g := &game{}
	g.fontLatin = canvas.NewFontFamily("DejaVu Serif")
	if err := g.fontLatin.LoadFontFileFS(resources.FS, "DejaVuSerif.ttf", canvas.FontRegular); err != nil {
		return nil, err
	}
	g.fontArabic = canvas.NewFontFamily("DejaVu Sans")
	if err := g.fontArabic.LoadFontFileFS(resources.FS, "DejaVuSans.ttf", canvas.FontRegular); err != nil {
		return nil, err
	}
	g.fontDevanagari = canvas.NewFontFamily("Noto Serif")
	if err := g.fontDevanagari.LoadFontFileFS(resources.FS, "NotoSerifDevanagari.ttf", canvas.FontRegular); err != nil {
		return nil, err
	}
	return g, nil
}

func (g *game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		switch g.mode {
		case pathMode:
			g.mode = canvasMode
		case canvasMode:
			g.mode = canvas2Mode
		case canvas2Mode:
			g.mode = pathMode
		}
	}
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	switch g.mode {
	case pathMode:
		g.drawPath(screen)
		ebitenutil.DebugPrint(screen, "[Path Example]\nPress space to toggle examples")

	case canvasMode:
		screen.Fill(color.White)
		r := renderer.New(screen)
		ctx := canvas.NewContext(r)
		g.drawCanvas(ctx)
		ebitenutil.DebugPrint(screen, "[Default Canvas Example]\nPress space to toggle examples")

	case canvas2Mode:
		screen.Fill(color.White)
		r := renderer.New(screen)
		r.Path.Type = renderer.CPURasterizer
		r.Text.Type = renderer.EvenOddTriangles
		ctx := canvas.NewContext(r)
		g.drawCanvas(ctx)
		ebitenutil.DebugPrint(screen, "[Another Canvas Example]\nPress space to toggle examples")
	}
}

func (g *game) Layout(w, h int) (int, int) {
	return vw, vh
}

func main() {
	ebiten.SetWindowSize(vw, vh)
	g, err := newGame()
	if err != nil {
		panic(err)
	}
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
