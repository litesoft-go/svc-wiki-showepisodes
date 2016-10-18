package htmlplus

type RowShape struct {
	mShapes []*CellShape
}

func (this *RowShape) String() (rShapes string) {
	for _, zShape := range this.getShapes() {
		if rShapes != "" {
			rShapes += "+"
		}
		rShapes += zShape.String()
	}
	return
}

func (this *RowShape) getShapes() (rShapes []*CellShape) {
	if (this != nil) {
		rShapes = this.mShapes
	}
	return
}

func (this *RowShape) add(pShape *CellShape) *RowShape {
	zShapes := append(this.getShapes(), pShape)
	return &RowShape{mShapes:zShapes}
}

func (this *RowShape) length() int {
	return len(this.getShapes())
}

func (this *RowShape) shouldComeBefore(them *RowShape) bool {
	return this.length() > them.length() // More Cells (longer) should be checked first!
}
