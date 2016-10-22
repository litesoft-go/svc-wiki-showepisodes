package htmlplus

import (
	"fmt"
	"lib-builtin/lib/ints"
)

type RowsProcessors struct {
	RowsShapeMapping
}

func NewRowsProcessors() *RowsProcessors {
	return &RowsProcessors{}
}

func (this *RowsProcessors) Add(pProcessor RowsProcessor) *RowsProcessors {
	this.add(pProcessor)
	return this
}

func (this *RowsProcessors) Process(pRows *RowStream) error {
	for zProxy := pRows.Next(); zProxy != nil; zProxy = pRows.Next() {
		zSet, err := createSet(pRows).populate(zProxy)
		if err != nil {
			return err
		}

		if zSet == nil {
			// zKey := asKey(zSet.mShapes)

			// TODO: XXX
		}
	}
	return nil
}

func createSet(pStream *RowStream) *rowSet {
	return &rowSet{mStream:pStream}
}

type rowSet struct {
	mStream       *RowStream
	m1stRowNumber int
	mRows         []*Row
	mShape        *RowsShape
}

func (this *rowSet) populate(pProxy *RowProxy) (rSet *rowSet, err error) {
	zRowNumber, zAdditionalRows, err := this.add(pProxy.GetRow())
	rSet, this.m1stRowNumber = this, zRowNumber
	if err == nil {
		zAdditionalRowsRemaining := zAdditionalRows
		for zAdditionalRowsRemaining > 0 {
			zAdditionalRowsRemaining--
			zProxy := this.mStream.Next()
			if zProxy == nil {
				err = fmt.Errorf("expected row %d, but wasn't one", zRowNumber + 1)
				return
			}
			zRowNumber, zAdditionalRows, err = this.add(pProxy.GetRow())
			zAdditionalRowsRemaining = ints.Max(zAdditionalRowsRemaining, zAdditionalRows)
		}
	}
	return
}

func (this *rowSet) add(pRowNumber int, pRow *Row) (rRowNumber, rAdditionalRows int, err error) {
	rRowNumber = pRowNumber
	zRowShape := pRow.GetRowShape()
	this.mShape = this.mShape.add(zRowShape)
	this.mRows = append(this.mRows, pRow)
	rAdditionalRows, err = zRowShape.calculateAdditionalRows(pRowNumber)
	return
}
