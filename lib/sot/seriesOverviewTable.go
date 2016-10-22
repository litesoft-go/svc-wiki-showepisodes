package sot

import (
	html "svc-wiki-showepisodes/lib/htmlplus"

	"svc-wiki-showepisodes/lib/utils"

	//"github.com/pkg/errors"

	//"lib-builtin/lib/augmentor"
	//"lib-builtin/lib/bools"

	"fmt"
	"lib-builtin/lib/slices"
	"lib-builtin/lib/augmentor"
)

var skipCell = html.S_IGNORED_CELL_PROCESSOR

func ctp(pTestProcessingFunc html.CellTextProcessingFunc) *html.ProxyCellProcessor {
	return html.NewProxyCellTextProcessor(pTestProcessingFunc)
}

func cp(pProcessingFunc html.CellProcessingFunc) *html.ProxyCellProcessor {
	return html.NewProxyCellProcessor(pProcessingFunc)
}

func r(pCellProcessors ...html.CellProcessor) html.RowProcessor {
	return html.NewSimpleRowProcessor(pCellProcessors...)
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
	return augmentor.Err(this.mCurSeason.SetFragmentRefAndID(slices.FirstOrEmpty(pCell.GetHrefs()), html.FirstTextOnly(pCell.GetText())),
		"SeasonID('%s')", pCell.GetText())
}

func (this *SeasonCollector) EpisodeCount(pCellText string) error {
	return augmentor.Err(this.mCurSeason.SetEpisodeCountFromText(html.FirstTextOnly(pCellText)), "EpisodeCount('%s')", pCellText)
}

func (this *SeasonCollector) SingleReleaseDate(pCellText string) error {
	zDate, err := utils.ExtractAirDate("SingleReleaseDate", pCellText)
	if err == nil {
		err = this.mCurSeason.SetSingleReleaseDate(zDate)
	}
	return augmentor.Err(err, "SingleReleaseDate('%s')", pCellText)
}

func (this *SeasonCollector) FirstAirDate(pCellText string) error {
	zDate, err := utils.ExtractAirDate("FirstAirDate", pCellText)
	if err == nil {
		err = this.mCurSeason.SetFirstAirDate(zDate)
	}
	return augmentor.Err(err, "FirstAirDate('%s')", pCellText)
}

func (this *SeasonCollector) SingleEpisodeDate(pCellText string) error {
	zDate, err := utils.ExtractAirDate("SingleEpisodeDate", pCellText)
	if err == nil {
		err = this.mCurSeason.SetSingleEpisodeDate(zDate)
	}
	return augmentor.Err(err, "SingleEpisodeDate('%s')", pCellText)
}

func (this *SeasonCollector) LastAirDate(pCellText string) error {
	zDate, err := utils.ExtractAirDate("LastAirDate", pCellText)
	if err == nil {
		err = this.mCurSeason.SetLastAirDate(zDate)
	}
	return augmentor.Err(err, "LastAirDate('%s')", pCellText)
}

func (this *SeasonCollector) newSingleRowProcessor(pCellProcessors ...html.CellProcessor) html.RowsProcessor {
	return this.newMultiRowProcessor(r(pCellProcessors...))
}

func (this *SeasonCollector) newMultiRowProcessor(pRowProcessors ...html.RowProcessor) html.RowsProcessor {
	return html.NewSimpleRowsProcessor(this.initSeason, this.finiSeason, pRowProcessors...)
}
