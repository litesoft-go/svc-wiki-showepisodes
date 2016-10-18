package htmlplus

type RowsProcessor interface {
	GetExpectedShape() []RowShape
	ProcessRowSet(pRowNumber int, pRows []*Row) error
}

