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