package page

import (
	client "lib-builtin/lib/clientadaptor"
	html "svc-wiki-showepisodes/lib/htmlplus"

	"lib-builtin/lib/augmentor"
	"lib-builtin/lib/ascii"
	"lib-builtin/lib/lines"

	"github.com/pkg/errors"

	"svc-wiki-showepisodes/lib/utils"
	"svc-wiki-showepisodes/lib/tv"

	"io/ioutil"
	"strings"
	"os"
	"svc-wiki-showepisodes/lib/sot"
)

type ShowEpisodes struct {
	mTitle       string `json:"title"`
	mEpisodesUrl string `json:"episodesUrl"`
	mSeasons     []*utils.Season `json:"seasons"`
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
	zCollector.Line(this.mTitle + "  (\"" + this.mEpisodesUrl + "\")")
	for _, zSeason := range this.mSeasons {
		zCollector.Indent()
		zSeason.AppendTo(zCollector)
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

func (this *ShowEpisodes) getSeasonsBasedOnDirectSeasonTables(pDocument *html.Document) (rSeasons []*utils.Season, err error) {
	rSeasons = append(rSeasons, &utils.Season{})
	//err = errors.New("niy: getSeasonsBasedOnDirectSeasonTables") // TODO: XXX
	return
}

func (this *ShowEpisodes) getSeasonsBasedOnSeriesOverviewTable(pDocument *html.Document) (rSeasons []*utils.Season, err error) {
	zTable, err := pDocument.GetTableWithId("Series_overview")
	if err != nil {
		err = html.ClearNodeNotFound(err)
	} else {
		rSeasons, err = sot.Process(zTable)
		err = augmentor.Err(err, "'Series_overview' table")
	}
	return
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
