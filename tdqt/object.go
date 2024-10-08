package tdqt

type Object interface {
	// Hash returns an ID suitable for use as a map key
	Hash() uint64

	// Overlaps indicate whether the object has *any* overlap with the specified
	// Rectangle (the first returned boolean), and whether the object is
	// *entirely contained within* the specified Rectangle (the second boolean)
	Overlaps(Rectangle) (bool, bool)
}
