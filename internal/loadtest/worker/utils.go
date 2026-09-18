/*
Package worker is responsible for carrying out the functionality
of sending HTTP requests to the configured 'target' and 'children'
in the Punch configuration file.
relies on. Workers are able to make HTTP requests both in parallel
and with concurrency.
*/
package worker

import (
	"fmt"
	"time"

	lg "charm.land/lipgloss/v2"
)

/*
----------------------
Stream-related helpers
----------------------
*/

// streamPreTestLogs outputs initial benchmark headers.
func streamPreTestLogs(
	logChan chan<- string,
	testID *string,
) {
	styledLeftBracket := primaryFaintStyle.Render("[ ")
	styledRightBracket := primaryFaintStyle.Render(" ]")

	logChan <- reverseLineBreak()
	logChan <- lg.NewStyle().Render(fmt.Sprintf(
		"%s%s Starting: %s%s",
		styledLeftBracket,
		faintStyle.Render("i."),
		primaryStyle.Bold(true).Render(fmt.Sprintf("TEST-%s", *testID)),
		styledRightBracket,
	),
	)
	logChan <- lineBreak("normal")
}

// streamPostTestLogs outputs post test metrics.
func streamPostTestLogs(
	logChan chan<- string,
	d *StreamPostTestData,
	c *GlobalCounts,
) {
	/* Capture current time before date to acquire as accurate of a
	result as possible. */
	currentTime := time.Now().Format(time.TimeOnly)
	currentDate := time.Now().Format(dateFormat)

	// Create post test metrics.
	var (
		// Implement "Status" field's value. Will output "Error" or "Success"
		// based off error counts
		completedTestLog     = createPostMetricLog("Status", "UNIMPLEMENTED")
		dateLog              = createPostMetricLog("Date", currentDate)
		timeLog              = createPostMetricLog("Time", currentTime)
		testIDLog            = createPostMetricLog("Test ID", fmt.Sprintf("TEST-%s", *d.TestID))
		testDurationLog      = createPostMetricLog("Duration", formatLatency(d.TestDur))
		gracePeriodLog       = createPostMetricLog("Grace Period", fmt.Sprintf("%d%s", d.GracePeriodPercent, "%"))
		fulfilledRequestsLog = createPostMetricLog("Fulfilled", fmt.Sprintf("%d requests", c.Curr.Load()))
		totalWorkersLog      = createPostMetricLog("Workers", fmt.Sprintf("%d workers", c.Workers.Load()))
	)

	// Stream post test metrics.
	logChan <- lineBreak("")
	logChan <- completedTestLog
	logChan <- testIDLog
	logChan <- testDurationLog
	logChan <- gracePeriodLog
	logChan <- fulfilledRequestsLog
	logChan <- totalWorkersLog
	logChan <- dateLog
	logChan <- timeLog
	logChan <- lineBreak("exit")
}

/*
----------------------
Stream-related helpers
----------------------
*/
