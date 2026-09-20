package ui

import (
	lg "charm.land/lipgloss/v2"
	"github.com/yuriongit/punch/internal/domain"
)

// Lipgloss Styles
var ()

var (
	// Sub-log styles; for sub-log tree branches and labels
	SubLogStyle = lg.
			NewStyle().
			Italic(true).
			Foreground(lg.
				Color(SubLogColor)).
			Faint(true)

	// Line breaks for program entry and exit
	PrimaryStyle = lg.
			NewStyle().
			Foreground(lg.Color(PrimaryColor))

	PrimaryFaintStyle = lg.
				NewStyle().
				Foreground(lg.Color(PrimaryColor))

	LineBreakStyle = lg.
			NewStyle().
			Foreground(lg.Color(PrimaryColor)).
			Bold(true)

	// Fainted
	FaintStyle          = lg.NewStyle().Faint(true)
	LineBreakFaintStyle = LineBreakStyle.Faint(true).Padding(0)

	// Primary log Level styles
	SuccessStyle    = lg.NewStyle().Foreground(lg.Color(SuccessColor)).Bold(true)
	RegularErrStyle = lg.NewStyle().Foreground(lg.Color(RegularErrColor)).Bold(true)
	FatalErrStyle   = lg.NewStyle().Foreground(lg.Color(FatalErrColor)).Bold(true)

	// Accent log Level styles
	SuccessAccentStyle  = lg.NewStyle().Foreground(lg.Color(SuccessAccentColor)).Faint(true)
	RegErrAccentStyle   = lg.NewStyle().Foreground(lg.Color("#E14C4C")).Faint(true)
	FatalErrAccentStyle = lg.NewStyle().Foreground(lg.Color(FatalErrAccentColor)).Bold(true)
)

// Colors
const (
	PrimaryColor        = "#8F45FF"
	SubLogColor         = ""
	SuccessColor        = "#1EAD17"
	SuccessAccentColor  = "#1D8C18"
	RegularErrColor     = "#D42424"
	FatalErrColor       = "#FF942E"
	FatalErrAccentColor = "#FFAB5C"
)

// styleLogLvl returns the styled log level.
func LogLevel(lvl string) lg.Style {
	switch lvl {
	case "LOG-SUCC":
		return SuccessStyle
	case "LOG-ERRO":
		return RegularErrStyle
	default:
		return FatalErrStyle
	}
}

// styleLogLvlAccent returns styles with an accent style based on log level.
func LogLevelAccent(lvl string) lg.Style {
	switch lvl {
	case "LOG-SUCC":
		return SuccessAccentStyle
	case "LOG-ERRO":
		return RegErrAccentStyle
	default:
		return FatalErrAccentStyle
	}
}

// Timestamp styles the timestamp in worker logs.
func Timestamp(timestamp string, logLvl domain.LogLevel) string {
	return LogLevelAccent(logLvl.Load()).
		Faint(true).Render(timestamp)
}
