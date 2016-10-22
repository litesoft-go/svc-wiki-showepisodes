package sot

import (
	html "svc-wiki-showepisodes/lib/htmlplus"

	"svc-wiki-showepisodes/lib/utils"

	//"github.com/pkg/errors"

	//"lib-builtin/lib/augmentor"
	//"lib-builtin/lib/bools"

	"fmt"
	"lib-builtin/lib/slices"
)

func ctp(pTestProcessingFunc html.CellTextProcessingFunc) *html.ProxyCellProcessor {
	return html.NewProxyCellTextProcessor(pTestProcessingFunc)
}

func cp(pProcessingFunc html.CellProcessingFunc) *html.ProxyCellProcessor {
	return html.NewProxyCellProcessor(pProcessingFunc)
}

// SOT == SeriesOverviewTable

func Process(pTable *html.Table) (rSeasons []*utils.Season, err error) {
	zFactory, err := determineFactory(pTable)
	if err == nil {
		zCollector := &SeasonCollector{}
		zRowsProcessors := zFactory(zCollector)
		err = zRowsProcessors.Process(pTable.GetBodyRowsAsStream())
		if err == nil {
			rSeasons = zCollector.mSeasons
		}
	}
	return
}

type Factory func(pCollector *SeasonCollector) *html.RowsProcessors

type FactorySet struct {
	mFactory    Factory
	mID         string
	mHeaderRows []html.HeaderRow
}

var sFactorySets []FactorySet

func addFactoryMapping(pFactory Factory, pID string, pHeaderRows ...html.HeaderRow) {
	sFactorySets = append(sFactorySets, FactorySet{mFactory:pFactory, mID:pID, mHeaderRows:pHeaderRows})
}

func determineFactory(pTable *html.Table) (Factory, error) {
	//fmt.Print(pTable.FormatHeader("Searching For:"))
	for _, zSet := range sFactorySets {
		if pTable.HeaderStartsWith(zSet.mHeaderRows) {
			fmt.Println(zSet.mID)
			return zSet.mFactory, nil
		}
	}
	return nil, pTable.ErrorHeaderNotMatched()
}

type SeasonCollector struct {
	mSeasons   []*utils.Season
	mCurSeason *utils.Season
}

func (this *SeasonCollector) initSeason() error {
	this.mCurSeason = &utils.Season{}
	return nil
}

func (this *SeasonCollector) finiSeason() error {
	this.mSeasons = append(this.mSeasons, this.mCurSeason)
	return nil
}

func (this *SeasonCollector) SeasonID(pCell *html.Cell) error {
	return this.mCurSeason.SetFragmentRefAndID(slices.FirstOrEmpty(pCell.GetHrefs()), html.FirstTextOnly(pCell.GetText()))
}

func (this *SeasonCollector) EpisodeCount(pCellText string) error {
	return this.mCurSeason.SetEpisodeCountFromText(html.FirstTextOnly(pCellText))
}

func (this *SeasonCollector) FirstAirDate(pCellText string) (err error) {
	zDate, err := utils.ExtractAirDate("FirstAirDate", pCellText)
	if err == nil {
		err = this.mCurSeason.SetFirstAirDate(zDate)
	}
	return
}

func (this *SeasonCollector) LastAirDate(pCellText string) (err error) {
	zDate, err := utils.ExtractAirDate("LastAirDate", pCellText)
	if err == nil {
		err = this.mCurSeason.SetLastAirDate(zDate)
	}
	return
}

func (this *SeasonCollector) newSingleRowProcessor(pCellProcessors ...html.CellProcessor) html.RowsProcessor {
	return html.NewSimpleRowsProcessor(this.initSeason, this.finiSeason, html.NewSimpleRowProcessor(pCellProcessors...))
}

