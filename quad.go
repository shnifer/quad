package quad

type QTree struct {
	bound      Rect
	splitCount int
	maxDepth   byte

	points map[V2]struct{}
	nodes  []node
}

func (Q *QTree) PinN(p V2, r node) bool {
	if p.X < r.Min.X || p.X > r.Max.X || (p.X == r.Max.X && p.X != Q.bound.Max.X) {
		return false
	}
	if p.Y < r.Min.Y || p.Y > r.Max.Y || (p.Y == r.Max.Y && p.Y != Q.bound.Max.Y) {
		return false
	}
	return true
}

type node struct {
	Rect
	lvl     byte
	points  Points
	subN    int
	parentN int
}

func NewQTree(bound Rect, splitCount int, maxDepth byte) *QTree {
	return &QTree{
		bound:      bound,
		splitCount: splitCount,
		maxDepth:   maxDepth,
		points:     make(map[V2]struct{}),
		nodes: []node{
			{
				parentN: -1,
				Rect:    bound,
			},
		},
	}
}

func (Q *QTree) Add(points ...V2) {
	for _, point := range points {
		if _, ok := Q.points[point]; !ok {
			Q.points[point] = struct{}{}
			Q.add(point, 0)
		}
	}
}

func (Q *QTree) add(p V2, iNode int) {
	node := Q.nodes[iNode]
	if !Q.PinN(p, node) {
		return
	}
	Q.nodes[iNode].points = append(Q.nodes[iNode].points, p)
	if len(node.points) > Q.splitCount && node.subN == 0 {
		Q.split(iNode)
	} else if node.subN > 0 {
		for part := 0; part < 4; part++ {
			Q.add(p, node.subN+part)
		}
	}
}

func (Q *QTree) split(iNode int) {
	n := Q.nodes[iNode]
	if n.lvl >= Q.maxDepth || n.subN != 0 {
		return
	}
	subN := len(Q.nodes)
	Q.nodes[iNode].subN = subN
	for i := 0; i < 4; i++ {
		Q.nodes = append(Q.nodes, node{
			Rect:    n.SubQuad(i),
			lvl:     n.lvl + 1,
			parentN: iNode,
			points:  nil,
			subN:    0,
		})
	}
	for _, p := range Q.nodes[iNode].points {
		for i := 0; i < 4; i++ {
			Q.add(p, subN+i)
		}
	}
}

func (Q *QTree) Near(p V2) V2 {
	if len(Q.points) == 0 {
		return V2{}
	}
	if _, ok := Q.points[p]; ok {
		return p
	}

	return Q.takeNearest(p, 1)[0]
}

func (Q *QTree) takeNearest(focus V2, num int) Points {
	//drop into most depth node, containing p AND num+ points
	var n node
	spot := 0
	current := 0
	for {
		n = Q.nodes[current]
		if len(n.points) < num {
			break
		}
		spot = current
		if n.subN != 0 {
			for i := n.subN; i < n.subN+4; i++ {
				if Q.PinN(focus, Q.nodes[i]) {
					current = i
				}
			}
		}
		if spot == current {
			break
		}
	}
	n = Q.nodes[spot]
	spotPoints := n.points.GetNClosest(focus, num)
	maxDist := spotPoints[len(spotPoints)-1].To(focus)
	maxRect := Rect{
		Min: V2{X: focus.X - maxDist, Y: focus.Y - maxDist},
		Max: V2{X: focus.X + maxDist, Y: focus.Y + maxDist}}.
		Bound(Q.bound)
	//check if dist radius is totally within spot node, if it is - we found num points we want
	if Q.PinN(maxRect.Min, n) && Q.PinN(maxRect.Max, n) {
		return spotPoints
	}

	empiricCap := len(n.points) * 4
	//rise up to node containing maxRect search area
	for n.parentN >= 0 && !(Q.PinN(maxRect.Min, n) && Q.PinN(maxRect.Max, n)) {
		spot = n.parentN
		n = Q.nodes[spot]
	}

	//now look down for all leaf nodes intersect with search area
	spotPoints = make(Points, 0, empiricCap)
	inds := []int{spot}
	for len(inds) > 0 {
		n = Q.nodes[inds[len(inds)-1]]
		inds = inds[:len(inds)-1]
		if !n.Rect.Intersects(maxRect) {
			continue
		}
		if n.subN == 0 {
			spotPoints = append(spotPoints, n.points...)
		} else {
			inds = append(inds, n.subN, n.subN+1, n.subN+2, n.subN+3)
		}
	}
	spotPoints = spotPoints.GetNClosest(focus, num)

	return spotPoints
}
