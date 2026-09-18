package worker

import (
	"fmt"
	"strings"
	"time"

	lg "charm.land/lipgloss/v2"
)

/*
----------------------
Constants
----------------------
*/

// Colors
const (
	primaryColor        = "#8F45FF"
	subLogColor         = ""
	successColor        = "#1EAD17"
	successAccentColor  = "#1D8C18"
	regularErrColor     = "#D42424"
	fatalErrColor       = "#FF942E"
	fatalErrAccentColor = "#FFAB5C"
)

/*
----------------------
Constants
----------------------
*/

/*
----------------------
Styles
----------------------
*/

// Lipgloss Styles
var (
	// Line breaks for program entry and exit
	primaryStyle = lg.
			NewStyle().
			Foreground(lg.Color(primaryColor))

	//
	primaryFaintStyle = lg.
				NewStyle().
				Foreground(lg.Color(primaryColor))

	lineBreakStyle = lg.
			NewStyle().
			Foreground(lg.Color(primaryColor)).
			Bold(true)

	// Fainted
	faintStyle          = lg.NewStyle().Faint(true)
	lineBreakFaintStyle = lineBreakStyle.Faint(true).Padding(0)

	// Sub-log styles; for sub-log tree branches and labels
	subLogStyle = lg.
			NewStyle().
			Italic(true).
			Foreground(lg.
				Color(subLogColor)).
			Faint(true)

	// Primary log Level styles
	successStyle    = lg.NewStyle().Foreground(lg.Color(successColor)).Bold(true)
	regularErrStyle = lg.NewStyle().Foreground(lg.Color(regularErrColor)).Bold(true)
	fatalErrStyle   = lg.NewStyle().Foreground(lg.Color(fatalErrColor)).Bold(true)

	// Accent log Level styles
	successAccentStyle  = lg.NewStyle().Foreground(lg.Color(successAccentColor)).Faint(true)
	regErrAccentStyle   = lg.NewStyle().Foreground(lg.Color("#E14C4C")).Faint(true)
	fatalErrAccentStyle = lg.NewStyle().Foreground(lg.Color(fatalErrAccentColor)).Bold(true)
)

/*
----------------------
Styles
----------------------
*/

/*
----------------------
Style-related helpers
----------------------
*/

// styleLogLvl returns the styled log level.
func styleLogLvl(lvl string) lg.Style {
	switch lvl {
	case "LOG-SUC":
		return successStyle
	case "LOG-ERR":
		return regularErrStyle
	default:
		return fatalErrStyle
	}
}

// styleLogLvlAccent returns strings with an accent style based on log level.
func styleLogLvlAccent(lvl string) lg.Style {
	switch lvl {
	case "LOG-SUC":
		return successAccentStyle
	case "LOG-ERR":
		return regErrAccentStyle
	default:
		return fatalErrAccentStyle
	}
}

// stylePostTestMetricField styles the metrics fields for the post test overview.
func stylePostTestMetricField(field string) string {
	if field == "Status" {
		return primaryStyle.Bold(true).Render("" + field + ":   ")
	}
	return lg.NewStyle().Faint(true).Italic(true).Render(" • " + field + ":")
}

// styleAndFormatTimestamp styles the timestamp in worker logs.
func styleAndFormatTimestamp(styledLogLvlAccent lg.Style, timestamp time.Time) string {
	return styledLogLvlAccent.Faint(true).Render(timestamp.Format(timeFormat))
}

/*
createPostMetricLog creates each worker's log. It provides a styled and
formatted presentation for the header, fields and values of the sub-logs.
*/
func createPostMetricLog(field string, value string) string {
	styledField := stylePostTestMetricField(field)

	switch value {
	case "Error":
		value = lg.NewStyle().Foreground(lg.Color(regularErrColor)).Render(value)
	case "Success":
		value = lg.NewStyle().Foreground(lg.Color(successColor)).Render(value)
	}

	var (
		maxFieldWidth = 20 // Width for the field column
	)

	// Calculate padding for the field (based on original field length, not styled length)
	fieldPadding := maxFieldWidth - len(field)
	if fieldPadding < 1 {
		fieldPadding = 1
	}

	return fmt.Sprintf(
		"%s%s%s",
		styledField,
		strings.Repeat(" ", fieldPadding),
		value,
	)
}

func reverseLineBreak() string {
	return lineBreakStyle.Render("Punch—————————————————————————————————————————————")
}

func lineBreak(lineType string) string {
	switch lineType {
	case "short":
		return lineBreakFaintStyle.Render("—————————")
	case "exit":
		return lineBreakStyle.Render("—————————————————————————————————————————————Punch")
	default:
		return lineBreakFaintStyle.Render("——————————————————————————————————————————————————")
	}
}

/*
----------------------
Style-related helpers
----------------------
*/
