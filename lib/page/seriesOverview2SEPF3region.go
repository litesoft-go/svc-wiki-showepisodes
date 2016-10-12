package page

import (
	html "svc-wiki-showepisodes/lib/htmlplus"
)

func init() {
	addSeriesOverview(process2SEPF3region, "2SEPF3region",
		html.HeaderRow{"Season", "Season", "Episodes", "Originally aired" /**/, "Originally aired" /**/, "DVD|||and|||Blu-ray|||release dates", "DVD|||and|||Blu-ray|||release dates", "DVD|||and|||Blu-ray|||release dates"},
		html.HeaderRow{"Season", "Season", "Episodes", "Season premiere" /* */, "Season finale" /* . */, "Region 1" /* . . . . . . . . . . .*/, "Region 2" /* . . . . . . . . . . .*/, "Region 4"})
}

func process2SEPF3region(pTable *html.Table) ([]*season, error) {
	return populateFromSOT(pTable, newSOTrowProcessors().
			add(newSimpleSOTrowProcessor(sSOTcellIgnored, sSTOcellSeasonNumber, sSTOcellEpisodeCount, sSTOcellFirstAirDate, sSTOcellLastAirDate, sSOTcellIgnored, sSOTcellIgnored, sSOTcellIgnored)))
}

