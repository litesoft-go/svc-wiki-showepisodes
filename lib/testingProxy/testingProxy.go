package testingProxy

import (
	"testing"
	"fmt"
)

type TP struct {
	mT *testing.T
}

func NewTestingProxy(pT *testing.T) *TP {
	return &TP{mT:pT}
}

func (this *TP) EqualsBool(pExpected, pActual bool, pWhat string) {
	this.commonChk((pExpected == pActual), pWhat, FuncBool(pExpected), FuncBool(pActual))
}

func FuncBool(pValue bool) func() string {
	return func() string {
		return fmt.Sprintf("%v", pValue)
	}
}

func (this *TP) EqualsInt(pExpected, pActual int, pWhat string) {
	 this.commonChk((pExpected == pActual), pWhat, FuncInt(pExpected), FuncInt(pActual))
}

func FuncInt(pValue int) func() string {
	return func() string {
		return fmt.Sprintf("%d", pValue)
	}
}

func (this *TP) EqualsString(pExpected, pActual, pWhat string) {
	 this.commonChk((pExpected == pActual), pWhat, FuncString(pExpected), FuncString(pActual))
}

func FuncString(pValue string) func() string {
	return func() string {
		zBytes := []byte(pValue)
		return fmt.Sprintf("'%s' (%v)", pValue, zBytes)
	}
}

func (this *TP) commonChk(pPassed bool, pWhat string, pFuncExpected, pFuncActual func() string) {
	if !pPassed {
		this.mT.Errorf("%s = %s, but expected: %s", pWhat, pFuncActual(), pFuncExpected())
	}
}
