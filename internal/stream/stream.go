package stream

import (
	"fmt"
	"sync"
	"time"

	"github.com/yuriongit/punch/internal/domain"
	"github.com/yuriongit/punch/internal/ui"
)

// PreTestLogs outputs initial benchmark headers.
func PreTestLogs(
	logChan chan<- string,
	testID *domain.TestID,
) {
	styledLeftBracket := ui.PrimaryFaintStyle.Render("[ ")
	styledRightBracket := ui.PrimaryFaintStyle.Render(" ]")

	logChan <- ui.ReverseLineBreak()
	logChan <- fmt.Sprintf(
		"%s%s Starting: %s%s",
		styledLeftBracket,
		ui.FaintStyle.Render("i."),
		ui.PrimaryStyle.Bold(true).Render(fmt.Sprintf("TEST-%s", *testID)),
		styledRightBracket,
	)
	logChan <- ui.LineBreak("normal")
}

// streamPostTestLogs outputs post test metrics.
func streamPostTestLogs(
	logChan chan<- string,
	m *domain.PostTestMetrics,
	c *domain.GlobalCounts,
) {
	/* Capture current time before date to acquire as accurate of a
	result as possible. */
	currentTime := time.Now().Format(time.TimeOnly)
	currentDate := ui.FormatDate(time.Now())

	testStatus := "Error"
	if c.RegularErr.Load() == 0 && c.FatalErr.Load() == 0 {
		testStatus = "Success"
	}

	// Create post test metrics.
	var (
		// Implement "Status" field's value. Will output "Error" or "Success"
		// based off error counts
		completedTestLog     = ui.CreatePostMetricLog("Status", testStatus)
		dateLog              = ui.CreatePostMetricLog("Date", currentDate)
		timeLog              = ui.CreatePostMetricLog("Time", currentTime)
		testIDLog            = ui.CreatePostMetricLog("Test ID", fmt.Sprintf("TEST-%s", *m.TestID))
		testDurationLog      = ui.CreatePostMetricLog("Duration", ui.FormatLatency(m.TestDuration))
		gracePeriodLog       = ui.CreatePostMetricLog("Grace Period", fmt.Sprintf("%d%s", m.GracePeriodPercent, "%"))
		fulfilledRequestsLog = ui.CreatePostMetricLog("Fulfilled", fmt.Sprintf("%d requests", c.Current.Load()))
		totalWorkersLog      = ui.CreatePostMetricLog("Workers", fmt.Sprintf("%d workers", c.Workers.Load()))
	)

	// Stream post test metrics.
	logChan <- ui.LineBreak("")
	logChan <- completedTestLog
	logChan <- testIDLog
	logChan <- testDurationLog
	logChan <- gracePeriodLog
	logChan <- fulfilledRequestsLog
	logChan <- totalWorkersLog
	logChan <- dateLog
	logChan <- timeLog
	logChan <- ui.LineBreak("exit")
}

// TestLogs streams logs to stdout concurrently.
func TestLogs(
	logChan chan string,
	logWg *sync.WaitGroup,
) {
	defer logWg.Done()
	for log := range logChan {
		fmt.Println(log)
	}
}
