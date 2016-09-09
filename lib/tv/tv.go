package tv

import "fmt"

type Episode interface {
	GetNumber() int
	GetTitle() string
	GetAirDate() string // ISO8601 format Date or "" if N/A
	fmt.Stringer
}

type Season interface {
	GetNumber() int
	GetEpisodeCount() int
	GetFirstAirDate() string // ISO8601 format Date or "" if N/A
	GetLastAirDate() string // ISO8601 format Date or "" if N/A
	GetEpisodes() []Episode
	fmt.Stringer
}

type Show interface {
	GetTitle() string
	GetSeasons() []Season
	fmt.Stringer
}
