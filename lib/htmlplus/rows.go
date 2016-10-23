package htmlplus

import (
	"fmt"
	"lib-builtin/lib/ints"
)

type TableMutator func(pTable *Table) error

type RowsProcessors struct {
	RowsShapeMapping
	mTableMutator TableMutator
}

func NewRowsProcessors() *RowsProcessors {
	return &RowsProcessors{}
}

func (this *RowsProcessors) With(pTableMutator TableMutator) *RowsProcessors {
	this.mTableMutator = pTableMutator
	return this
}

func (this *RowsProcessors) Add(pProcessor RowsProcessor) *RowsProcessors {
	this.add(pProcessor)
	return this
}

func (this *RowsProcessors) GetTableMutator() TableMutator {
	return this.mTableMutator
}

func (this *RowsProcessors) Process(pRows *RowStream) error {
	for zProxy := pRows.Next(); zProxy != nil; zProxy = pRows.Next() {
		zSet, err := createSet(pRows).populate(zProxy)
		if err != nil {
			return err
		}
		zRowsProcessor := this.getProcessorFor(zSet.mShape)
		if zRowsProcessor == nil {
			return this.noProcessorFoundFor(zSet.mShape)
		}
		err = zRowsProcessor.ProcessRowSet(zSet.m1stRowNumber, zSet.mRows)
		if err != nil {
			return err
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
	//fmt.Println("RowSet:")
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
			zRowNumber, zAdditionalRows, err = this.add(zProxy.GetRow())
			zAdditionalRowsRemaining = ints.Max(zAdditionalRowsRemaining, zAdditionalRows)
		}
	}
	return
}

func (this *rowSet) add(pRowNumber int, pRow *Row) (rRowNumber, rAdditionalRows int, err error) {
	//fmt.Printf("   row[%d]: %v\n", pRowNumber, pRow)
	rRowNumber = pRowNumber
	zRowShape := pRow.GetRowShape()
	this.mShape = this.mShape.add(zRowShape)
	this.mRows = append(this.mRows, pRow)
	rAdditionalRows, err = zRowShape.calculateAdditionalRows(pRowNumber)
	return
}
