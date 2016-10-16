package sot

import (
	html "svc-wiki-showepisodes/lib/htmlplus"
	"svc-wiki-showepisodes/lib/utils"
)

func init() {
	addProcessorMapping(processor2SE, "2SE-FL",
		html.HeaderRow{"Season", "Season", "Episodes", "Originally aired", "Originally aired"},
		html.HeaderRow{"Season", "Season", "Episodes", "First aired" /**/, "Last aired"})
	addProcessorMapping(processor2SE, "2SE-FL",
		html.HeaderRow{"Series", "Series", "Episodes", "Originally aired", "Originally aired"},
		html.HeaderRow{"Series", "Series", "Episodes", "First aired" /**/, "Last aired"})
	addProcessorMapping(processor2SE, "2SE-PF",
		html.HeaderRow{"Season", "Season", "Episodes", "Originally aired" /**/, "Originally aired"},
		html.HeaderRow{"Season", "Season", "Episodes", "Season premiere" /* */, "Season finale"})
	addProcessorMapping(processor2SE, "2SE-PF",
		html.HeaderRow{"Series", "Series", "Episodes", "Originally aired" /**/, "Originally aired"},
		html.HeaderRow{"Series", "Series", "Episodes", "Series premiere" /* */, "Series finale"})
	addProcessorMapping(processor2SE, "2SE-PF",
		html.HeaderRow{"Season", "Season", "Episodes", "Season premiere" /* */, "Season finale"})
	addProcessorMapping(processor2SE, "2SE-PF",
		html.HeaderRow{"Season", "Season", "Episodes", "Season Premiere" /* */, "Season Finale"})
	addProcessorMapping(processor2SE, "2SE-FL",
		html.HeaderRow{"Series", "Series", "Episodes", "Series premiere" /* */, "Series finale"})
	addProcessorMapping(processor2SE, "2SE-FL",
		html.HeaderRow{"Series", "Series", "Episodes", "Series Premiere" /* */, "Series Finale"})
}

var rowProcs2SE = newRowProcessors().
		add(newSimpleSOTrowProcessor(sSOTcellIgnored,
	sSOTcellSeasonID, sSOTcellEpisodeCount, sSOTcellFirstAirDate, sSOTcellLastAirDate))

func processor2SE(pTable *html.Table) ([]*utils.Season, error) {
	return populateFromSOT(pTable, rowProcs2SE)
}
