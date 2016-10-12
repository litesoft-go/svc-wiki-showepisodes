package page

import (
	html "svc-wiki-showepisodes/lib/htmlplus"

	"github.com/pkg/errors"

	"lib-builtin/lib/augmentor"
	"lib-builtin/lib/fatal"
	"lib-builtin/lib/ints"

	"strconv"
	"fmt"
)

// SOT == SeriesOverviewTable

type SeriesOverviewTableProcessor func(pTable *html.Table) (rSeasons []*season, err error)

type SOT_set struct {
	mProcessor  SeriesOverviewTableProcessor
	mID         string
	mHeaderRows []html.HeaderRow
}

var sSOT_sets []SOT_set

func addSeriesOverview(pProcessor SeriesOverviewTableProcessor, pID string, pHeaderRows ...html.HeaderRow) {
	sSOT_sets = append(sSOT_sets, SOT_set{mProcessor:pProcessor, mID:pID, mHeaderRows:pHeaderRows})
}

func determineProcessorSOT(pTable *html.Table) (SeriesOverviewTableProcessor, error) {
	//fmt.Print(pTable.FormatHeader("Searching For:"))
	for _, zSet := range sSOT_sets {
		if pTable.HeaderMatches(zSet.mHeaderRows) {
			fmt.Println(zSet.mID)
			return zSet.mProcessor, nil
		}
	}
	return nil, pTable.ErrorHeaderNotMatched()
}

type SOTrowProcessor interface {
	GetExpectedCells() int
	ProcessRow(pNumber, pExpectedSeasonNumber int, pRow *html.Row, pStream *html.RowStream) (*season, error)
}

type SOTrowProcessors struct {
	mRowProcessorsByCellCount map[int]SOTrowProcessor
	mAcceptableLengths        string
}

func newSOTrowProcessors() *SOTrowProcessors {
	return &SOTrowProcessors{mRowProcessorsByCellCount:make(map[int]SOTrowProcessor)}
}

func (this *SOTrowProcessors) add(pProcessor SOTrowProcessor) *SOTrowProcessors {
	zCellCount := pProcessor.GetExpectedCells()
	_, zExisting := this.mRowProcessorsByCellCount[zCellCount]
	fatal.IfTrue(zExisting, "duplicate 'SOTrowProcessor' for cell count: %d", zCellCount)
	this.mRowProcessorsByCellCount[zCellCount] = pProcessor
	if this.mAcceptableLengths != "" {
		this.mAcceptableLengths = this.mAcceptableLengths + ", "
	}
	this.mAcceptableLengths = this.mAcceptableLengths + strconv.Itoa(zCellCount)
	return this
}

func (this *SOTrowProcessors) getProcessor(pRow *html.Row) (rProcessor SOTrowProcessor, err error) {
	zCellCount := len(pRow.GetCells())
	rProcessor, ok := this.mRowProcessorsByCellCount[zCellCount]
	if !ok {
		err = errors.Errorf("expected (%s) cells, but got: %d", this.mAcceptableLengths, zCellCount)
	}
	return
}

func (this *SOTrowProcessors) process(pProxy *html.RowProxy, pRows *html.RowStream) (rSeason *season, err error) {
	i, zRow := pProxy.GetRow()
	zProcessor, err := this.getProcessor(zRow)
	if err == nil {
		rSeason, err = zProcessor.ProcessRow(i, i + 1, zRow, pRows)
	}
	return
}

type cell4SOT struct {
	mColspan               int // Expected
	mSeasonAttributeSetter setSeasonAttributeFromText
}

func newSOTcell(pSeasonAttributeSetter setSeasonAttributeFromText) cell4SOT {
	return cell4SOT{mColspan:1, mSeasonAttributeSetter:pSeasonAttributeSetter}
}

func (this cell4SOT) colspan(pColspan int) cell4SOT {
	this.mColspan = pColspan // mutating copy
	return this
}

var sSOTcellIgnored = newSOTcell(nil)
var sSTOcellSeasonNumber = newSOTcell(setSeasonNumber)
var sSTOcellEpisodeCount = newSOTcell(setEpisodeCount)
var sSTOcellFirstAirDate = newSOTcell(setFirstAirDate)
var sSTOcellLastAirDate = newSOTcell(setLastAirDate)

func populateFromSOT(pTable *html.Table, pSOTrowProcessors *SOTrowProcessors) (rSeasons []*season, err error) {
	var zSeason *season
	zRows := pTable.GetBodyRowsAsStream()
	for zProxy := zRows.Next(); zProxy != nil; zProxy = zRows.Next() {
		zSeason, err = pSOTrowProcessors.process(zProxy, zRows)
		if err != nil {
			err = augmentor.Err(err, "body row %d", zProxy.GetNumber())
			return
		}
		if zSeason != nil {
			rSeasons = append(rSeasons, zSeason)
		}
	}
	if len(rSeasons) == 0 {
		err = errors.New("but no seasons")
	}
	return
}

func processCell(pCell *html.Cell, pSOTcell cell4SOT, pSeason *season) (err error) {
	err = pCell.AssertRowColSpan(1, pSOTcell.mColspan)
	if (err == nil) && (pSOTcell.mSeasonAttributeSetter != nil) {
		err = pSOTcell.mSeasonAttributeSetter(pSeason, pCell.GetText())
	}
	return
}

func populateFromSOTrow(pRowCells []*html.Cell, pExpectedSeasonNumber int, pSOTcells []cell4SOT) (*season, error) {
	if len(pSOTcells) != len(pRowCells) {
		return nil, errors.Errorf("expected %d cells, but got %d", len(pSOTcells), len(pRowCells))
	}
	zSeason := &season{}
	for i, zCell := range pRowCells {
		if err := processCell(zCell, pSOTcells[i], zSeason); err != nil {
			return nil, augmentor.Err(err, "col %d", i)
		}
	}
	return zSeason, ints.AssertEqual(pExpectedSeasonNumber, zSeason.mNumber, "expected season %d, but got %d")
}

func newSimpleSOTrowProcessor(pSOTcells ...cell4SOT) SOTrowProcessor {
	return &SimpleSOTrowProcessor{mSOTcells:pSOTcells}
}

type SimpleSOTrowProcessor struct {
	mSOTcells []cell4SOT
}

func (this *SimpleSOTrowProcessor) GetExpectedCells() int {
	return len(this.mSOTcells)
}

func (this *SimpleSOTrowProcessor) ProcessRow(pNumber, pExpectedSeasonNumber int, pRow *html.Row, pStream *html.RowStream) (*season, error) {
	return populateFromSOTrow(pRow.GetCells(), pExpectedSeasonNumber, this.mSOTcells)
}
