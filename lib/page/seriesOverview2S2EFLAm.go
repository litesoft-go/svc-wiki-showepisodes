package page

import (
	html "svc-wiki-showepisodes/lib/htmlplus"
)

func init() {
	addSeriesOverview(process2S2EFLAm, "2S2EFLAm",
		html.HeaderRow{"Season", "Season", "Episodes", "Episodes", "Originally aired", "Originally aired", "Average viewers|||(in millions)"},
		html.HeaderRow{"Season", "Season", "Episodes", "Episodes", "First aired" /**/, "Last aired" /* */, "Average viewers|||(in millions)"})
}

func process2S2EFLAm(pTable *html.Table) ([]*season, error) {
	return populateFromSOT(pTable, newSOTrowProcessors().
			add(newSimpleSOTrowProcessor(sSOTcellIgnored, sSTOcellSeasonNumber, sSTOcellEpisodeCount.colspan(2), sSTOcellFirstAirDate, sSTOcellLastAirDate, sSOTcellIgnored)))
}

