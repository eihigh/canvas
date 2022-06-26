package ebiten

import (
	"math"

	"github.com/eihigh/canvas"
	"github.com/hajimehoshi/ebiten/v2"
)

// TODO: use functions in path_util.go

// triangulator converts a path into triangles (vertices + indices).
type triangulator struct {
	m, n    int
	ax, ay  float32
	vs      []ebiten.Vertex
	indices []uint16
}

func (T *triangulator) moveTo(x, y float32) {
	T.ax, T.ay = x, y
	T.vs = append(T.vs, ebiten.Vertex{
		SrcX: 0, SrcY: 0, DstX: x, DstY: y,
		ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1,
	})
	T.n++
	T.m = T.n
}

func (T *triangulator) closePath() {
	T.m = T.n + 1
}

func (T *triangulator) lineTo(x, y float32) {
	T.ax, T.ay = x, y
	T.vs = append(T.vs, ebiten.Vertex{
		SrcX: 0, SrcY: 0, DstX: x, DstY: y,
		ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1,
	})
	T.n++
	// Append indices after the 3rd vertex
	if (T.n - T.m) >= 2 {
		T.indices = append(T.indices, uint16(T.m), uint16(T.n-1), uint16(T.n))
	}
}

func (T *triangulator) quadTo(bx, by, cx, cy float32) {
	devsq := devSquared(T.ax, T.ay, bx, by, cx, cy)
	if devsq >= 0.333 {
		const tol = 3
		n := int(math.Sqrt(math.Sqrt(tol * float64(devsq))))
		t, nInv := float32(0), 1/float32(n)
		for i := 0; i < n-1; i++ {
			t += nInv
			dx := lerp(t, lerp(t, T.ax, bx), lerp(t, bx, cx))
			dy := lerp(t, lerp(t, T.ay, by), lerp(t, by, cy))
			T.lineTo(dx, dy)
		}
	}
	T.lineTo(cx, cy)
}

func (T *triangulator) cubeTo(bx, by, cx, cy, dx, dy float32) {
	devsq := devSquared(T.ax, T.ay, bx, by, cx, cy)
	if devsqAlt := devSquared(T.ax, T.ay, cx, cy, dx, cy); devsq < devsqAlt {
		devsq = devsqAlt
	}
	if devsq >= 0.333 {
		const tol = 3
		n := int(math.Sqrt(math.Sqrt(tol * float64(devsq))))
		t, nInv := float32(0), 1/float32(n)
		for i := 0; i < n-1; i++ {
			t += nInv
			abx, aby := lerp(t, T.ax, bx), lerp(t, T.ay, by)
			bcx, bcy := lerp(t, bx, cx), lerp(t, by, cy)
			cdx, cdy := lerp(t, cx, dx), lerp(t, cy, dy)
			abcx, abcy := lerp(t, abx, bcx), lerp(t, aby, bcy)
			bcdx, bcdy := lerp(t, bcx, cdx), lerp(t, bcy, cdy)
			T.lineTo(lerp(t, abcx, bcdx), lerp(t, abcy, bcdy))
		}
	}
	T.lineTo(dx, dy)
}

func devSquared(ax, ay, bx, by, cx, cy float32) float32 {
	devx := ax - 2*bx + cx
	devy := ay - 2*by + cy
	return devx*devx + devy*devy
}

func lerp(t, a, b float32) float32 {
	return a + t*(b-a)
}

func AppendVerticesAndIncides(p *canvas.Path, vs []ebiten.Vertex, indices []uint16) ([]ebiten.Vertex, []uint16) {
	p = p.ReplaceArcs()
	t := &triangulator{m: len(vs), n: len(vs) - 1, vs: vs, indices: indices}
	s := p.Scanner()
	for s.Scan() {
		d := s.Values()
		switch s.Cmd() {
		case canvas.MoveToCmd:
			t.moveTo(float32(d[0]), float32(d[1]))
		case canvas.LineToCmd:
			t.lineTo(float32(d[0]), float32(d[1]))
		case canvas.QuadToCmd:
			t.quadTo(float32(d[0]), float32(d[1]), float32(d[2]), float32(d[3]))
		case canvas.CubeToCmd:
			t.cubeTo(float32(d[0]), float32(d[1]), float32(d[2]), float32(d[3]), float32(d[4]), float32(d[5]))
		case canvas.CloseCmd:
			t.closePath()
		}
	}
	return t.vs, t.indices
}
