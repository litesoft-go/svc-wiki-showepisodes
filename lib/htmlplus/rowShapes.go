package htmlplus

type RowShapes struct {
	mShapes []RowShape
}

func (this RowShapes) add(pShape RowShape) RowShapes {
	zRowShapes := append(this.mShapes, pShape)
	return RowShapes{mShapes:zRowShapes}
}

func (this RowShapes) String() (rShapes string) {
	for _, zShape := range this.mShapes {
		if rShapes != "" {
			rShapes += "r:"
		}
		rShapes += zShape.String()
	}
	return
}

func (this RowShapes) shouldComeBefore(them RowShapes) bool {
	for i := range this.mShapes {
		if !this.mShapes[i].shouldComeBefore(them.mShapes[i]) {
			return false
		}
	}
	return true
}
