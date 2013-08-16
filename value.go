package fastheap

type Value interface {
	Index() int
	SetIndex(int)
}

type IntValue interface {
	Value
	Value() int64
}

type UintValue interface {
	Value
	Value() uint64
}

type FloatValue interface {
	Value
	Value() float64
}
