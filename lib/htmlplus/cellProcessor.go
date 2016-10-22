package htmlplus

type CellProcessingFunc func(pCell *Cell) error

type CellTextProcessingFunc func(pCell string) error

type CellProcessor interface {
	GetExpectedShape() *CellShape
	ProcessCell(pCell *Cell) error
}

type ProxyCellProcessor struct {
	mProcessingFunc CellProcessingFunc
	CellShape
}

func NewProxyCellTextProcessor(pTextProcessingFunc CellTextProcessingFunc) *ProxyCellProcessor {
	return NewProxyCellProcessor(func(pCell *Cell) error {
		return pTextProcessingFunc(pCell.GetText())
	})
}

func NewProxyCellProcessor(pProcessingFunc CellProcessingFunc) *ProxyCellProcessor {
	return &ProxyCellProcessor{mProcessingFunc:pProcessingFunc, CellShape:CellShape{mRowspan:1, mColspan:1}}
}

func (this *ProxyCellProcessor) Rowspan(pRowspan int) *ProxyCellProcessor {
	return &ProxyCellProcessor{mProcessingFunc:this.mProcessingFunc, CellShape:CellShape{mRowspan:pRowspan, mColspan:this.mColspan}}
}

func (this *ProxyCellProcessor) Colspan(pColspan int) *ProxyCellProcessor {
	return &ProxyCellProcessor{mProcessingFunc:this.mProcessingFunc, CellShape:CellShape{mRowspan:this.mRowspan, mColspan:pColspan}}
}

func (this *ProxyCellProcessor) GetExpectedShape() *CellShape {
	return &this.CellShape
}

func (this *ProxyCellProcessor) ProcessCell(pCell *Cell) (err error) {
	if this.mProcessingFunc != nil {
		err = this.mProcessingFunc(pCell)
	}
	return
}

var S_IGNORED_CELL_PROCESSOR = NewProxyCellProcessor(nil)