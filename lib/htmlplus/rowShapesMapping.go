package htmlplus

type RowsShapeProcessorPair struct {
	mShape     *RowsShape
	mProcessor RowsProcessor
}

func (this *RowsShapeProcessorPair) shouldComeBefore(them *RowsShapeProcessorPair) bool {
	return this.mShape.shouldComeBefore(them.mShape);
}

type RowsShapeProcessors struct {
	mPairs []*RowsShapeProcessorPair
}

func (this *RowsShapeProcessors) add(pShape *RowsShape, pProcessor RowsProcessor) {
	zNew := &RowsShapeProcessorPair{mShape:pShape, mProcessor:pProcessor}
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

func (this *RowsShapeMapping) getRowsShapeProcessors(pShape *RowsShape) (rRowsShapeProcessors *RowsShapeProcessors) {
	zKey := pShape.length()
	rRowsShapeProcessors, ok := this.mRowProcessorsByCellCount[zKey]
	if !ok {
		rRowsShapeProcessors = this.addEmpty(zKey)
	}
	return
}

func (this *RowsShapeMapping) add(pProcessor RowsProcessor) {
	zShape := pProcessor.GetExpectedShape()
	this.getRowsShapeProcessors(zShape).add(zShape, pProcessor)
}
