package htmlplus

type RowsProcessor interface {
	GetExpectedShape() *RowsShape
	ProcessRowSet(pRowNumber int, pRows []*Row) error
}

type SimpleRowsProcessor struct {
	mPreRowFunc, mPostRowFunc func() error
	mRowsShape                *RowsShape
	mRowProcessors            []RowProcessor
}

func NewSimpleRowsProcessor(pPreRowFunc, pPostRowFunc func() error, pRowProcessors ...RowProcessor) RowsProcessor {
	zSRP := &SimpleRowsProcessor{mPreRowFunc:pPreRowFunc, mPostRowFunc:pPostRowFunc}
	for _, zRP := range pRowProcessors {
		zSRP.add(zRP)
	}
	return zSRP
}

func (this *SimpleRowsProcessor) add(pRowProcessor RowProcessor) {
	this.mRowsShape = this.mRowsShape.add(pRowProcessor.GetExpectedShape())
	this.mRowProcessors = append(this.mRowProcessors, pRowProcessor)
}

func (this *SimpleRowsProcessor) GetExpectedShape() *RowsShape {
	return this.mRowsShape
}

func (this *SimpleRowsProcessor) ProcessRowSet(pRowNumber int, pRows []*Row) (err error) {
	err = runFunc(err, this.mPreRowFunc)
	for i := range pRows {
		if err != nil {
			err = this.mRowProcessors[i].ProcessRow(pRowNumber + i, pRows[i])
		}
	}
	err = runFunc(err, this.mPostRowFunc)
	return
}

func runFunc(pError error, pFunc func() error) (err error) {
	err = pError
	if (err == nil) && pFunc != nil {
		err = pFunc()
	}
	return
}
