package htmlplus

type RowShapes struct {
	mShapes []*RowShape
}

func (this *RowShapes) getShapes() (rShapes []*RowShape) {
	if (this != nil) {
		rShapes = this.mShapes
	}
	return
}

func (this *RowShapes) add(pShape *RowShape) *RowShapes {
	zShapes := append(this.getShapes(), pShape)
	return &RowShapes{mShapes:zShapes}
}

func (this *RowShapes) String() (rShapes string) {
	for _, zShape := range this.getShapes() {
		if rShapes != "" {
			rShapes += "r:"
		}
		rShapes += zShape.String()
	}
	return
}

func (this *RowShapes) shouldComeBefore(them *RowShapes) bool {
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
