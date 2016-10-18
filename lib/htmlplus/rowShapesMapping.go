package htmlplus

type RowsShapeProcessorPair struct {
	mShapes    []RowShape
	mProcessor RowsProcessor
}

func (this *RowsShapeProcessorPair) shouldComeBefore(them RowsShapeProcessorPair) bool {
	for i := range this.mShapes {
		if !this.mShapes[i].shouldComeBefore(them.mShapes[i]) {
			return false
		}
	}
	return true
}

type RowsShapeProcessors struct {
	mPairs []*RowsShapeProcessorPair
}

func (this *RowsShapeProcessors) add(pShapes []RowShape, pProcessor RowsProcessor) {
	zNew := &RowsShapeProcessorPair{mShapes:pShapes, mProcessor:pProcessor}
	if len(this.mPairs) == 0 {
		this.mPairs = append(this.mPairs, zNew)
		return
	}
	var zLeft, zRight []*RowsShapeProcessorPair
	for _, zCurrent := range this.mPairs {
		if zCurrent.shouldComeBefore(zNew) {
			zLeft = append(zLeft, zCurrent)
		} else {
			zRight = append(zRight, zCurrent)
		}
	}
	this.mPairs = append(append(zLeft, zNew), zRight...)
}

type RowsShapeMapping struct {
	mRowProcessorsByCellCount map[int]*RowsShapeProcessors // Row Shape Processors by Row Count
}

func (this *RowsShapeMapping) addEmpty(pKey int) (rRowsShapeProcessors *RowsShapeProcessors) {
	if len(this.mRowProcessorsByCellCount) == 0 {
		this.mRowProcessorsByCellCount = make(map[int]*RowsShapeProcessors)
	}
	rRowsShapeProcessors = &RowsShapeProcessors{}
	this.mRowProcessorsByCellCount[pKey] = rRowsShapeProcessors
	return
}

func (this *RowsShapeMapping) getRowsShapeProcessors(pShapes []RowShape) (rRowsShapeProcessors *RowsShapeProcessors) {
	zKey := len(pShapes)
	rRowsShapeProcessors, ok := this.mRowProcessorsByCellCount[zKey]
	if !ok {
		rRowsShapeProcessors = this.addEmpty(zKey)
	}
	return
}

func (this *RowsShapeMapping) add(pProcessor RowsProcessor) {
	zShapes := pProcessor.GetExpectedShape()
	this.getRowsShapeProcessors(zShapes).add(zShapes, pProcessor)
}
