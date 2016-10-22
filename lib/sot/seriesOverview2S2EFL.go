package sot

import (
	html "svc-wiki-showepisodes/lib/htmlplus"
)

func init() {
	addFactoryMapping(factory2S2E, "2S2E-FL",
		html.HeaderRow{"Season", "Season", "Episodes", "Episodes", "Originally aired", "Originally aired"},
		html.HeaderRow{"Season", "Season", "Episodes", "Episodes", "First aired" /**/, "Last aired"})
	addFactoryMapping(factory2S2E, "2S2E-PF",
		html.HeaderRow{"Season", "Season", "Episodes", "Episodes", "Season Premiere", "Season Finale"})
}

func factory2S2E(sc *SeasonCollector) *html.RowsProcessors {
	return html.NewRowsProcessors().
			Add(sc.newSingleRowProcessor(html.S_IGNORED_CELL_PROCESSOR,
		cp(sc.SeasonID), ctp(sc.EpisodeCount).Colspan(2), ctp(sc.FirstAirDate), ctp(sc.LastAirDate)))
}
