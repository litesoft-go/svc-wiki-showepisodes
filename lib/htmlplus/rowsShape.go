package htmlplus

type RowsShape struct {
	mShapes []*RowShape
}

func (this *RowsShape) getShapes() (rShapes []*RowShape) {
	if (this != nil) {
		rShapes = this.mShapes
	}
	return
}

func (this *RowsShape) add(pShape *RowShape) *RowsShape {
	zShapes := append(this.getShapes(), pShape)
	return &RowsShape{mShapes:zShapes}
}

func (this *RowsShape) String() (rShapes string) {
	for _, zShape := range this.getShapes() {
		if rShapes != "" {
			rShapes += "r:"
		}
		rShapes += zShape.String()
	}
	return
}

func (this *RowsShape) length() int {
	return len(this.getShapes())
}

func (this *RowsShape) shouldComeBefore(them *RowsShape) bool {
	thisShapes := this.getShapes()
	thisLength := len(thisShapes)
	themShapes := them.getShapes()
	themLength := len(themShapes)
	// for robustness we will treat non-equal lengths as longer should come Before shorter
	if thisLength > themLength {
		return true
	}
	if thisLength == themLength {
		for i := range thisShapes {
			if thisShapes[i].shouldComeBefore(themShapes[i]) {
				return true
			}
		}
	}
	return false
}
