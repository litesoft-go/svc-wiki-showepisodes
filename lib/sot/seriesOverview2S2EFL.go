package sot

import (
	html "svc-wiki-showepisodes/lib/htmlplus"
	"svc-wiki-showepisodes/lib/utils"
)

func init() {
	addProcessorMapping(processor2S2E, "2S2E-FL",
		html.HeaderRow{"Season", "Season", "Episodes", "Episodes", "Originally aired", "Originally aired"},
		html.HeaderRow{"Season", "Season", "Episodes", "Episodes", "First aired" /**/, "Last aired"})
	addProcessorMapping(processor2S2E, "2S2E-PF",
		html.HeaderRow{"Season", "Season", "Episodes", "Episodes", "Season Premiere", "Season Finale"})
}

var rowProcs2S2E = newRowProcessors().
		add(newSimpleSOTrowProcessor(sSOTcellIgnored,
	sSOTcellSeasonID, sSOTcellEpisodeCount.colspan(2), sSOTcellFirstAirDate, sSOTcellLastAirDate))

func processor2S2E(pTable *html.Table) ([]*utils.Season, error) {
	return populateFromSOT(pTable, rowProcs2S2E)
}
