package page

import (
	html "svc-wiki-showepisodes/lib/htmlplus"
)

func init() {
	addSeriesOverview(process2S2EFL, "2S2EFL",
		html.HeaderRow{"Season", "Season", "Episodes", "Episodes", "Originally aired", "Originally aired"},
		html.HeaderRow{"Season", "Season", "Episodes", "Episodes", "First aired" /**/, "Last aired"})
}

var rowProcs2S2EFL = newSOTrowProcessors().
		addDefault(newSimpleSOTrowProcessor(sSOTcellIgnored,
	sSOTcellSeasonNumber, sSOTcellEpisodeCount.colspan(2), sSOTcellFirstAirDate, sSOTcellLastAirDate))

func process2S2EFL(pTable *html.Table) ([]*season, error) {
	return populateFromSOT(pTable, rowProcs2S2EFL)
}
