package utils

import (
	html "svc-wiki-showepisodes/lib/htmlplus"

	"lib-builtin/lib/augmentor"
	"lib-builtin/lib/iso8601"
	"lib-builtin/lib/slices"
	"lib-builtin/lib/dates"

	"strings"
)

var JUNK_DATES = []string{"()"}

// FirstAirDate: October 10, 2012||| (|||2012-10-10|||)
// LastAirDate:      May 15, 2013||| (|||2013-05-15|||)
// FirstAirDate:  October 9, 2013||| (|||2013-10-09|||)
// LastAirDate:      May 14, 2014||| (|||2014-05-14|||)
// FirstAirDate:  October 8, 2014||| (|||2014-10-08|||)
// LastAirDate:      May 13, 2015||| (|||2015-05-13|||)
// FirstAirDate:  October 7, 2015||| (|||2015-10-07|||)
// LastAirDate:      May 25, 2016||| (|||2016-05-25|||)
// FirstAirDate:  October 5, 2016||| (|||2016-10-05|||)|||[2]
// LastAirDate: data-sort-TBA
// FirstAirDate:             2017||| (|||2017|||)|||[2]
//
// "December 25, 2012" OR "15 January 2012"
//
// Return: ISO8601 format Date or "" if N/A
func ExtractAirDate(pWhat, pCellText string) (rDate string, err error) {
	if strings.Contains(pCellText, "TBA") || strings.Contains(pCellText, "TBD") || (pCellText == ""){
		return
	}
	if strings.Contains(pCellText, html.CELL_TEXT_SEPARATOR) {
		rDate, err = extractMultiTextAirDate(pCellText)
	} else {
		rDate, err = extractSingleTextAirDate(pCellText)
	}
	err = augmentor.Err(err, pWhat)
	return
}

// Return: ISO8601 format Date or "" if N/A
func extractMultiTextAirDate(pCellText string) (rDate string, err error) {
	zDate := html.FirstTextOnly(pCellText)
	if len(zDate) == 4 {
		rDate, err = iso8601.ValidateToYear(zDate)
		if err == nil {
			return
		}
	}
	zDate, err = html.NthTextOnly(pCellText, 2)
	if err == nil {
		if len(zDate) == 4 {
			rDate, err = iso8601.ValidateToYear(zDate)
		} else {
			rDate, err = iso8601.ValidateToDay(zDate)
		}
		err = augmentor.Err(err, "'%s'", pCellText)
	}
	return
}

// "December 25, 2012" OR "15 January 2012"
//
// Return: ISO8601 format Date or "" if N/A
func extractSingleTextAirDate(pDate string) (rDate string, err error) {
	pDate = strings.TrimSpace(pDate)
	if !slices.Contains(JUNK_DATES, pDate) {
		if len(pDate) == 4 {
			rDate, err = iso8601.ValidateToYear(pDate)
		} else {
			rDate, err = fmtISO8601(dates.ParseTextualMonthData(pDate))
		}
	}
	return
}

func fmtISO8601(pYear, pMonth, pDay int, pError error) (rDate string, err error) {
	err = pError
	if err == nil {
		rDate, err = iso8601.Format(pYear, pMonth, pDay)
	}
	return
}
