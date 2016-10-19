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
	mShapes       []*RowShape
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
	rRowNumber, rAdditionalRows, zRowShape, err := calcShapeAndAdditionalRows(pRowNumber, pRow.GetCellShapes())
	if err == nil {
		this.mShapes = append(this.mShapes, zRowShape)
		this.mRows = append(this.mRows, pRow)
	}
	return
}

func calcShapeAndAdditionalRows(pRowNumber int, pCellShapes []*CellShape) (rRowNumber, rAdditionalRows int, rRowShape *RowShape, err error) {
	if len(pCellShapes) == 0 {
		err = fmt.Errorf("row %d : no cells", pRowNumber)
		return
	}
	for _, zCellShape := range pCellShapes {
		rRowShape = rRowShape.add(zCellShape)
		rAdditionalRows = ints.Max(rAdditionalRows, zCellShape.mRowspan - 1)
	}
	rRowNumber = pRowNumber
	return
}