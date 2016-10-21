package htmlplus

import (
	"lib-builtin/lib/ints"
	"fmt"
)

type RowsShapeProcessorPair struct {
	mShape     *RowsShape
	mProcessor RowsProcessor
}

func (this *RowsShapeProcessorPair) shouldComeBefore(them *RowsShapeProcessorPair) bool {
	return this.mShape.shouldComeBefore(them.mShape);
}

func (this *RowsShapeProcessorPair) canProcess(pShape *RowsShape) bool {
	return this.mShape.accepts(pShape);
}

type RowsShapeProcessors struct {
	mPairs []*RowsShapeProcessorPair
}

func (this *RowsShapeProcessors) String() (rv string) {
	for _, zPair := range this.mPairs {
		rv += "\n      " + zPair.mShape.String()
	}
	return
}

func (this *RowsShapeProcessors) findProcessor(pShape *RowsShape) (rProcessor RowsProcessor) {
	//fmt.Println("findProcessor(", pShape, ")")
	for _, zPair := range this.mPairs {
		//fmt.Print("     checking: ", zPair.mShape)
		if zPair.canProcess(pShape) {
			rProcessor = zPair.mProcessor
			//fmt.Println("  --  FOUND")
			return
		}
		//fmt.Println(" n/a")
	}
	return
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
	mKeys                     []int
	mRowProcessorsByCellCount map[int]*RowsShapeProcessors // Row Shape Processors by Row Count
}

func (this *RowsShapeMapping) String() (rv string) {
	rv = fmt.Sprintf("RowsShapeMapping (%d):\n", len(this.mKeys))
	for _,zKey := range this.mKeys {
		rv += fmt.Sprintf("   %d:%v\n", zKey, this.mRowProcessorsByCellCount[zKey])
	}
	return
}

func (this *RowsShapeMapping) addEmpty(pKey int) (rRowsShapeProcessors *RowsShapeProcessors) {
	if len(this.mRowProcessorsByCellCount) == 0 {
		this.mRowProcessorsByCellCount = make(map[int]*RowsShapeProcessors)
	}
	this.mKeys = ints.AddOrdered(this.mKeys, pKey)
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

func (this *RowsShapeMapping) getProcessorFor(pShape *RowsShape) (rProcessor RowsProcessor) {
	zKey := pShape.length()
	zRowsShapeProcessors, ok := this.mRowProcessorsByCellCount[zKey]
	if ok {
		rProcessor = zRowsShapeProcessors.findProcessor(pShape)
	}
	return
}
