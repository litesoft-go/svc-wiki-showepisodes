package htmlplus

import (
	"fmt"
	"lib-builtin/lib/ints"
)

type RowShape struct {
	mAcceptExactLengthOnly bool
	mShapes                []*CellShape
}

func (this *RowShape) String() (rShapes string) {
	for _, zShape := range this.getShapes() {
		if rShapes != "" {
			rShapes += "+"
		}
		rShapes += zShape.String()
	}
	if this.getAcceptExactLengthOnly() {
		rShapes += "!"
	}
	return
}

func (this *RowShape) getAcceptExactLengthOnly() bool {
	return (this != nil) && this.mAcceptExactLengthOnly // left to right!
}

func (this *RowShape) getShapes() (rShapes []*CellShape) {
	if (this != nil) {
		rShapes = this.mShapes
	}
	return
}

func (this *RowShape) calculateAdditionalRows(pRowNumber int) (rAdditionalRows int, err error) {
	zCellShapes := this.getShapes()
	if len(zCellShapes) == 0 {
		err = fmt.Errorf("row %d : no cells", pRowNumber)
		return
	}
	for _, zCellShape := range zCellShapes {
		rAdditionalRows = ints.Max(rAdditionalRows, zCellShape.mRowspan - 1)
	}
	return
}

func (this *RowShape) acceptExactLengthOnly() *RowShape {
	return &RowShape{mAcceptExactLengthOnly:true,mShapes:this.getShapes()}
}

func (this *RowShape) add(pShape *CellShape) *RowShape {
	zShapes := append(this.getShapes(), pShape)
	return &RowShape{mAcceptExactLengthOnly:this.getAcceptExactLengthOnly(),mShapes:zShapes}
}

func (this *RowShape) length() int {
	return len(this.getShapes())
}

// Accepts is a "first success wins" algorithm.
// What constitutes success is that:
// 1) exact matches are always true
// 2) AcceptExactLengthOnly is false and all the "ProcessorCellShapes" == the start of the "ActualRowCellShapes" then that is true!
func (this *RowShape) accepts(pActualShape *RowShape) bool {
	zProcessorCellShapes := this.getShapes()
	zActualRowCellShapes := pActualShape.getShapes()
	if len(zActualRowCellShapes) < len(zProcessorCellShapes) {
		return false
	}
	if this.getAcceptExactLengthOnly() && (len(zActualRowCellShapes) != len(zProcessorCellShapes)) {
		return false
	}
	for i := range zProcessorCellShapes {
		if !zProcessorCellShapes[i].Equals(zActualRowCellShapes[i]) {
			return false
		}
	}
	return true
}

// determine "order" that "accepts" will check.  Accepts is a "first success wins" algorithm.  As such "normally" longer should be checked first
func (this *RowShape) shouldComeBefore(them *RowShape) bool {
	return them.getAcceptExactLengthOnly() || (this.length() > them.length()) // More Cells (longer) should be checked first!
}
