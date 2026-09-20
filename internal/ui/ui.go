package ui

import (
	"fmt"
	"strings"
	"time"

	lg "charm.land/lipgloss/v2"
	"github.com/yuriongit/punch/internal/domain"
)

// Miscellaneous
const (
	// Formats
	timeFormat = "15:04:05"
	dateFormat = "2006-01-02"
)

func ReverseLineBreak() string {
	return LineBreakStyle.Render("Punch—————————————————————————————————————————————")
}

func LineBreak(lineType string) string {
	switch lineType {
	case "short":
		return LineBreakFaintStyle.Render("—————————")
	case "exit":
		return LineBreakStyle.Render("—————————————————————————————————————————————Punch")
	default:
		return LineBreakFaintStyle.Render("——————————————————————————————————————————————————")
	}
}

// FormatPostTestMetricField styles the metrics fields for the post test overview.
func FormatPostTestMetricField(field string) string {
	if field == "Status" {
		return PrimaryStyle.Bold(true).Render("" + field + ":   ")
	}
	return lg.NewStyle().Faint(true).Italic(true).Render(" • " + field + ":")
}

/*
formatLatency dynamically formats a time.Duration into
human-readable units (nanoseconds, microseconds, milliseconds,
seconds, minutes, hours).
*/
func FormatLatency(
	d time.Duration,
) string {
	switch {
	case d < time.Microsecond:
		ns := d.Nanoseconds()
		if ns == 1 {
			return "1 nanosecond"
		}
		return fmt.Sprintf("%d nanoseconds", ns)

	case d < time.Millisecond:
		us := float64(d.Nanoseconds()) / 1000.0
		if us == 1.0 {
			return "1 microsecond"
		}
		return fmt.Sprintf("%.2f microseconds", us)

	case d < time.Second:
		ms := float64(d.Nanoseconds()) / 1000000.0
		if ms == 1.0 {
			return "1 millisecond"
		}
		return fmt.Sprintf("%.2f milliseconds", ms)

	case d < time.Minute:
		s := d.Seconds()
		if s == 1.0 {
			return "1 second"
		}
		return fmt.Sprintf("%.2f seconds", s)

	case d < time.Hour:
		m := d.Minutes()
		if m == 1.0 {
			return "1 minute"
		}
		return fmt.Sprintf("%.2f minutes", m)

	default:
		h := d.Hours()
		if h == 1.0 {
			return "1 hour"
		}
		return fmt.Sprintf("%.2f hours", h)
	}
}

func FormatLogLevel(logLvl domain.LogLevel, style bool) string {
  
  return ""
}

/*
CreatePostMetricLog creates each worker's log. It provides a styled and
formatted presentation for the header, fields and values of the sub-logs.
*/
func CreatePostMetricLog(field string, value string) string {
	styledField := FormatPostTestMetricField(field)

	switch value {
	case "Error":
		value = lg.NewStyle().Foreground(lg.Color(RegularErrColor)).Render(value)
	case "Success":
		value = lg.NewStyle().Foreground(lg.Color(SuccessColor)).Render(value)
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

/*
----------------------
Style-related helpers
----------------------
*/

// TODO
func FormatTimestamp(timestamp time.Time, logLvl domain.LogLevel, style bool) string {
  time := timestamp.Format(timeFormat)
  if style {
    time = Timestamp(time, logLvl)
  }
  return time
}

func FormatDate(date time.Time) string {
  return date.Format(dateFormat)
}