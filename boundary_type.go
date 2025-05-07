package ipdb

type boundaryType uint8

func (b boundaryType) String() string {
	switch b {
	case lb:
		return "lower boundary"
	case ub:
		return "upper boundary"
	case db:
		return "double boundary"
	default:
		return "unknown boundary"
	}
}
