package sot

import (
	html "svc-wiki-showepisodes/lib/htmlplus"
)

func init() {
	addFactoryMapping(factory2S2EOr, "2S2EOr",
		html.HeaderRow{"Season", "Season", "Episodes", "Episodes", "Originally released"})
}

func factory2S2EOr(sc *SeasonCollector) *html.RowsProcessors {
	return html.NewRowsProcessors().
			Add(sc.newSingleRowProcessor(
		skipCell, cp(sc.SeasonID), ctp(sc.EpisodeCount).Colspan(2), ctp(sc.SingleReleaseDate)))
}
