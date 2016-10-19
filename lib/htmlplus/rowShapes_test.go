package htmlplus

import (
	"testing"
	"svc-wiki-showepisodes/lib/testingProxy"
)

//noinspection GoUnusedFunction
func TestRowShapes(t *testing.T) {
	assert, errored := testingProxy.New()

	var zRow0 *RowShape // ZeroValue == nil
	var zShapes0 *RowShapes // ZeroValue == nil

	assert.EqualsString("", zShapes0.String(), "zShapes0.String()")
	assert.EqualsInt(0, len(zShapes0.getShapes()), "len(zShapes0.getShapes())")
	assert.EqualsBool(false, zShapes0.shouldComeBefore(nil), "zShapes0.shouldComeBefore(nil)")
	assert.EqualsBool(false, zShapes0.shouldComeBefore(zShapes0), "zShapes0.shouldComeBefore(zShapes0)")

	zRow1 := zRow0.add(&CellShape{mRowspan:1, mColspan:2})
	zShapes1 := zShapes0.add(zRow1)

	assert.EqualsString("[1x2]", zShapes1.String(), "zShapes1.String()")
	assert.EqualsInt(1, len(zShapes1.getShapes()), "len(zShapes1.getShapes())")
	assert.EqualsBool(true, zShapes1.shouldComeBefore(nil), "zShapes1.shouldComeBefore(nil)")
	assert.EqualsBool(true, zShapes1.shouldComeBefore(zShapes0), "zShapes1.shouldComeBefore(zShapes0)")
	assert.EqualsBool(false, zShapes1.shouldComeBefore(zShapes1), "zShapes1.shouldComeBefore(zShapes1)")
	assert.EqualsBool(false, zShapes0.shouldComeBefore(zShapes1), "zShapes0.shouldComeBefore(zShapes1)")

	zRow2 := zRow1.add(&CellShape{mRowspan:3, mColspan:4})
	zShapes2 := zShapes1.add(zRow2)

	assert.EqualsString("[1x2]r:[1x2]+[3x4]", zShapes2.String(), "zShapes2.String()")
	assert.EqualsInt(2, len(zShapes2.getShapes()), "len(zShapes2.getShapes())")
	assert.EqualsBool(true, zShapes2.shouldComeBefore(nil), "zShapes2.shouldComeBefore(nil)")
	assert.EqualsBool(true, zShapes2.shouldComeBefore(zShapes0), "zShapes2.shouldComeBefore(zShapes0)")
	assert.EqualsBool(true, zShapes2.shouldComeBefore(zShapes1), "zShapes2.shouldComeBefore(zShapes1)")
	assert.EqualsBool(false, zShapes2.shouldComeBefore(zShapes2), "zShapes2.shouldComeBefore(zShapes2)")
	assert.EqualsBool(false, zShapes1.shouldComeBefore(zShapes2), "zShapes1.shouldComeBefore(zShapes2)")
	assert.EqualsBool(false, zShapes0.shouldComeBefore(zShapes2), "zShapes0.shouldComeBefore(zShapes2)")

	zRow3 := zRow2.add(&CellShape{mRowspan:1, mColspan:1})
	zShapes3 := zShapes1.add(zRow3)

	assert.EqualsBool(false, zShapes2.shouldComeBefore(zShapes3), "zShapes2.shouldComeBefore(zShapes3)")
	assert.EqualsBool(true, zShapes3.shouldComeBefore(zShapes2), "zShapes3.shouldComeBefore(zShapes2)")

	if err := errored(); err != nil {
		t.Error(err)
	}
}
