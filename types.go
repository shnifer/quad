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
		l := len(ps)
		if l == 0 {
			return nil
		}
		bestDist := ps[0].To(center)
		bestPoint := Points{ps[0]}
		for i := 1; i < l; i++ {
			d := ps[i].To(center)
			if d < bestDist {
				bestDist = d
				bestPoint[0] = ps[i]
			}
		}
		return bestPoint
	}
	res := make(Points, len(ps))
	copy(res, ps)
	res.SortAround(center)
	if len(res) > n {
		res = res[:n]
	}
	return res
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

func (ps Points) PartQuads(center V2) (parts [4]Points) {
	sort.Slice(ps, func(i, j int) bool {
		return center.PQuad(ps[i]) < center.PQuad(ps[j])
	})
	ind1 := sort.Search(len(ps), func(i int) bool {
		return center.PQuad(ps[i]) > 0
	})
	ind2 := sort.Search(len(ps), func(i int) bool {
		return center.PQuad(ps[i]) > 1
	})
	ind3 := sort.Search(len(ps), func(i int) bool {
		return center.PQuad(ps[i]) > 2
	})
	parts[0] = ps[:ind1]
	parts[1] = ps[ind1:ind2]
	parts[2] = ps[ind2:ind3]
	parts[3] = ps[ind3:]
	return parts
}

func (c V2) PQuad(p V2) (q int) {
	if p.X >= c.X {
		q += 1
	}
	if p.Y >= c.Y {
		q += 2
	}
	return q
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

func (r Rect) Mid() V2 {
	return V2{
		X: (r.Min.X + r.Max.X) / 2,
		Y: (r.Min.Y + r.Max.Y) / 2,
	}
}

func (r Rect) SubQuad(i int) Rect {
	mid := r.Mid()
	switch i {
	//sw
	case 0:
		return Rect{
			Min: V2{X: r.Min.X, Y: r.Min.Y},
			Max: V2{X: mid.X, Y: mid.Y},
		}
	//se
	case 1:
		return Rect{
			Min: V2{X: mid.X, Y: r.Min.Y},
			Max: V2{X: r.Max.X, Y: mid.Y},
		}
	//nw
	case 2:
		return Rect{
			Min: V2{X: r.Min.X, Y: mid.Y},
			Max: V2{X: mid.X, Y: r.Max.Y},
		}
	//ne
	case 3:
		return Rect{
			Min: V2{X: mid.X, Y: mid.Y},
			Max: V2{X: r.Max.X, Y: r.Max.Y},
		}
	default:
		return Rect{}
	}
}

func (r Rect) BoundTo(b Rect) Rect {
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

func (r Rect) Contains(s Rect) bool {
	return s.Min.X >= r.Min.X && s.Max.X <= r.Max.X &&
		s.Min.Y >= r.Min.Y && s.Max.Y <= r.Max.Y
}

func (r Rect) Area() float64 {
	w := r.Max.X - r.Min.X
	h := r.Max.Y - r.Min.Y
	return w * h
}

func (v V2) String() string {
	return fmt.Sprintf("(%v, %v)", v.X, v.Y)
}
