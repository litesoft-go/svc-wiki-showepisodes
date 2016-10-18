package htmlplus

import (
	"testing"
	"fmt"
)

//noinspection GoUnusedFunction
func TestRowShape(t *testing.T) {
	var err error // nil

	var zShape0 *RowShape // ZeroValue == nil

	err = EqualsString(err, "", zShape0.String(), "zShape0.String()")
	err = EqualsInt(err, 0, zShape0.length(), "zShape0.length()")
	err = EqualsInt(err, 0, len(zShape0.getShapes()), "len(zShape0.getShapes())")
	err = EqualsBool(err, false, zShape0.shouldComeBefore(nil), "zShape0.shouldComeBefore(nil)")
	err = EqualsBool(err, false, zShape0.shouldComeBefore(zShape0), "zShape0.shouldComeBefore(zShape0)")

	zShape1 := zShape0.add(&CellShape{mRowspan:1,mColspan:2})

	err = EqualsString(err, "[1x2]", zShape1.String(), "zShape1.String()")
	err = EqualsInt(err, 1, zShape1.length(), "zShape1.length()")
	err = EqualsInt(err, 1, len(zShape1.getShapes()), "len(zShape1.getShapes())")
	err = EqualsBool(err, true, zShape1.shouldComeBefore(nil), "zShape1.shouldComeBefore(nil)")
	err = EqualsBool(err, true, zShape1.shouldComeBefore(zShape0), "zShape1.shouldComeBefore(zShape0)")
	err = EqualsBool(err, false, zShape1.shouldComeBefore(zShape1), "zShape1.shouldComeBefore(zShape1)")
	err = EqualsBool(err, false, zShape0.shouldComeBefore(zShape1), "zShape0.shouldComeBefore(zShape1)")

	zShape2 := zShape1.add(&CellShape{mRowspan:3,mColspan:4})

	err = EqualsString(err, "[1x2]+[3x4]", zShape2.String(), "zShape2.String()")
	err = EqualsInt(err, 2, zShape2.length(), "zShape2.length()")
	err = EqualsInt(err, 2, len(zShape2.getShapes()), "len(zShape2.getShapes())")
	err = EqualsBool(err, true, zShape2.shouldComeBefore(nil), "zShape2.shouldComeBefore(nil)")
	err = EqualsBool(err, true, zShape2.shouldComeBefore(zShape0), "zShape2.shouldComeBefore(zShape0)")
	err = EqualsBool(err, true, zShape2.shouldComeBefore(zShape1), "zShape2.shouldComeBefore(zShape1)")
	err = EqualsBool(err, false, zShape2.shouldComeBefore(zShape2), "zShape2.shouldComeBefore(zShape2)")
	err = EqualsBool(err, false, zShape1.shouldComeBefore(zShape2), "zShape1.shouldComeBefore(zShape2)")
	err = EqualsBool(err, false, zShape0.shouldComeBefore(zShape2), "zShape0.shouldComeBefore(zShape2)")

	if err != nil {
		t.Error(err)
	}
}

func EqualsBool(err error, pExpected, pActual bool, pWhat string) error {
	return commonChk(err, (pExpected == pActual), pWhat, FuncBool(pExpected), FuncBool(pActual))
}

func FuncBool(pValue bool) func() string {
	return func() string {
		return fmt.Sprintf("%v", pValue)
	}
}

func EqualsInt(err error, pExpected, pActual int, pWhat string) error {
	return commonChk(err, (pExpected == pActual), pWhat, FuncInt(pExpected), FuncInt(pActual))
}

func FuncInt(pValue int) func() string {
	return func() string {
		return fmt.Sprintf("%d", pValue)
	}
}

func EqualsString(err error, pExpected, pActual, pWhat string) error {
	return commonChk(err, (pExpected == pActual), pWhat, FuncString(pExpected), FuncString(pActual))
}

func FuncString(pValue string) func() string {
	return func() string {
		zBytes := []byte(pValue)
		return fmt.Sprintf("'%s' (%v)", pValue, zBytes)
	}
}

func commonChk(err error, pPassed bool, pWhat string, pFuncExpected, pFuncActual func() string) error {
	if (err == nil) && !pPassed {
			err = fmt.Errorf("%s = %s, but expected: %s", pWhat, pFuncActual(), pFuncExpected())
	}
	return err
}
