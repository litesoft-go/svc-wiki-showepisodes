package utils

import (
	html "svc-wiki-showepisodes/lib/htmlplus"

	"svc-wiki-showepisodes/lib/tv"

	"lib-builtin/lib/lines"

	"strconv"
	"fmt"
)

type SetSeasonAttributeFromText func(*Season, string) error

func SetSeasonID(pSeason *Season, pCellText string) (err error) {
	pSeason.mID = pCellText
	return
}

func SetEpisodeCount(pSeason *Season, pCellText string) (err error) {
	zText := html.FirstTextOnly(pCellText)
	if zText == "TBA" {
		pSeason.mEpisodesTBA = true
	} else {
		pSeason.mEpisodeCount, err = strconv.Atoi(zText)
	}
	return
}

func SetFirstAirDate(pSeason *Season, pCellText string) (err error) {
	pSeason.mFirstAirDate, err = extractAirDate("FirstAirDate", pCellText)
	return
}

func SetLastAirDate(pSeason *Season, pCellText string) (err error) {
	pSeason.mLastAirDate, err = extractAirDate("LastAirDate", pCellText)
	return
}

type Season struct {
	mID           string `json:"ID"`
	mEpisodeCount int `json:"episodeCount"`
	mEpisodesTBA  bool `json:"episodesTBA"`
	mFirstAirDate string `json:"firstAirDate"`
	mLastAirDate  string `json:"lastAirDate"`
	mFragmentRef  string `json:"fragmentRef"`
	mEpisodes     []*Episode `json:"episodes"`
}

func (this *Season) GetID() string {
	return this.mID
}

func (this *Season) GetEpisodeCount() int {
	return this.mEpisodeCount
}

// ISO8601 format Date or "" if N/A
func (this *Season) GetFirstAirDate() string {
	return this.mFirstAirDate
}

// ISO8601 format Date or "" if N/A
func (this *Season) GetLastAirDate() string {
	return this.mLastAirDate
}

func (this *Season) GetEpisodes() (rEpisodes []tv.Episode) {
	if len(this.mEpisodes) > 0 {
		rEpisodes = make([]tv.Episode, len(this.mEpisodes))
		for i, zEpisode := range this.mEpisodes {
			rEpisodes[i] = tv.Episode(zEpisode)
		}
	}
	return
}

func (this *Season) String() string {
	zCollector := lines.NewCollector()
	this.AppendTo(zCollector)
	return this.String()
}

func (this *Season) AppendTo(pLines *lines.Collector) {
	zLine := fmt.Sprintf("%-10s - %-10s %s: %2s", this.mFirstAirDate, this.mLastAirDate, this.episodeAppendValue(), this.mID)
	if this.mFragmentRef != "" {
		zLine += " | " + this.mFragmentRef
	}
	pLines.Line(zLine)
	for _, zEpisode := range this.mEpisodes {
		pLines.Indent()
		pLines.Line(zEpisode.String())
		pLines.Outdent()
	}
}

func (this *Season) episodeAppendValue() string {
	if this.mEpisodesTBA {
		return "(TBA)"
	}
	if this.mEpisodeCount != 0 {
		return fmt.Sprintf("  %2d ", this.mEpisodeCount)
	}
	return "(???)"
}

type Episode struct {
	mNumber  int `json:"number"`
	mTitle   string `json:"title"`
	mAirDate string `json:"airDate"`
}

func (this *Episode) GetNumber() int {
	return this.mNumber
}

func (this *Episode) GetTitle() string {
	return this.mTitle
}

// ISO8601 format Date or "" if N/A
func (this *Episode) GetAirDate() string {
	return this.mAirDate
}

func (this *Episode) String() string {
	return fmt.Sprintf("%2d - %10s : %s", this.mNumber, this.mAirDate, this.mTitle)
}
