package htmlplus

import (
	"testing"
	"svc-wiki-showepisodes/lib/testingProxy"
)

//noinspection GoUnusedFunction
func TestRowShape(t *testing.T) {
	assert := testingProxy.New()

	var zShape0 *RowShape // ZeroValue == nil

	assert.EqualsString("", zShape0.String(), "zShape0.String()")
	assert.EqualsInt(0, zShape0.length(), "zShape0.length()")
	assert.EqualsInt(0, len(zShape0.getShapes()), "len(zShape0.getShapes())")
	assert.EqualsBool(false, zShape0.shouldComeBefore(nil), "zShape0.shouldComeBefore(nil)")
	assert.EqualsBool(false, zShape0.shouldComeBefore(zShape0), "zShape0.shouldComeBefore(zShape0)")

	zShape1 := zShape0.add(&CellShape{mRowspan:1,mColspan:2})

	assert.EqualsString("[1x2]", zShape1.String(), "zShape1.String()")
	assert.EqualsInt(1, zShape1.length(), "zShape1.length()")
	assert.EqualsInt(1, len(zShape1.getShapes()), "len(zShape1.getShapes())")
	assert.EqualsBool(true, zShape1.shouldComeBefore(nil), "zShape1.shouldComeBefore(nil)")
	assert.EqualsBool(true, zShape1.shouldComeBefore(zShape0), "zShape1.shouldComeBefore(zShape0)")
	assert.EqualsBool(false, zShape1.shouldComeBefore(zShape1), "zShape1.shouldComeBefore(zShape1)")
	assert.EqualsBool(false, zShape0.shouldComeBefore(zShape1), "zShape0.shouldComeBefore(zShape1)")

	zShape2 := zShape1.add(&CellShape{mRowspan:3,mColspan:4})

	assert.EqualsString("[1x2]+[3x4]", zShape2.String(), "zShape2.String()")
	assert.EqualsInt(2, zShape2.length(), "zShape2.length()")
	assert.EqualsInt(2, len(zShape2.getShapes()), "len(zShape2.getShapes())")
	assert.EqualsBool(true, zShape2.shouldComeBefore(nil), "zShape2.shouldComeBefore(nil)")
	assert.EqualsBool(true, zShape2.shouldComeBefore(zShape0), "zShape2.shouldComeBefore(zShape0)")
	assert.EqualsBool(true, zShape2.shouldComeBefore(zShape1), "zShape2.shouldComeBefore(zShape1)")
	assert.EqualsBool(false, zShape2.shouldComeBefore(zShape2), "zShape2.shouldComeBefore(zShape2)")
	assert.EqualsBool(false, zShape1.shouldComeBefore(zShape2), "zShape1.shouldComeBefore(zShape2)")
	assert.EqualsBool(false, zShape0.shouldComeBefore(zShape2), "zShape0.shouldComeBefore(zShape2)")

	if err := assert.Error(); err != nil {
		t.Error(err)
	}
}
