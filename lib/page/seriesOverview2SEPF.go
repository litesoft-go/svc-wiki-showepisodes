package page

import (
	html "svc-wiki-showepisodes/lib/htmlplus"
)

func init() {
	addSeriesOverview(process2SEPF, "2SEPF",
		html.HeaderRow{"Season", "Season", "Episodes", "Originally aired" /**/, "Originally aired"},
		html.HeaderRow{"Season", "Season", "Episodes", "Season premiere" /* */, "Season finale"})
	addSeriesOverview(process2SEPF, "2SEPF",
		html.HeaderRow{"Series", "Series", "Episodes", "Originally aired" /**/, "Originally aired"},
		html.HeaderRow{"Series", "Series", "Episodes", "Series premiere" /* */, "Series finale"})
	addSeriesOverview(process2SEPF, "2SEPF",
		html.HeaderRow{"Season", "Season", "Episodes", "Season premiere" /* */, "Season finale"})
	addSeriesOverview(process2SEPF, "2SEPF",
		html.HeaderRow{"Season", "Season", "Episodes", "Season Premiere" /* */, "Season Finale"})
	addSeriesOverview(process2SEFL, "2SEFL",
		html.HeaderRow{"Series", "Series", "Episodes", "Series premiere" /* */, "Series finale"})
	addSeriesOverview(process2SEFL, "2SEFL",
		html.HeaderRow{"Series", "Series", "Episodes", "Series Premiere" /* */, "Series Finale"})
}

var rowProcs2SEPF = newSOTrowProcessors().
		addDefault(newSimpleSOTrowProcessor(sSOTcellIgnored,
	sSOTcellSeasonNumber, sSOTcellEpisodeCount, sSOTcellFirstAirDate, sSOTcellLastAirDate))

func process2SEPF(pTable *html.Table) ([]*season, error) {
	return populateFromSOT(pTable, rowProcs2SEPF)
}