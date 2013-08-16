package fastheap

// Value is interface for manipulating value indices
type Value interface {
	Index() int
	SetIndex(int)
}

// IntValue is a 'reference' to a integer item
type IntValue interface {
	Value
	Value() int64
}

// UintValue is a 'reference' to a unsigned integer item
type UintValue interface {
	Value
	Value() uint64
}

// FloatValue is a 'reference' to a float item
type FloatValue interface {
	Value
	Value() float64
}
