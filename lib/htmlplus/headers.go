package htmlplus

import (
	"lib-builtin/lib/slices"
	"lib-builtin/lib/lines"
)

type HeaderRow []string

func (this HeaderRow) String() string {
	return slices.AsOptions([]string(this))
}

func (this HeaderRow) Equals(them HeaderRow) bool {
	return slices.Equals([]string(this), []string(them)...)
}

func (this HeaderRow) StartsWith(them HeaderRow) bool {
	return slices.StartsWith([]string(this), []string(them)...)
}

func (this HeaderRow) shouldComeBefore(them HeaderRow) bool {
	return len([]string(this)) > len([]string(them)) // More Cells (longer) should be checked first!
}

type HeaderRows struct {
	mRows []HeaderRow
}

func NewHeaderRows(pHeaderRows ...HeaderRow) *HeaderRows {
	return &HeaderRows{mRows:pHeaderRows}
}

func (this *HeaderRows) getRows() (rRows []HeaderRow) {
	if (this != nil) {
		rRows = this.mRows
	}
	return
}

func (this *HeaderRows) Equals(them *HeaderRows) bool {
	zThisRows, zThemRows, zThisLength, zThemLength := this.getPairData(them)
	if zThisLength != zThemLength {
		return false
	}
	for i := range zThisRows {
		if !zThisRows[i].Equals(zThemRows[i]) {
			return false
		}
	}
	return true
}

func (this *HeaderRows) StartsWith(them *HeaderRows) bool {
	zThisRows, zThemRows, zThisLength, zThemLength := this.getPairData(them)
	if zThisLength != zThemLength {
		return false
	}
	for i := range zThisRows {
		if !zThisRows[i].StartsWith(zThemRows[i]) {
			return false
		}
	}
	return true
}

func (this *HeaderRows) ShouldComeBefore(them *HeaderRows) bool {
	zThisRows, zThemRows, zThisLength, zThemLength := this.getPairData(them)
	if zThisLength > zThemLength {
		return true // More Rows (longer) should be checked first!
	}
	if zThisLength < zThemLength {
		return false
	}
	for i := range this.mRows {
		if zThemRows[i].shouldComeBefore(zThisRows[i]) {
			return false
		}
	}
	return true // same # or more Cells should be checked first!
}

func (this *HeaderRows) getPairData(them *HeaderRows) (rThisRows, rThemRows []HeaderRow, rThisLength, rThemLength int) {
	rThisRows = this.getRows()
	rThemRows = them.getRows()
	rThisLength = len(rThisRows)
	rThemLength = len(rThemRows)
	return
}

func (this *HeaderRows) String() string {
	zCollector := lines.NewCollector()
	this.addHeaders(zCollector, "")
	return zCollector.String()
}

func (this *HeaderRows) addHeaders(pCollector *lines.Collector, pWhat string) {
	zIndent := pWhat != ""
	if zIndent {
		pCollector.Line(pWhat)
	}
	for _, zRow := range this.getRows() {
		if zIndent {
			pCollector.Indent()
		}
		pCollector.Line(zRow.String())
		if zIndent {
			pCollector.Outdent()
		}
	}
}


