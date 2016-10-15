package page

import (
	html "svc-wiki-showepisodes/lib/htmlplus"
)

func init() {
	addSeriesOverview(process2SEFL, "2SEFL",
		html.HeaderRow{"Season", "Season", "Episodes", "Originally aired", "Originally aired"},
		html.HeaderRow{"Season", "Season", "Episodes", "First aired" /**/, "Last aired"})
	addSeriesOverview(process2SEFL, "2SEFL",
		html.HeaderRow{"Series", "Series", "Episodes", "Originally aired", "Originally aired"},
		html.HeaderRow{"Series", "Series", "Episodes", "First aired" /**/, "Last aired"})
}

var rowProcs2SEFL = newSOTrowProcessors().
		addDefault(newSimpleSOTrowProcessor(sSOTcellIgnored,
	sSOTcellSeasonNumber, sSOTcellEpisodeCount, sSOTcellFirstAirDate, sSOTcellLastAirDate))

func process2SEFL(pTable *html.Table) ([]*season, error) {
	return populateFromSOT(pTable, rowProcs2SEFL)
}
