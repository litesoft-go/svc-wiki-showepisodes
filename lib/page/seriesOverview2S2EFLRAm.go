package page

import (
	html "svc-wiki-showepisodes/lib/htmlplus"
)

func init() {
	addSeriesOverview(process2S2EFLRAm, "2S2EFLRAm",
		html.HeaderRow{"Season", "Season", "Episodes", "Episodes", "Originally aired", "Originally aired", "Nielsen ratings", "Nielsen ratings"},
		html.HeaderRow{"Season", "Season", "Episodes", "Episodes", "First aired" /**/, "Last aired" /* */, "Rank" /* . . .*/, "Average viewers|||(in millions)"})
}

func process2S2EFLRAm(pTable *html.Table) ([]*season, error) {
	return populateFromSOT(pTable, newSOTrowProcessors().
			add(newSimpleSOTrowProcessor(sSOTcellIgnored, sSTOcellSeasonNumber, sSTOcellEpisodeCount.colspan(2), sSTOcellFirstAirDate, sSTOcellLastAirDate, sSOTcellIgnored, sSOTcellIgnored)))
}

