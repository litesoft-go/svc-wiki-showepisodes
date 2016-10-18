package htmlplus

type RowShape struct {
	mShapes []*CellShape
}

func (this *RowShape) String() (rShapes string) {
	for _, zShape := range this.getCells() {
		if rShapes != "" {
			rShapes += "+"
		}
		rShapes += zShape.String()
	}
	return
}

func (this *RowShape) getCells() (rShapes []*CellShape) {
	if (this != nil) {
		rShapes = this.mShapes
	}
	return
}

func (this *RowShape) add(pShape *CellShape) *RowShape {
	var zCellShapes []*CellShape
	if (this != nil) {
		zCellShapes = this.mShapes
	}
	zCellShapes := append(this.mShapes, pShape)
	return RowShape{mShapes:zCellShapes}
}

func (this *RowShape) length() (rLength int) {
	if (this != nil) {
		rLength = len(this.mShapes)
	}
	return
}

func (this *RowShape) shouldComeBefore(them *RowShape) bool {
	return this.length() > them.length() // More Cells (longer) should be checked first!
}
