package quad

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
)

func TestNewQTree(t *testing.T) {
	Q := NewQTree(Rect{V2{0, 0}, V2{10, 10}}, 2, 10)
	Q.Add(V2{1, 1}, V2{2, 2}, V2{3, 3})

	got := Q.Near(V2{0, 0})
	require.Equal(t, V2{1, 1}, got)
	got = Q.Near(V2{10, 10})
	require.Equal(t, V2{3, 3}, got)
	got = Q.Near(V2{1.9, 2.1})
	require.Equal(t, V2{2, 2}, got)
}

func TestFuzz(t *testing.T) {
	Q := NewQTree(Rect{V2{0, 0}, V2{100, 100}}, 4, 10)
	points := make(Points, 1000)
	for n := 0; n < 1000; n++ {
		points[n] = randP(100)
		Q.Add(points[n])

		//check duplicate points
		for _, node := range Q.nodes {
			m := make(map[V2]struct{})
			for _, p := range node.points {
				_, duplicate := m[p]
				require.False(t, duplicate)
				m[p] = struct{}{}
			}
		}
	}
	require.Len(t, Q.points, 1000)
	require.Len(t, Q.nodes[0].points, 1000)

	//check points near itself
	for _, p := range points {
		got := Q.Near(p)
		require.Equal(t, got, p)
	}

	//check Near()
	for i := 0; i < 1000; i++ {
		focus := randP(200)
		focus.X -= 50
		focus.Y -= 50

		got := Q.Near(focus)
		bestDist := 1000000.0
		bestP := V2{}
		for j := 0; j < 1000; j++ {
			dist := points[j].To(focus)
			if dist < bestDist {
				bestDist = dist
				bestP = points[j]
			}
		}
		require.Equal(t, got, bestP)
	}
}

func BenchmarkQTree_Near(b *testing.B) {
	params := []struct {
		split int
		depth byte
	}{
		{1, 2},
		{1, 3},
		{1, 4},
		//{1, 6},
		//{1, 8},
		{2, 3},
		//{2, 4},
		//{2, 6},
		//{2, 8},
		{3, 3},
		//{3, 4},
		//{3, 6},
		//{3, 8},
		{4, 3},
		//{4, 4},
		//{4, 6},
		//{4, 8},
		{5, 3},
		{6, 3},
		//{7, 2},
		{7, 3},
		//{7, 4},
		//{8, 2},
		{8, 3},
		//{8, 4},
		//{8, 6},
		//{8, 8},
		//{10, 2},
		//{10, 3},
		//{10, 4},
		//{12, 2},
		//{12, 3},
		//{12, 4},
		//{12, 6},
		//{12, 8},
	}
	for _, param := range params {
		split, depth := param.split, param.depth
		Q := NewQTree(Rect{V2{0, 0}, V2{100, 100}}, split, depth)
		points := make(Points, 1000)
		for n := 0; n < 1000; n++ {
			points[n] = randP(100)
			Q.Add(randP(100))
		}
		s := 0
		for _, n := range Q.nodes {
			s += len(n.points)
		}

		b.Run(fmt.Sprintf("with split %d and depth %d, approx %d kb", split, depth, (len(Q.nodes)*80+s*16)/1000), func(b *testing.B) {

			for i := 0; i < b.N/1000; i++ {
				for j := 0; j < 1000; j++ {
					_ = Q.Near(points[j])
				}
			}
		})
	}
}

func BenchmarkQTree_NaiveNear(b *testing.B) {
	points := make(Points, 1000)
	focuses := make(Points, 1000)
	for n := 0; n < 1000; n++ {
		points[n] = randP(100)
		focuses[n] = randP(100)
	}
	b.ResetTimer()

	for i := 0; i < b.N/1000; i++ {
		for _, focus := range focuses {
			bestDist := points[0].To(focus)
			bestP := points[0]
			for _, p := range points {
				d := p.To(focus)
				if d < bestDist {
					bestDist = d
					bestP = p
				}
			}
			_ = bestP
		}
	}
}

func randP(s float64) V2 {
	return V2{X: rand.Float64() * s, Y: rand.Float64() * s}
}
