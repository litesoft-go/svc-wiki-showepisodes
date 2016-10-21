package htmlplus

import (
	"testing"
	"svc-wiki-showepisodes/lib/testingProxy"
)

var Rbuilder *RowShape // ZeroValue == nil

var ROWS_PROCESSOR_R_11 = newTestProcessor(Rbuilder.add(Cbuilder(1, 1)))
var ROWS_PROCESSOR_R_11_11 = newTestProcessor(Rbuilder.add(Cbuilder(1, 1)).add(Cbuilder(1, 1)))
var ROWS_PROCESSOR_R_12_R_11_11 = newTestProcessor(
	Rbuilder.add(Cbuilder(1, 2)),
	Rbuilder.add(Cbuilder(1, 1)).add(Cbuilder(1, 1)))
var ROWS_PROCESSOR_R_21_11_R_11 = newTestProcessor(
	Rbuilder.add(Cbuilder(2, 1)).add(Cbuilder(1, 1)),
	Rbuilder.add(Cbuilder(1, 1)))
var ROWS_PROCESSOR_R_11_21_R_11 = newTestProcessor(
	Rbuilder.add(Cbuilder(1, 1)).add(Cbuilder(2, 1)),
	Rbuilder.add(Cbuilder(1, 1)))
var ROWS_PROCESSOR_R_11_21_R_21_R_11 = newTestProcessor(
	Rbuilder.add(Cbuilder(1, 1)).add(Cbuilder(2, 1)),
	Rbuilder.add(Cbuilder(2, 1)),
	Rbuilder.add(Cbuilder(1, 1)))

var R_11_21_11_R_11_11 = RSbuilder(
	Rbuilder.add(Cbuilder(1, 1)).add(Cbuilder(2, 1)).add(Cbuilder(1, 1)),
	Rbuilder.add(Cbuilder(1, 1)).add(Cbuilder(1, 1)))

//noinspection GoUnusedFunction
func TestRowsShapeMapping(t *testing.T) {
	assert, errored := testingProxy.New()

	zMapping := &RowsShapeMapping{}
	zMapping.add(ROWS_PROCESSOR_R_11)
	zMapping.add(ROWS_PROCESSOR_R_11_11)
	zMapping.add(ROWS_PROCESSOR_R_12_R_11_11)
	zMapping.add(ROWS_PROCESSOR_R_21_11_R_11)
	zMapping.add(ROWS_PROCESSOR_R_11_21_R_11)
	zMapping.add(ROWS_PROCESSOR_R_11_21_R_21_R_11)

	//fmt.Print(zMapping)

	assert.SameInstance(ROWS_PROCESSOR_R_11_21_R_21_R_11, zMapping.getProcessorFor(ROWS_PROCESSOR_R_11_21_R_21_R_11.GetExpectedShape()), "ROWS_PROCESSOR_R_11_21_R_21_R_11")
	assert.SameInstance(ROWS_PROCESSOR_R_11_21_R_11, zMapping.getProcessorFor(ROWS_PROCESSOR_R_11_21_R_11.GetExpectedShape()), "ROWS_PROCESSOR_R_11_21_R_11")
	assert.SameInstance(ROWS_PROCESSOR_R_21_11_R_11, zMapping.getProcessorFor(ROWS_PROCESSOR_R_21_11_R_11.GetExpectedShape()), "ROWS_PROCESSOR_R_21_11_R_11")
	assert.SameInstance(ROWS_PROCESSOR_R_12_R_11_11, zMapping.getProcessorFor(ROWS_PROCESSOR_R_12_R_11_11.GetExpectedShape()), "ROWS_PROCESSOR_R_12_R_11_11")
	assert.SameInstance(ROWS_PROCESSOR_R_11_11, zMapping.getProcessorFor(ROWS_PROCESSOR_R_11_11.GetExpectedShape()), "ROWS_PROCESSOR_R_11_11")
	assert.SameInstance(ROWS_PROCESSOR_R_11, zMapping.getProcessorFor(ROWS_PROCESSOR_R_11.GetExpectedShape()), "ROWS_PROCESSOR_R_11")

	assert.SameInstance(ROWS_PROCESSOR_R_11_21_R_11, zMapping.getProcessorFor(R_11_21_11_R_11_11), "ROWS_PROCESSOR_R_11_21_R_11 from R_11_21_11_R_11_11")

	if err := errored(); err != nil {
		t.Error(err)
	}
}

func Cbuilder(pRowspan, pColspan int) *CellShape {
	return &CellShape{mRowspan:pRowspan, mColspan:pColspan}
}

func newTestProcessor(pRowShapes ...*RowShape) RowsProcessor {
	return &TestProcessor{mShapes:RSbuilder(pRowShapes...)}
}

func RSbuilder(pRowShapes ...*RowShape) (rShapes *RowsShape) {
	for _, zShape := range pRowShapes {
		rShapes = rShapes.add(zShape)
	}
	return
}

type TestProcessor struct {
	mShapes *RowsShape
}

func (this *TestProcessor) String() string {
	return "TestProcessor: " + this.mShapes.String()
}

func (this *TestProcessor) GetExpectedShape() *RowsShape {
	return this.mShapes
}

func (this *TestProcessor) ProcessRowSet(pRowNumber int, pRows []*Row) error {
	panic("NIY")
}