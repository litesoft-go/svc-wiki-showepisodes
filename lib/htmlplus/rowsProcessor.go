package htmlplus

type RowsProcessor interface {
	GetExpectedShape() *RowsShape
	ProcessRowSet(pRowNumber int, pRows []*Row) error
}

