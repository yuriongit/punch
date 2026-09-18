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
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/yuriongit/punch/internal/config"
)

/*
Global HTTP Client configured with timeouts and connection
pooling to prevent goroutine leaks and stalled requests during
high load tests.
*/
var httpClient = &http.Client{
	Timeout: 5 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        1000,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	},
}

/*
formatLatency dynamically formats a time.Duration into
human-readable units (nanoseconds, microseconds, milliseconds,
seconds, minutes, hours).
*/
func formatLatency(
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

/*
createWorkerLog creates a formatted worker log
with custom colored brackets and light grey sub-logs.
*/
func createWorkerLog(
	logChan chan<- string,
	l *Log,
	c *CreateWorkerResLogCounts,
	latency time.Duration,
) {
	logLvl := l.Lvl.Load()
	styledLogLvl := styleLogLvl(logLvl)
	styledTimestamp := styleAndFormatTimestamp(styledLogLvl.Bold(false), time.Now())
	formattedLatency := formatLatency(latency)

	styleLogLvlAccent := styleLogLvlAccent(logLvl)

	// 2. Colorize status level text
	coloredLvl := styledLogLvl.Render(logLvl)

	// 3. Highlight Got and Want status codes using the log level's style
	gotAndWantStatusCodes := styledLogLvl.Render(fmt.Sprintf("%d/%d", l.GotStatusCode, l.WantStatusCode))

	// 4. Accent-color the opening and closing brackets
	openBracket := styleLogLvlAccent.Render("[ ")
	closeBracket := styleLogLvlAccent.Render(" ]")

	if l.Error == "" { // TODO: Change to l.Response
		l.Error = "Punch – N/A"
	} else {
		l.Error = fmt.Sprintf("'%s'", l.Error)
	}

	header := fmt.Sprintf(
		"%s %s%s %s Child #%d, Worker #%d%s\n",
		styledTimestamp,
		openBracket,
		coloredLvl,
		faintStyle.Render("|"),
		l.ChildID,
		l.WorkerID,
		closeBracket,
	)

	// 5. Build raw tree structure
	rawSubLogs := fmt.Sprintf(
		"  %s %s\n"+
			"  ├── Global Request: #%d\n"+
			"  ├── My request: #%d\n"+
			"  ├── Latency: %s\n"+
			"  └── Response Body: %s",
		styleLogLvl(logLvl).Bold(true).Blink(true).Render("├── Got/Want:"),
		gotAndWantStatusCodes,
		c.GlobalReqCounter,
		c.Curr,
		formattedLatency,
		l.Error, // Change to l.ResponseBody
	)

	// 6. Render the sub-logs in light grey (248)
	coloredSubLogs := subLogStyle.Render(rawSubLogs)

	logChan <- (header + coloredSubLogs)
}

// StreamWorkerLogs streams logs to stdout concurrently.
func StreamWorkerLogs(
	logChan chan string,
	logWg *sync.WaitGroup,
) {
	defer logWg.Done()
	for msg := range logChan {
		fmt.Println(msg)
	}
}

func executeWorker(
	logChan chan<- string,
	wg *sync.WaitGroup,
	globalCounts *GlobalCounts,
	workerID uint32,
	childID uint32,
	req ReqInfo,
	targetRequests uint32,
	baseDuration uint16,
	workerReqDelay time.Duration,
) {
	defer wg.Done()

	if targetRequests == 0 {
		return
	}

	switch req.Method {
	case "POST":
		panic("not yet implemented")
	case "GET":
		baseDuration := time.Duration(baseDuration) * time.Second
		gracePeriod := (baseDuration * time.Duration(req.GracePeriodPercent)) / 100

		// Maximum time allotted to send requests
		maxTime := baseDuration + gracePeriod
		deadline := time.Now().Add(maxTime)

		// Request counters for each worker's child
		counts := ChildCounts{}

		for time.Now().Before(deadline) {
			// Pause to stretch requests over full BaseDuration
			time.Sleep(workerReqDelay)

			// Measure pure HTTP round-trip latency
			requestStartTime := time.Now()
			resp, err := httpClient.Get(req.URL)
			elapsedTime := time.Since(requestStartTime)

			counts.Curr++
			globalCounts.Curr.Add(1)

			switch {
			case err != nil:
				counts.FatErr++
				globalCounts.FatErr.Add(1)

				l := Log{
					Lvl:            LogFat,
					WorkerID:       workerID,
					ChildID:        childID,
					ReqMethod:      req.Method,
					WantStatusCode: req.WantStatusCode,
					GotStatusCode:  0,
					Error:          err.Error(),
				}
				c := CreateWorkerResLogCounts{
					GlobalReqCounter: globalCounts.Curr.Load(),
					Curr:             counts.Curr, // Rename to localCounts for improved clarity.
				}

				createWorkerLog(logChan, &l, &c, elapsedTime)

			case resp.StatusCode == int(req.WantStatusCode):
				counts.Suc++
				globalCounts.Suc.Add(1)

				l := Log{
					Lvl:            LogSuc,
					WorkerID:       workerID,
					ChildID:        childID,
					ReqMethod:      req.Method,
					WantStatusCode: req.WantStatusCode,
					GotStatusCode:  uint16(resp.StatusCode), //nolint:gosec // Status codes safely fit
				}
				c := CreateWorkerResLogCounts{
					GlobalReqCounter: globalCounts.Curr.Load(),
					Curr:             counts.Curr,
				}

				createWorkerLog(logChan, &l, &c, elapsedTime)
				_ = resp.Body.Close()

			default:
				counts.RegErr++
				globalCounts.RegErr.Add(1)

				l := Log{
					Lvl:            LogErr,
					WorkerID:       workerID,
					ChildID:        childID,
					ReqMethod:      req.Method,
					WantStatusCode: req.WantStatusCode,
					GotStatusCode:  uint16(resp.StatusCode), //nolint:gosec // Status codes safely fit
				}
				c := CreateWorkerResLogCounts{
					GlobalReqCounter: globalCounts.Curr.Load(),
					Curr:             counts.Curr,
				}

				createWorkerLog(logChan, &l, &c, elapsedTime)
				_ = resp.Body.Close()
			}

			// Stop when worker fulfills assigned quota
			if counts.Curr == targetRequests {
				return
			}
		}
	}
}

// RunWorkers ...
func RunWorkers(
	wg *sync.WaitGroup,
	logChan chan<- string,
	config *config.PunchConfig,
	testID *string,
) {
	streamPreTestLogs(logChan, testID)

	globalCounts := GlobalCounts{
		Curr:    atomic.Uint32{},
		Workers: atomic.Uint32{},
		FatErr:  atomic.Uint32{},
		RegErr:  atomic.Uint32{},
	}

	testStartTime := time.Now()

	for childID, child := range config.Children {
		if child.BaseDurationSecs == 0 || child.TotalRequests == 0 {
			continue
		}

		requestsPerSecond := float32(child.TotalRequests) / float32(child.BaseDurationSecs)
		numWorkers := uint32(requestsPerSecond)
		if numWorkers < 1 {
			numWorkers = 1
		}

		// Cap workers to total requests if total requests is smaller than numWorkers
		if numWorkers > child.TotalRequests {
			numWorkers = child.TotalRequests
		}

		// Calculate base distribution and remainders so no requests drop out
		baseReqsPerWorker := child.TotalRequests / numWorkers
		remainderReqs := child.TotalRequests % numWorkers

		if child.WantStatusCode == nil {
			childExpectedStatus := uint16(200)
			child.WantStatusCode = &childExpectedStatus
		}

		URL := fmt.Sprintf("%s://%s%s", config.Protocol, config.Target, child.Name)

		for i := uint32(0); i < numWorkers; i++ {
			reqsForThisWorker := baseReqsPerWorker
			if i < remainderReqs {
				reqsForThisWorker++
			}

			// Delay calculation per request to fill BaseDurationSecs properly
			var workerReqDelay time.Duration
			if reqsForThisWorker > 0 {
				workerReqDelay = time.Duration((float64(child.BaseDurationSecs) / float64(reqsForThisWorker)) * float64(time.Second))
			}

			globalCounts.Workers.Add(1)
			wg.Add(1)

			req := ReqInfo{
				Method:             child.Method,
				URL:                URL,
				ChildName:          child.Name,
				WantStatusCode:     *child.WantStatusCode,
				GracePeriodPercent: config.GracePeriodPercent,
			}

			workerID := globalCounts.Workers.Add(1)

			go executeWorker(
				logChan,
				wg,
				&globalCounts,
				workerID,
				uint32(childID+1),
				req,
				reqsForThisWorker,
				child.BaseDurationSecs,
				workerReqDelay,
			)
		}
	}

	wg.Wait()
	testDuration := time.Since(testStartTime)

	postTestMetrics := PostTestMetrics{testID, time.Duration(testDuration), config.GracePeriodPercent}

	streamPostTestLogs(
		logChan,
		&postTestMetrics,
		&globalCounts,
	)
}
