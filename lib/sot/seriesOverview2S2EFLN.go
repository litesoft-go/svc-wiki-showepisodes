package sot

import (
	html "svc-wiki-showepisodes/lib/htmlplus"
)

func init() {
	addFactoryMapping(factory2S2EFLN, "2S2E-FL-N",
		html.HeaderRow{"Season", "Season", "Episodes", "Episodes", "Originally aired", "Originally aired", "Originally aired"},
		html.HeaderRow{"Season", "Season", "Episodes", "Episodes", "First aired" /**/, "Last aired" /* */, "Network"})
}

func factory2S2EFLN(sc *SeasonCollector) *html.RowsProcessors {
	return html.NewRowsProcessors().
			With(removeNetwork).
			Add(sc.newSingleRowProcessor(
		skipCell, cp(sc.SeasonID), ctp(sc.EpisodeCount).Colspan(2), ctp(sc.SingleReleaseDate).Colspan(2))).
			Add(sc.newSingleRowProcessor(
		skipCell, cp(sc.SeasonID), ctp(sc.EpisodeCount).Colspan(2), ctp(sc.FirstAirDate), ctp(sc.LastAirDate)))
}
