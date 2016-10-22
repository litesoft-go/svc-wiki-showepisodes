package htmlplus

import "lib-builtin/lib/augmentor"

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
	for zCellNunmber, zCell := range pRow.GetCells() {
		err := this.mCellProcessors[zCellNunmber].ProcessCell(zCell)
		if err != nil {
			return augmentor.Err(err, "row[%d].cell[%d]", pRowNumber, zCellNunmber)
		}
	}
	return nil
}
