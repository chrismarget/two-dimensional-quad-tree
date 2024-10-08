# two-dimensional-quad-tree

A quadtree for fast retrieval of two (or fewer) dimensional objects in 2D space.

Zero-dimension objects are straightforward: Because they're dimensionless points
(a coordinate pair) they exist neatly inside or outside of any defined region of
a coordinate plane.

Lines and polygons are more complicated: They might fall *partly* inside or
outside of a region of the coordinate plane.

Further complicating things, as insertions into the tree lead to splitting
regions of the tree, an object may find itself existing in multiple leaf
regions.

![image](diagram.svg) – 

This package solves the dilemma of insertion and retrieval of
larger-than-leaf-region objects by flipping the usual "Region, do you contain
point?" question around: It is now "Object, do you overlap region?"

This way, the burden of determining whether an arbitrarily complicated shape
"belongs to" any given region falls not on the region, but on the object.

## Objects

This package includes two simple sample implementations:
- `ColorPoint` an (x,y) coordinate pair and a color
- `ColorLine` two (x,y) coordinate pairs and a color

It is assumed that callers will provide their own implementations of the
`Object` interface suited to their needs.

Implementations of the `Object` interface require just two methods:
- `Hash() uint64` returns a value suitable for use as a map key. The hash value
 is primarily useful for eliminating duplicate results when tree searches find
 the same object in multiple regions. `Hash()` will be called multiple times as
 the tree regions are split, so it may be useful for implementations to return
 a stored value rather than recalculating it on the fly.
- `Overlaps(Rectangle) (bool, bool)` This method is used to determine whether an
 `Object` should be considered "part of" any rectangular area. It is used in
 insertions, reinsertions (during split operations), and when selecting
 `Objects` for retrieval with `Tree.Search()`.
