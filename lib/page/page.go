package page

import (
	client "lib-builtin/lib/clientadaptor"
	html "svc-wiki-showepisodes/lib/htmlplus"

	"lib-builtin/lib/ascii"
	"lib-builtin/lib/lines"

	"github.com/pkg/errors"

	"svc-wiki-showepisodes/lib/tv"

	"io/ioutil"
	"strings"
	"fmt"
	"os"
)

type ShowEpisodes struct {
	mTitle       string `json:"title"`
	mEpisodesUrl string `json:"episodesUrl"`
	mSeasons     []*season `json:"seasons"`
}

func NewWikiShowEpisodes(pEpisodesUrl string) (rEpisodes *ShowEpisodes, err error) {
	zTitle, err := ExtractTitle(pEpisodesUrl)
	rEpisodes = &ShowEpisodes{mTitle:zTitle, mEpisodesUrl:pEpisodesUrl}
	return
}

func (this *ShowEpisodes) GetTitle() string {
	return this.mTitle
}

func (this *ShowEpisodes) GetSeasons() (rSeasons []tv.Season) {
	if len(this.mSeasons) > 0 {
		rSeasons = make([]tv.Season, len(this.mSeasons))
		for i, zSeason := range this.mSeasons {
			rSeasons[i] = tv.Season(zSeason)
		}
	}
	return
}
func (this *ShowEpisodes) String() string {
	zCollector := lines.NewCollector()
	zCollector.Line(this.mTitle + "  (" + this.mEpisodesUrl + ")")
	for _, zSeason := range this.mSeasons {
		zCollector.Indent()
		zSeason.appendTo(zCollector)
		zCollector.Outdent()
	}
	return zCollector.String()
}

func (this *ShowEpisodes) PullAndParse(pFileName string) error {
	zEndpoint, err := client.NewDomainEndpoint().WithFullUrl(this.mEpisodesUrl).
		WithStatusFilter(client.Non200StatusFilter).UsingGet()
	if err != nil {
		return err
	}
	zBody, _, err := zEndpoint.GetResponseBodyAndCode()
	if (pFileName != "") && (err == nil) {
		err = ioutil.WriteFile(pFileName, []byte(zBody), os.ModePerm)
	}
	return this.Parse([]byte(zBody), err)
}

func (this *ShowEpisodes) Parse(pBody []byte, pError error) (err error) {
	err = pError
	if err != nil {
		return
	}
	zDocument, err := html.NewDocument(pBody)
	if err != nil {
		return
	}
	zSeasons, err := this.getSeasonsBasedOnSeriesOverviewTable(zDocument)
	if (len(zSeasons) == 0) && (err == nil) {
		zSeasons, err = this.getSeasonsBasedOnDirectSeasonTables(zDocument)
	}
	if err == nil {
		this.mSeasons = zSeasons
		if len(zSeasons) == 0 {
			err = errors.New("no Seasons found")
		} else {
			err = this.fillEpisodes(zDocument)
		}
	}
	return
}

func (this *ShowEpisodes) fillEpisodes(pDocument *html.Document) (err error) {
	return // TODO: XXX
}

func (this *ShowEpisodes) getSeasonsBasedOnDirectSeasonTables(pDocument *html.Document) (rSeasons []*season, err error) {
	err = errors.New("niy: getSeasonsBasedOnDirectSeasonTables") // TODO: XXX
	return
}

func (this *ShowEpisodes) getSeasonsBasedOnSeriesOverviewTable(pDocument *html.Document) (rSeasons []*season, err error) {
	zTable, err := pDocument.GetTableWithId("Series_overview")
	if err != nil {
		err = html.ClearNodeNotFound(err)
		return
	}
	zProcessor, err := determineProcessorSOT(zTable)
	if err == nil {
		rSeasons, err = zProcessor(zTable)
	}
	return
}

func determineProcessorSOT(pTable *html.Table) (SeriesOverviewTableProcessor, error) {
	for _, zSet := range sSOT_sets {
		if pTable.HeaderMatches(zSet.mHeaderRows) {
			return zSet.mProcessor, nil
		}
	}
	return nil, pTable.ErrorHeaderNotMatched()
}

type SeriesOverviewTableProcessor func(pTable *html.Table) (rSeasons []*season, err error)

type SOT_set struct {
	mProcessor  SeriesOverviewTableProcessor
	mHeaderRows []html.HeaderRow
}

