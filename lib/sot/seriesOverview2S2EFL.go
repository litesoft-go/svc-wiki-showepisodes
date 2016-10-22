package sot

import (
	html "svc-wiki-showepisodes/lib/htmlplus"
)

func init() {
	addFactoryMapping(factory2S2E, "2S2E-FL",
		html.HeaderRow{"Series", "Series", "Episodes", "Episodes", "Originally aired", "Originally aired"},
		html.HeaderRow{"Series", "Series", "Episodes", "Episodes", "First aired" /**/, "Last aired"})
	addFactoryMapping(factory2S2E, "2S2E-FL",
		html.HeaderRow{"Season", "Season", "Episodes", "Episodes", "Originally aired", "Originally aired"},
		html.HeaderRow{"Season", "Season", "Episodes", "Episodes", "First aired" /**/, "Last aired"})
	addFactoryMapping(factory2S2E, "2S2E-PF",
		html.HeaderRow{"Season", "Season", "Episodes", "Episodes", "Season Premiere", "Season Finale"})
}

func factory2S2E(sc *SeasonCollector) *html.RowsProcessors {
	return html.NewRowsProcessors().
			Add(sc.newMultiRowProcessor(
		r(skipCell.Rowspan(2), cp(sc.SeasonID).Rowspan(2), ctp(sc.EpisodeCount).Rowspan(2), skipCell, ctp(sc.FirstAirDate)),
		r(skipCell, skipCell, ctp(sc.LastAirDate)))).
			Add(sc.newMultiRowProcessor(
		r(skipCell, cp(sc.SeasonID).Rowspan(2), ctp(sc.EpisodeCount).Rowspan(2), skipCell, ctp(sc.FirstAirDate)),
		r(skipCell, skipCell, skipCell, ctp(sc.LastAirDate)))).
			Add(sc.newSingleRowProcessor(
		skipCell, cp(sc.SeasonID).Colspan(3), ctp(sc.SingleEpisodeDate).Colspan(2))).
			Add(sc.newSingleRowProcessor(
		skipCell, cp(sc.SeasonID), ctp(sc.EpisodeCount).Colspan(2), ctp(sc.FirstAirDate), ctp(sc.LastAirDate)))
}
