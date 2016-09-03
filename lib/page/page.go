package page

import (
	client "lib-builtin/lib/clientadaptor"

	"lib-builtin/lib/ascii"

	"github.com/pkg/errors"

	"strings"
	"fmt"
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

func (this *ShowEpisodes) PullAndParse() error {
	zEndpoint, err := client.NewDomainEndpoint().WithFullUrl("https://en.wikipedia.org/wiki/List_of_Major_Crimes_episodes").
		WithStatusFilter(client.Non200StatusFilter).UsingGet()
	if err != nil {
		return err
	}
	zBody, _, err := zEndpoint.GetResponseBodyAndCode()
	return this.Parse([]byte(zBody), err)
}

func (this *ShowEpisodes) Parse(pBody []byte, pError error) (err error) {
	err = pError
	fmt.Println("Body len:", len(pBody))
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

type episode struct {
	mNumber  int `json:"number"`
	mTitle   string `json:"title"`
	mAirDate string `json:"airDate"`
}

type season struct {
	mNumber       int `json:"number"`
	mEpisodeCount int `json:"episodeCount"`
	mFirstAirDate string `json:"firstAirDate"`
	mLastAirDate  string `json:"lastAirDate"`
	mEpisodes     []*episode `json:"episodes"`
}