var sSOT_sets []SOT_set

func addSeriesOverview(pProcessor SeriesOverviewTableProcessor, pHeaderRows ...html.HeaderRow) {
	sSOT_sets = append(sSOT_sets, SOT_set{mProcessor:pProcessor, mHeaderRows:pHeaderRows})
}

const (
	LIST_OF_PREFIX = "List_of_"
	LIST_OF_SUFFIX = "_episodes"
)

// https://en.wikipedia.org/wiki/List_of_Major_Crimes_episodes#??? -> "Major Crimes"
func ExtractTitle(pEpisodesUrl string) (rTitle string, err error) {
	zEpisodeUrl := ascii.RemoveEverythingFrom(pEpisodesUrl, "#")
	zEpisodeUrl = ascii.RemoveEverythingFrom(zEpisodeUrl, "?")
	zLastPath, err := ascii.JustAfterLast(zEpisodeUrl, "/")
	if err == nil {
		rTitle, err = toTitle(pEpisodesUrl, zLastPath)
	}
	return
}

func toTitle(pEpisodesUrl, pLastPath string) (rTitle string, err error) {
	if strings.HasPrefix(pLastPath, LIST_OF_PREFIX) && strings.HasSuffix(pLastPath, LIST_OF_SUFFIX) {
		zFrom := len(LIST_OF_PREFIX)
		zUpTo := len(pLastPath) - len(LIST_OF_SUFFIX)
		if zUpTo <= zFrom {
			err = errors.Errorf("url ('%s') last path did not start with '%s' AND end with '%s', but was: %s", pEpisodesUrl, LIST_OF_PREFIX, LIST_OF_SUFFIX, pLastPath)
			return
		}
		pLastPath = pLastPath[zFrom:zUpTo]
	}
	rTitle = normalizeTitle(pLastPath)
	return
}

func normalizeTitle(pRawTitle string) (rTitle string) {
	rTitle = strings.Replace(pRawTitle, "_", " ", -1)
	rTitle = ascii.RemoveEverythingFrom(rTitle, " (TV series)")
	rTitle = ascii.ReplacePercentEscapedChars(rTitle)
	rTitle = strings.Replace(rTitle, ".", "", -1)
	return
}

type season struct {
	mNumber       int `json:"number"`
	mEpisodeCount int `json:"episodeCount"`
	mFirstAirDate string `json:"firstAirDate"`
	mLastAirDate  string `json:"lastAirDate"`
	mEpisodes     []*episode `json:"episodes"`
}

func (this *season) GetNumber() int {
	return this.mNumber
}

func (this *season) GetEpisodeCount() int {
	return this.mEpisodeCount
}

// ISO8601 format Date or "" if N/A
func (this *season) GetFirstAirDate() string {
	return this.mFirstAirDate
}

// ISO8601 format Date or "" if N/A
func (this *season) GetLastAirDate() string {
	return this.mLastAirDate
}

func (this *season) GetEpisodes() (rEpisodes []tv.Episode) {
	if len(this.mEpisodes) > 0 {
		rEpisodes = make([]tv.Episode, len(this.mEpisodes))
		for i, zEpisode := range this.mEpisodes {
			rEpisodes[i] = tv.Episode(zEpisode)
		}
	}
	return
}

func (this *season) String() string {
	zCollector := lines.NewCollector()
	this.appendTo(zCollector)
	return this.String()
}

func (this *season) appendTo(pLines *lines.Collector) {
	pLines.Line(fmt.Sprintf("%2d - %10s - %10s - %2d episodes", this.mNumber, this.mFirstAirDate, this.mLastAirDate, this.mEpisodeCount))
	for _, zEpisode := range this.mEpisodes {
		pLines.Indent()
		pLines.Line(zEpisode.String())
		pLines.Outdent()
	}
}

type episode struct {
	mNumber  int `json:"number"`
	mTitle   string `json:"title"`
	mAirDate string `json:"airDate"`
}

func (this *episode) GetNumber() int {
	return this.mNumber
}

func (this *episode) GetTitle() string {
	return this.mTitle
}

// ISO8601 format Date or "" if N/A
func (this *episode) GetAirDate() string {
	return this.mAirDate
}

func (this *episode) String() string {
	return fmt.Sprintf("%2d - %10s : %s", this.mNumber, this.mAirDate, this.mTitle)
}