//type RowProcessors struct {
//	mRowProcessorsByCellCount map[int]RowProcessor // TODO: Switch to Row "shaped" map!
//	mAcceptableLengths        string
//	mDefaultMinCellCount      int
//	mDefaultRowProcessor      RowProcessor
//}
//
//func newRowProcessors() *RowProcessors {
//	return &RowProcessors{mRowProcessorsByCellCount:make(map[int]RowProcessor), mDefaultMinCellCount:-1}
//}
//
//func (this *RowProcessors) add(pProcessor RowProcessor) *RowProcessors {
//	this.mDefaultRowProcessor = pProcessor
//	this.mDefaultMinCellCount = pProcessor.GetExpectedCells()
//	return this
//	//zCellCount := pProcessor.GetExpectedCells()
//	//_, zExisting := this.mRowProcessorsByCellCount[zCellCount]
//	//fatal.IfTrue(zExisting, "duplicate 'SOTrowProcessor' for cell count: %d", zCellCount)
//	//this.mRowProcessorsByCellCount[zCellCount] = pProcessor
//	//if this.mAcceptableLengths != "" {
//	//	this.mAcceptableLengths = this.mAcceptableLengths + ", "
//	//}
//	//this.mAcceptableLengths = this.mAcceptableLengths + strconv.Itoa(zCellCount)
//	//return this
//}
//
//func (this *RowProcessors) getProcessor(pRow *html.Row) (rProcessor RowProcessor, err error) {
//	zCellCount := len(pRow.GetCells())
//	rProcessor, ok := this.mRowProcessorsByCellCount[zCellCount]
//	if ok {
//		return
//	}
//	if (this.mDefaultRowProcessor != nil) && (this.mDefaultMinCellCount <= zCellCount) {
//		rProcessor = this.mDefaultRowProcessor
//		return
//	}
//	err = this.generateError(zCellCount)
//	return
//}
//
//func (this *RowProcessors) hasCellCountProcessors() bool {
//	return 0 < len(this.mRowProcessorsByCellCount)
//}
//
//func (this *RowProcessors) hasDefaultRowProcessor() bool {
//	return 0 <= this.mDefaultMinCellCount
//}
//
//func (this *RowProcessors) generateError(pCellCount int) error {
//	switch bools.TwoBools(this.hasCellCountProcessors(), this.hasDefaultRowProcessor()) {
//	case bools.LEFT_ONLY:
//		return errors.Errorf("expected (%s) cells, but got: %d", this.mAcceptableLengths, pCellCount)
//	case bools.RIGHT_ONLY:
//		return errors.Errorf("got row with %d cells, which is less than the default minimum of %d",
//			pCellCount, this.mDefaultMinCellCount)
//	case bools.BOTH:
//		return errors.Errorf("expected (%s) cells, but got: %d, which is less than the default minimum of %d",
//			this.mAcceptableLengths, pCellCount, this.mDefaultMinCellCount)
//	default: // Neither
//		return errors.Errorf("no processors, but got row with %d cells", pCellCount)
//	}
//}
//
//func (this *RowProcessors) process(pProxy *html.RowProxy, pRows *html.RowStream) (rSeason *utils.Season, err error) {
//	zRowNumber, zRow := pProxy.GetRow()
//	zProcessor, err := this.getProcessor(zRow)
//	if err == nil {
//		rSeason, err = zProcessor.ProcessRow(zRowNumber, zRow, pRows)
//	}
//	return
//}
//
//type cell4SOT struct {
//	mColspan               int // Expected
//	mSeasonAttributeSetter utils.SetSeasonAttributeFromText
//}
//
//func newSOTcell(pSeasonAttributeSetter utils.SetSeasonAttributeFromText) cell4SOT {
//	return cell4SOT{mColspan:1, mSeasonAttributeSetter:pSeasonAttributeSetter}
//}
//
//func (this cell4SOT) colspan(pColspan int) cell4SOT {
//	this.mColspan = pColspan // mutating copy
//	return this
//}
//
//var sSOTcellIgnored = newSOTcell(nil)
//var sSOTcellSeasonID = newSOTcell(utils.SetSeasonID)
//var sSOTcellEpisodeCount = newSOTcell(utils.SetEpisodeCount)
//var sSOTcellFirstAirDate = newSOTcell(utils.SetFirstAirDate)
//var sSOTcellLastAirDate = newSOTcell(utils.SetLastAirDate)
//
//func populateFromSOT(pTable *html.Table, pSOTrowProcessors *RowProcessors) (rSeasons []*utils.Season, err error) {
//	var zSeason *utils.Season
//	zRows := pTable.GetBodyRowsAsStream()
//	for zProxy := zRows.Next(); zProxy != nil; zProxy = zRows.Next() {
//		zSeason, err = pSOTrowProcessors.process(zProxy, zRows)
//		if err != nil {
//			err = augmentor.Err(err, "body row %d", zProxy.GetNumber())
//			return
//		}
//		if zSeason != nil {
//			rSeasons = append(rSeasons, zSeason)
//		}
//	}
//	if len(rSeasons) == 0 {
//		err = errors.New("but no seasons")
//	}
//	return
//}
//
//func processCell(pCell *html.Cell, pSOTcell cell4SOT, pSeason *utils.Season) (err error) {
//	err = pCell.AssertRowColSpan(1, pSOTcell.mColspan)
//	if (err == nil) && (pSOTcell.mSeasonAttributeSetter != nil) {
//		err = pSOTcell.mSeasonAttributeSetter(pSeason, pCell.GetText())
//	}
//	return
//}
//
//func populateFromSOTrow(pRowCells []*html.Cell, pSOTcells []cell4SOT) (rSeason *utils.Season, err error) {
//	if len(pRowCells) < len(pSOTcells) {
//		err = errors.Errorf("expected at least %d cells, but got %d", len(pSOTcells), len(pRowCells))
//		return
//	}
//	rSeason = &utils.Season{}
//	for i, zSOTcell := range pSOTcells {
//		if err = processCell(pRowCells[i], zSOTcell, rSeason); err != nil {
//			err = augmentor.Err(err, "col %d", i)
//			return
//		}
//	}
//	return
//}
//
//func newSimpleSOTrowProcessor(pSOTcells ...cell4SOT) RowProcessor {
//	return &SimpleSOTrowProcessor{mSOTcells:pSOTcells}
//}
//
//type SimpleSOTrowProcessor struct {
//	mSOTcells []cell4SOT
//}
//
//func (this *SimpleSOTrowProcessor) GetExpectedCells() int {
//	return len(this.mSOTcells)
//}
//
//func (this *SimpleSOTrowProcessor) ProcessRow(pRowNumber int, pRow *html.Row, pStream *html.RowStream) (*utils.Season, error) {
//	return populateFromSOTrow(pRow.GetCells(), this.mSOTcells)
//}
