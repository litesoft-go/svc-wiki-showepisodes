package sot

import (
	html "svc-wiki-showepisodes/lib/htmlplus"
)

func init() {
	addFactoryMapping(factory2SEFLN, "2SE-FL-N",
		html.HeaderRow{"Season", "Season", "Episodes", "Originally aired", "Originally aired", "Originally aired"},
		html.HeaderRow{"Season", "Season", "Episodes", "First aired" /**/, "Last aired" /* */, "Network"})
}

func factory2SEFLN(sc *SeasonCollector) *html.RowsProcessors {
	return html.NewRowsProcessors().
			With(removeNetwork).
			Add(sc.newSingleRowProcessor(
		skipCell, cp(sc.SeasonID), ctp(sc.EpisodeCount), ctp(sc.SingleReleaseDate).Colspan(2))).
			Add(sc.newSingleRowProcessor(
		skipCell, cp(sc.SeasonID), ctp(sc.EpisodeCount), ctp(sc.FirstAirDate), ctp(sc.LastAirDate)))
}
