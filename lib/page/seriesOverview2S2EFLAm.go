package page

import (
	html "svc-wiki-showepisodes/lib/htmlplus"
	"fmt"
)

func init() {
	addSeriesOverview(process2S2EFLAm,
		html.HeaderRow{"Season", "Season", "Episodes", "Episodes", "Originally aired", "Originally aired", "Nielsen ratings", "Nielsen ratings"},
		html.HeaderRow{"Season", "Season", "Episodes", "Episodes", "First aired", "Last aired", "Rank", "Average viewers|||(in millions)"})
}

func process2S2EFLAm(pTable *html.Table) ([]*season, error) {
	fmt.Println("2S2EFLAm")
	return populateFromSOT(pTable,newSOTrowProcessors().
		add(newSimpleSOTrowProcessor(sSOTcellIgnored, sSTOcellSeasonNumber, sSTOcellEpisodeCount.colspan(2), sSTOcellFirstAirDate, sSTOcellLastAirDate, sSOTcellIgnored, sSOTcellIgnored)))
}

