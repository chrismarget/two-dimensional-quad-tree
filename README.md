# two-dimensional-quad-tree

A quadtree for keeping track of zero, one and two dimensional objects for fast 
retrieval from regions of a coordinate plane.

Zero-dimension objects are straightforward: They're just a point (coordinate
pair) and therefore clearly fall inside or outside of any defined region of a
coordinate plane.

One dimensional and two dimensional objects (lines and polygons) are more
complicated: They might fall partly inside or outside of any region of the
coordinate plane.

![image](diagram.svg) – 

This package solves the dilemma of insertion and retrieval of larger-than-point
objects by flipping the usual "Region, do you contain point?" question around:
It is now "Object, do you overlap region?"

This way, the burden of determining whether an arbitrarily complicated shape
"belongs to" any given region falls not on the region, but on the object.