package page

import (
	html "svc-wiki-showepisodes/lib/htmlplus"
	"fmt"
)

func init() {
	addSeriesOverview(process2S2EFL,
		html.HeaderRow{"Season", "Season", "Episodes", "Episodes", "Originally aired", "Originally aired"},
		html.HeaderRow{"Season", "Season", "Episodes", "Episodes", "First aired", "Last aired"})
}

func process2S2EFL(pTable *html.Table) ([]*season, error) {
	fmt.Println("2S2EFL")
	return populateFromSOT(pTable, newSOTrowProcessors().
		add(newSimpleSOTrowProcessor(sSOTcellIgnored, sSTOcellSeasonNumber, sSTOcellEpisodeCount.colspan(2), sSTOcellFirstAirDate, sSTOcellLastAirDate))) // ,
	// SOTcellSet(sSOTcellIgnored, sSTOcellSeasonNumber, sSTOcellEpisodeCount.colspan(2), sSTOcellFirstAirDate, sSTOcellLastAirDate))
}

