package htmlplus

import (
	"lib-builtin/lib/augmentor"
)

type RowProcessor interface {
	GetExpectedShape() *RowShape
	ProcessRow(pRowNumber int, pRow *Row) error
}

type SimpleRowProcessor struct {
	mShape          *RowShape
	mCellProcessors []CellProcessor
}

func NewSimpleRowProcessor(pCellProcessors ...CellProcessor) RowProcessor {
	zRP := &SimpleRowProcessor{}
	for _, zCP := range pCellProcessors {
		zRP.mShape = zRP.mShape.add(zCP.GetExpectedShape())
		zRP.mCellProcessors = append(zRP.mCellProcessors, zCP)
	}
	return zRP
}

func (this *SimpleRowProcessor) GetExpectedShape() *RowShape {
	return this.mShape
}

func (this *SimpleRowProcessor) ProcessRow(pRowNumber int, pRow *Row) error {
	//fmt.Printf("ProcessRow -> row[%d]: %v\n", pRowNumber, pRow)
	zCells := pRow.GetCells()
	for zCellNumber, zProcessor := range this.mCellProcessors {
		err := zProcessor.ProcessCell(zCells[zCellNumber])
		if err != nil {
			return augmentor.Err(err, "row[%d].cell[%d]", pRowNumber, zCellNumber)
		}
	}
	return nil
}
