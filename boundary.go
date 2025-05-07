package ipdb

const (
	lb boundaryType = 0 // lower boundary
	ub boundaryType = 1 // upper boundary
	db boundaryType = 2 // double boundary - single ip range
)

type boundary struct {
	IP           IPv6
	BoundaryType boundaryType
	ValueRowID   int64
}

func (b boundary) IsUpperBoundary() bool {
	return b.BoundaryType == ub
}

func (b *boundary) SetUpperBoundary() {
	b.BoundaryType = ub
}

func (b boundary) IsLowerBoundary() bool {
	return b.BoundaryType == lb
}

func (b *boundary) SetLowerBoundary() {
	b.BoundaryType = lb
}

func (b boundary) IsDoubleBoundary() bool {
	return b.BoundaryType == db
}

func (b *boundary) SetDoubleBoundary() {
	b.BoundaryType = db
}
