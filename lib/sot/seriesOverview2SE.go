package sot

import (
	html "svc-wiki-showepisodes/lib/htmlplus"
)

func init() {
	addFactoryMapping(factory2SE, "2SE-FL",
		html.HeaderRow{"Season", "Season", "Episodes", "Originally aired", "Originally aired"},
		html.HeaderRow{"Season", "Season", "Episodes", "First aired" /**/, "Last aired"})
	addFactoryMapping(factory2SE, "2SE-FL",
		html.HeaderRow{"Series", "Series", "Episodes", "Originally aired", "Originally aired"},
		html.HeaderRow{"Series", "Series", "Episodes", "First aired" /**/, "Last aired"})
	addFactoryMapping(factory2SE, "2SE-PF",
		html.HeaderRow{"Season", "Season", "Episodes", "Originally aired" /**/, "Originally aired"},
		html.HeaderRow{"Season", "Season", "Episodes", "Season premiere" /* */, "Season finale"})
	addFactoryMapping(factory2SE, "2SE-PF",
		html.HeaderRow{"Series", "Series", "Episodes", "Originally aired" /**/, "Originally aired"},
		html.HeaderRow{"Series", "Series", "Episodes", "Series premiere" /* */, "Series finale"})
	addFactoryMapping(factory2SE, "2SE-PF",
		html.HeaderRow{"Season", "Season", "Episodes", "Season premiere" /* */, "Season finale"})
	addFactoryMapping(factory2SE, "2SE-PF",
		html.HeaderRow{"Season", "Season", "Episodes", "Season Premiere" /* */, "Season Finale"})
	addFactoryMapping(factory2SE, "2SE-FL",
		html.HeaderRow{"Series", "Series", "Episodes", "Series premiere" /* */, "Series finale"})
	addFactoryMapping(factory2SE, "2SE-FL",
		html.HeaderRow{"Series", "Series", "Episodes", "Series Premiere" /* */, "Series Finale"})
}

func factory2SE(sc *SeasonCollector) *html.RowsProcessors {
	return html.NewRowsProcessors().
			Add(sc.newSingleRowProcessor(
		skipCell, cp(sc.SeasonID).Colspan(2), ctp(sc.SingleEpisodeDate).Colspan(2))).
			Add(sc.newSingleRowProcessor(
		skipCell, cp(sc.SeasonID), ctp(sc.EpisodeCount), ctp(sc.SingleReleaseDate).Colspan(2))).
			Add(sc.newSingleRowProcessor(html.S_IGNORED_CELL_PROCESSOR,
		cp(sc.SeasonID), ctp(sc.EpisodeCount), ctp(sc.FirstAirDate), ctp(sc.LastAirDate)))
}
