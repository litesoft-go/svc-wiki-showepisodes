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

func process2S2EFL(pTable *html.Table) (rSeasons []*season, err error) {
	fmt.Println("2S2EFL")
	return // todo: XXX
}

