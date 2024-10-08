package tdqt

type Tree struct {
	area            Rectangle
	cannotSubdivide bool
	depth           uint8
	insertCallback  func(tree *Tree, depth uint8)
	maxObjects      uint16
	objects         map[uint64]Object
	subTrees        [4]*Tree
}

func (t *Tree) Insert(obj Object) {
	h := obj.Hash()
	t.insert(h, obj, 0)
}

func (t *Tree) Search(area Rectangle) map[uint64]Object {
	result := make(map[uint64]Object)

	t.search(area, result)

	return result
}

func (t *Tree) search(area Rectangle, result map[uint64]Object) {
	if !t.overlaps(area) {
		return
	}

	for _, st := range t.subTrees {
		if st == nil {
			break // any nil subTree means we won't find subsequent subTrees
		}

		st.search(area, result)
	}

	for k, v := range t.objects {
		if overlap, _ := v.Overlaps(area); overlap {
			result[k] = v
		}
	}
}

func (t *Tree) createSubtrees() {
	xMid := t.area.xRange.midpoint()
	yMid := t.area.yRange.midpoint()

	var cannotSplitX bool
	if xMid == t.area.xRange.max {
		cannotSplitX = true
	}

	var cannotSplitY bool
	if yMid == t.area.yRange.max {
		cannotSplitY = true
	}

	var subTreeAreas []Rectangle

	// Calculate the Limits of each subtree. We don't need to worry about
	// being unable to split both X and Y. This is handled by the tree's
	// cannotSubdivide boolean.
	switch {
	case cannotSplitX:
		subTreeAreas = []Rectangle{
			NewRectangle(NewLimits(t.area.xRange.min, t.area.xRange.max), NewLimits(yMid, t.area.yRange.max)), // quadrant I and II
			NewRectangle(NewLimits(t.area.xRange.min, t.area.xRange.max), NewLimits(t.area.yRange.min, yMid)), // quadrant III and IV
		}
	case cannotSplitY:
		subTreeAreas = []Rectangle{ // Calculate the Limits of each subtree
			NewRectangle(NewLimits(xMid, t.area.xRange.max), NewLimits(t.area.yRange.min, t.area.yRange.max)), // quadrant I and IV
			NewRectangle(NewLimits(t.area.xRange.min, xMid), NewLimits(t.area.yRange.min, t.area.yRange.max)), // quadrant I and IV
		}
	default:
		subTreeAreas = []Rectangle{ // Calculate the Limits of each subtree
			NewRectangle(NewLimits(xMid, t.area.xRange.max), NewLimits(yMid, t.area.yRange.max)), // quadrant I
			NewRectangle(NewLimits(t.area.xRange.min, xMid), NewLimits(yMid, t.area.yRange.max)), // quadrant II
			NewRectangle(NewLimits(t.area.xRange.min, xMid), NewLimits(t.area.yRange.min, yMid)), // quadrant III
			NewRectangle(NewLimits(xMid, t.area.xRange.max), NewLimits(t.area.yRange.min, yMid)), // quadrant IV
		}
	}

	// Create each subtree using the calculated Limits
	for i, sta := range subTreeAreas {
		xMin, xMax, yMin, yMax := sta.xyMinMax()
		t.subTrees[i] = newTree(newTreeCfg{
			xMin:               xMin,
			xMax:               xMax,
			yMin:               yMin,
			yMax:               yMax,
			insertCallbackFunc: t.insertCallback,
			maxObjects:         t.maxObjects,
		})
	}
}

func (t *Tree) SetInsertCallback(f func(tree *Tree, depth uint8)) {
	t.insertCallback = f
}

func (t *Tree) insert(key uint64, obj Object, depth uint8) {
	if t.cannotSubdivide {
		// This node has been divided as far as we're going to take it.
		// Add this point without regard for the usual capacity limit.
		t.objects[key] = obj
		if t.insertCallback != nil {
			t.insertCallback(t, depth)
		}
		return
	}

	// Maybe we've reached the slice capacity the tipping point?
	if uint16(len(t.objects)) >= t.maxObjects {
		t.subdivide(depth)
		t.insertIntoSubtree(key, obj, depth+1)
		return
	}

	// Trees which have been subdivided will have a non-nil subtrees at index 0
	if t.subTrees[0] != nil {
		t.insertIntoSubtree(key, obj, depth+1)
		return
	}

	// just store the point
	t.objects[key] = obj
	if t.insertCallback != nil {
		t.insertCallback(t, depth)
	}
}

// insertIntoSubtree determines which subtree to use, and calls Insert() on that subtree.
func (t *Tree) insertIntoSubtree(key uint64, obj Object, depth uint8) {
	for _, st := range t.subTrees {
		if overlap, fullyContained := obj.Overlaps(st.area); overlap {
			st.insert(key, obj, depth)
			if fullyContained {
				return
			}
		}
	}
}

func (t *Tree) overlaps(a Rectangle) bool {
	return t.area.Overlaps(a)
}

func (t *Tree) subdivide(depth uint8) {
	// create subtrees
	t.createSubtrees()

	// redistribute objects among new subtrees
	for k, v := range t.objects {
		t.insertIntoSubtree(k, v, depth+1)
	}

	// the objects slice is not going to be used again
	t.objects = nil
}

func NewTree(xMin, xMax, yMin, yMax int64, maxObjects uint16) *Tree {
	nt := newTree(newTreeCfg{
		xMin:       xMin,
		xMax:       xMax,
		yMin:       yMin,
		yMax:       yMax,
		maxObjects: maxObjects,
	})
	return nt
}

type newTreeCfg struct {
	xMin               int64
	xMax               int64
	yMin               int64
	yMax               int64
	insertCallbackFunc func(tree *Tree, depth uint8)
	maxObjects         uint16
}

func newTree(cfg newTreeCfg) *Tree {
	area := NewRectangle(NewLimits(cfg.xMin, cfg.xMax), NewLimits(cfg.yMin, cfg.yMax))

	return &Tree{
		area:            area,
		cannotSubdivide: area.cannotSubdivide(),
		insertCallback:  cfg.insertCallbackFunc,
		maxObjects:      cfg.maxObjects,
		objects:         make(map[uint64]Object),
	}
}
