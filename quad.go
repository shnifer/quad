package quad

type QTree struct {
	bound      Rect
	splitCount int
	maxDepth   byte

	points map[V2]struct{}
	nodes  []node
}

func (Q *QTree) PinN(p V2, r node) bool {
	return !(p.X < r.Min.X || p.X > r.Max.X || (p.X == r.Max.X && p.X != Q.bound.Max.X) ||
		p.Y < r.Min.Y || p.Y > r.Max.Y || (p.Y == r.Max.Y && p.Y != Q.bound.Max.Y))
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
	newPoints := make(Points, 0, len(points))
	for _, p := range points {

		if Q.PinN(p, Q.nodes[0]) {
			if _, ok := Q.points[p]; !ok {
				Q.points[p] = struct{}{}
				newPoints = append(newPoints, p)
			}
		}
	}
	Q.add(newPoints, 0)
}

func (Q *QTree) add(ps Points, iNode int) {
	node := Q.nodes[iNode]
	Q.nodes[iNode].points = append(Q.nodes[iNode].points, ps...)
	if len(node.points) > Q.splitCount && node.subN == 0 {
		Q.split(iNode)
	} else if node.subN > 0 {
		parts := ps.PartQuads(node.Mid())
		for part := 0; part < 4; part++ {
			Q.add(parts[part], node.subN+part)
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
	parts := Q.nodes[iNode].points.PartQuads(n.Mid())
	for i := 0; i < 4; i++ {
		Q.add(parts[i], subN+i)
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
		if n.subN == 0 {
			break
		}
		quad := n.Mid().PQuad(focus)
		current = n.subN + quad
	}
	n = Q.nodes[spot]
	spotPoints := n.points.GetNClosest(focus, num)
	maxDist := spotPoints[len(spotPoints)-1].To(focus)
	maxRect := Rect{
		Min: V2{X: focus.X - maxDist, Y: focus.Y - maxDist},
		Max: V2{X: focus.X + maxDist, Y: focus.Y + maxDist}}.
		BoundTo(Q.bound)
	//check if dist radius is totally within spot node, if it is - we found num points we want
	if Q.PinN(maxRect.Min, n) && Q.PinN(maxRect.Max, n) {
		return spotPoints
	}

	empiricCap := uint(float64(len(Q.nodes[0].points))*maxRect.Area()/Q.bound.Area()) * 3
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
		if n.subN == 0 || maxRect.Contains(n.Rect) {
			spotPoints = append(spotPoints, n.points...)
		} else {
			inds = append(inds, n.subN, n.subN+1, n.subN+2, n.subN+3)
		}
	}
	spotPoints = spotPoints.GetNClosest(focus, num)

	return spotPoints
}
