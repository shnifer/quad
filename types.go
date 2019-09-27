package quad

import (
	"fmt"
	"math"
	"sort"
)

type V2 struct {
	X, Y float64
}

func (p1 V2) To(p2 V2) float64 {
	dx := p1.X - p2.X
	dy := p1.Y - p2.Y
	return math.Sqrt(dx*dx + dy*dy)
}

type Points []V2

func (ps Points) GetNClosest(center V2, n int) Points {
	if n == 1 {
		if len(ps) == 0 {
			return nil
		}
		bestDist := ps[0].To(center)
		bestInd := 0
		for i := range ps {
			d := ps[i].To(center)
			if d < bestDist {
				bestDist = d
				bestInd = i
			}
		}
		return ps[bestInd : bestInd+1]
	}
	ps.SortAround(center)
	if len(ps) > n {
		ps = ps[:n]
	}
	return ps
}

type t struct {
	dists  []float64
	points Points
}

func (t t) Len() int {
	return len(t.dists)
}

func (t t) Less(i, j int) bool {
	return t.dists[i] < t.dists[j]
}

func (t t) Swap(i, j int) {
	t.dists[i], t.dists[j] = t.dists[j], t.dists[i]
	t.points[i], t.points[j] = t.points[j], t.points[i]
}

func (ps Points) SortAround(center V2) {
	t := t{
		dists:  make([]float64, len(ps)),
		points: ps,
	}
	for i := range ps {
		t.dists[i] = ps[i].To(center)
	}
	sort.Sort(t)
}

func (ps Points) MaxDist(to V2) float64 {
	maxDist := 0.0
	for _, p := range ps {
		dist := p.To(to)
		if dist > maxDist {
			maxDist = dist
		}
	}
	return maxDist
}

func (ps Points) Bounds() Rect {
	if len(ps) == 0 {
		return Rect{}
	}
	r := Rect{V2{X: ps[0].X, Y: ps[0].Y}, V2{X: ps[0].X, Y: ps[0].Y}}
	for i := 1; i < len(ps); i++ {
		if r.Min.X > ps[i].X {
			r.Min.X = ps[i].X
		}
		if r.Min.Y > ps[i].Y {
			r.Min.Y = ps[i].Y
		}
		if r.Max.X < ps[i].X {
			r.Max.X = ps[i].X
		}
		if r.Max.Y < ps[i].X {
			r.Max.Y = ps[i].X
		}
	}
	return r
}

type Rect struct {
	Min, Max V2
}

func (r Rect) SubQuad(i int) Rect {
	mid := V2{
		X: (r.Min.X + r.Max.X) / 2,
		Y: (r.Min.Y + r.Max.Y) / 2,
	}
	switch i {
	//nw
	case 0:
		return Rect{
			Min: V2{X: r.Min.X, Y: mid.Y},
			Max: V2{X: mid.X, Y: r.Max.Y},
		}
	//ne
	case 1:
		return Rect{
			Min: V2{X: mid.X, Y: mid.Y},
			Max: V2{X: r.Max.X, Y: r.Max.Y},
		}
	//se
	case 2:
		return Rect{
			Min: V2{X: mid.X, Y: r.Min.Y},
			Max: V2{X: r.Max.X, Y: mid.Y},
		}
	//sw
	case 3:
		return Rect{
			Min: V2{X: r.Min.X, Y: r.Min.Y},
			Max: V2{X: mid.X, Y: mid.Y},
		}
	default:
		return Rect{}
	}
}

func (r Rect) Bound(b Rect) Rect {
	if r.Min.X < b.Min.X {
		r.Min.X = b.Min.X
	}
	if r.Min.Y < b.Min.Y {
		r.Min.Y = b.Min.Y
	}
	if r.Max.X > b.Max.X {
		r.Max.X = b.Max.X
	}
	if r.Max.Y > b.Max.Y {
		r.Max.Y = b.Max.Y
	}
	return r
}

func (r Rect) Intersects(s Rect) bool {
	return r.Min.X <= s.Max.X && s.Min.X <= r.Max.X &&
		r.Min.Y <= s.Max.Y && s.Min.Y <= r.Max.Y
}

func (v V2) String() string {
	return fmt.Sprintf("(%v, %v)", v.X, v.Y)
}
