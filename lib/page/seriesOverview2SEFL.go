package page

import (
	html "svc-wiki-showepisodes/lib/htmlplus"
)

func init() {
	addSeriesOverview(process2SEFL, "2SEFL",
		html.HeaderRow{"Season", "Season", "Episodes", "Originally aired", "Originally aired"},
		html.HeaderRow{"Season", "Season", "Episodes", "First aired" /**/, "Last aired"})
}

func process2SEFL(pTable *html.Table) ([]*season, error) {
	return populateFromSOT(pTable, newSOTrowProcessors().
			add(newSimpleSOTrowProcessor(sSOTcellIgnored, sSTOcellSeasonNumber, sSTOcellEpisodeCount, sSTOcellFirstAirDate, sSTOcellLastAirDate)))
}

