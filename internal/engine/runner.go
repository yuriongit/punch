/*
Package worker is responsible for carrying out the functionality
of sending HTTP requests to the configured 'target' and 'children'
in the Punch configuration file.
relies on. Workers are able to make HTTP requests both in parallel
and with concurrency.
*/
package engine

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/yuriongit/punch/internal/domain"
	"github.com/yuriongit/punch/internal/ui"
)

type Logger interface {
	Log(entry domain.LogEntry) // create LogEntry struct
}

type Request struct {
	Method             string
	URL                string
	WantStatusCode     domain.StatusCode
	GracePeriodPercent uint8
}

type IDs struct {
	Worker domain.ID
	Child  domain.ID
}

// Worker holds all state needed for a single goroutine worker to execute requests.
type Worker struct {
	ids        IDs
	childId    domain.ID
	req        Request
	targetReqs uint32
	delay      time.Duration
	duration   time.Duration

	client *http.Client
	counts *domain.GlobalCounts
	logger Logger
}

func NewWorker(
	ids IDs,
	req Request,
	targetReqs uint32,

	delay time.Duration,
	duration time.Duration,

	client *http.Client,
	counts *domain.GlobalCounts,
	logger Logger,
) *Worker {
	return &Worker{
		ids:        ids,
		req:        req,
		targetReqs: targetReqs,
		delay:      delay,
		duration:   duration,

		client: client,
		counts: counts,
		logger: logger,
	}
}

func (w *Worker) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	if w.targetReqs == 0 {
		return
	}

	if w.req.Method != http.MethodGet {
		panic("not yet implemented")
	}

	gracePeriod := (w.duration * time.Duration(w.req.GracePeriodPercent)) / 100
	deadline := time.Now().Add(w.duration + gracePeriod)

	var localReqCount uint32

	for time.Now().Before(deadline) {
		time.Sleep(w.delay)

		start := time.Now()
		resp, err := w.client.Get(w.req.URL)
		latency := time.Since(start)

		localReqCount++
		w.counts.Current.Add(1)

		w.handleResponse(resp, err, latency, localReqCount)

		if localReqCount == w.targetReqs {
			return
		}
	}
}

func (w *Worker) handleResponse(
	resp *http.Response,
	err error,
	latency time.Duration,
	localCount uint32,
) {
	var gotStatus domain.StatusCode
	var lvl domain.LogLevel

	switch {
	case err != nil:
		w.counts.FatalErr.Add(1)
		lvl = domain.LogFata
	case resp != nil && resp.StatusCode == int(w.req.WantStatusCode):
		w.counts.Success.Add(1)
		lvl = domain.LogSucc
		_ = resp.Body.Close()
		gotStatus = domain.StatusCode(resp.StatusCode)
	default:
		w.counts.RegularErr.Add(1)
		lvl = domain.LogErro // relative
		if resp != nil {
			gotStatus = domain.StatusCode(resp.StatusCode)
			_ = resp.Body.Close()
		}
	}

	if w.logger != nil {
		w.logger.Log(domain.LogEntry{
			Level:      lvl,
			WorkerID:   w.ids.Worker,
			ChildID:    w.ids.Child,
			WantStatus: w.req.WantStatusCode,
			GotStatus:  gotStatus,
			Latency:    latency,
			GlobalReq:  domain.RequestCount(w.counts.Current.Load()),
			WorkerReq:  domain.RequestCount(localCount),
		})
	}
}

// Runner orchestrates all workers.
type Runner struct {
	cfg    *domain.ConfigFile
	client *http.Client
	logger Logger
}

func NewRunner(cfg *domain.ConfigFile, logger Logger, client *http.Client) *Runner {
	return &Runner{
		cfg:    cfg,
		logger: logger,
		client: client,
	}
}

// Run executes the full load test pool.
func (r *Runner) Run(wg *sync.WaitGroup, testID *domain.TestID) (domain.PostTestMetrics, *domain.GlobalCounts) {
	var counts domain.GlobalCounts
	testStartTime := time.Now()

	for childID, child := range r.cfg.Children {
		if child.BaseDurationSecs == 0 || child.TotalRequests == 0 {
			continue
		}

		numWorkers := uint32(float32(child.TotalRequests) / float32(child.BaseDurationSecs))
		if numWorkers < 1 {
			numWorkers = 1
		}
		if numWorkers > child.TotalRequests {
			numWorkers = child.TotalRequests
		}

		baseReqs := child.TotalRequests / numWorkers
		remainder := child.TotalRequests % numWorkers

		wantStatus := domain.StatusCode(200)
		if child.WantStatusCode != nil {
			wantStatus = *child.WantStatusCode
		}

		url := fmt.Sprintf("%s://%s%s", r.cfg.Protocol, r.cfg.Target, child.Name)

		req := Request{
			Method:             child.Method,
			URL:                url,
			WantStatusCode:     wantStatus,
			GracePeriodPercent: r.cfg.GracePeriodPercent,
		}

		baseDuration := time.Duration(child.BaseDurationSecs) * time.Second

		for i := uint32(0); i < numWorkers; i++ {
			workerReqs := baseReqs
			if i < remainder {
				workerReqs++
			}

			var delay time.Duration
			if workerReqs > 0 {
				delay = time.Duration((float64(child.BaseDurationSecs) / float64(workerReqs)) * float64(time.Second))
			}

			counts.Workers.Add(1)
			wg.Add(1)

			worker := NewWorker(
				IDs{
					Worker: domain.ID(counts.Workers.Load()),
					Child:  domain.ID(childID + 1),
				},
				req,
				workerReqs,
				delay,
				baseDuration,
				r.client,
				&counts,
				r.logger,
			)

			go worker.Run(wg)
		}
	}

	wg.Wait()

	return domain.PostTestMetrics{
		TestID:             testID,
		TestDuration:       time.Since(testStartTime),
		GracePeriodPercent: r.cfg.GracePeriodPercent,
	}, &counts
}

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
createWorkerLog creates a formatted worker log
with custom colored brackets and light grey sub-logs.
*/
func createWorkerLog(
	logChan chan<- string,
	l *Log,
	c *CurrentRequestCounts,
) {
	formattedLogLevel := ui.FormatLogLevel(l.Lvl, true)
	formattedTimestamp := ui.FormatTimestamp(time.Now(), l.Lvl, true) // Bold(false) on timestamp
	formattedLatency := ui.FormatLatency(l.Latency)

	// TODO: implement and use foreseen "color" parameter (bool):
	formatLogLvlAccent := ui.LogLevelAccent(formattedLogLevel)

	// 2. Colorize status level text
	coloredLvl := ui.LogLevelAccent(formattedLogLevel)

	// 3. Highlight Got and Want status codes using the log level's style
	formattedGotWantStatusCodes := fmt.Sprintf("%d/%d", l.StatusCodes.Got, l.StatusCodes.Want)

	// 4. Accent-color the opening and closing brackets
	openBracket := formatLogLvlAccent.Render("[ ")
	closeBracket := formatLogLvlAccent.Render(" ]")

	header := fmt.Sprintf(
		"%s %s%s %s Child #%d, Worker #%d%s\n",
		formattedTimestamp,
		openBracket,
		coloredLvl,
		ui.FaintStyle.Render("|"),
		l.IDs.Child,
		l.IDs.Worker,
		closeBracket,
	)

	// 5. Build raw tree structure
	rawSubLogs := fmt.Sprintf(
		"  %s\n"+
			"  ├── Global Request: #%d\n"+
			"  ├── My request: #%d\n"+
			"  └── Latency: %s",
		ui.LogLevel(l.Lvl.Load()).Render(
			fmt.Sprintf("├── Got/Want: %s", formattedGotWantStatusCodes),
		),
		c.GlobalCurrent,
		c.WorkerCurrent,
		formattedLatency,
	)

	// 6. Render the sub-logs in light grey (248)
	coloredSubLogs := ui.SubLogStyle.Render(rawSubLogs)

	logChan <- (header + coloredSubLogs)
}

func executeWorker(
	logChan chan<- string,
	wg *sync.WaitGroup,
	globalCounts *domain.GlobalCounts,
	IDs IDs,
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
		workerCounts := ChildCounts{}

		for time.Now().Before(deadline) {
			// Pause to stretch requests over full BaseDuration
			time.Sleep(workerReqDelay)

			// Measure pure HTTP round-trip latency
			requestStartTime := time.Now()
			resp, err := httpClient.Get(req.URL)
			latency := time.Since(requestStartTime)

			workerCounts.Current++
			globalCounts.Current.Add(1)

			switch {
			case err != nil:
				workerCounts.FatalErr++
				globalCounts.FatalErr.Add(1)

				l := Log{
					Lvl:       domain.LogFata,
					IDs:       IDs,
					ReqMethod: req.Method,
					StatusCodes: StatusCodes{
						Want: req.WantStatusCode,
						Got:  0,
					},
					Latency: latency,
				}
				c := CurrentRequestCounts{
					GlobalCurrent: globalCounts.Current.Load(),
					WorkerCurrent: workerCounts.Current,
				}

				createWorkerLog(logChan, &l, &c)

			case resp.StatusCode == int(req.WantStatusCode):
				workerCounts.Success++
				globalCounts.Success.Add(1)

				l := Log{
					Lvl:       domain.LogSucc,
					IDs:       IDs,
					ReqMethod: req.Method,
					StatusCodes: StatusCodes{
						Want: domain.StatusCode(req.WantStatusCode),
						Got:  domain.StatusCode(resp.StatusCode), //nolint:gosec // Status codes safely fit
					},
					Latency: latency,
				}
				c := CurrentRequestCounts{
					GlobalCurrent: globalCounts.Current.Load(),
					WorkerCurrent: workerCounts.Current,
				}

				createWorkerLog(logChan, &l, &c)
				_ = resp.Body.Close()

			default:
				workerCounts.RegularErr++
				globalCounts.RegularErr.Add(1)

				l := Log{
					Lvl:       domain.LogErro,
					IDs:       IDs,
					ReqMethod: req.Method,
					StatusCodes: StatusCodes{
						Want: domain.StatusCode(req.WantStatusCode),
						Got:  domain.StatusCode(resp.StatusCode), //nolint:gosec // Status codes safely fit
					},
					Latency: latency,
				}
				c := CurrentRequestCounts{
					GlobalCurrent: globalCounts.Current.Load(),
					WorkerCurrent: workerCounts.Current,
				}

				createWorkerLog(logChan, &l, &c)
				_ = resp.Body.Close()
			}

			// Stop when worker fulfills assigned quota
			if workerCounts.Current == targetRequests {
				return
			}
		}
	}
}

// RunWorkers ...
func RunWorkers(
	wg *sync.WaitGroup,
	logChan chan<- string,
	config *domain.ConfigFile,
	testID *domain.TestID,
) (metrics domain.PostTestMetrics, counts *domain.GlobalCounts) {
	globalCounts := domain.GlobalCounts{
		Current:    atomic.Uint32{},
		Workers:    atomic.Uint32{},
		FatalErr:   atomic.Uint32{},
		RegularErr: atomic.Uint32{},
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
			childExpectedStatus := domain.StatusCode(200)
			child.WantStatusCode = &childExpectedStatus
		}

		// Start constructing NewWorker data
		workersIDs := IDs{
			Worker: domain.ID((globalCounts.Workers.Load())),
			Child:  domain.ID(childID + 1),
		}
		req := Request{
			Method:             config.Children[workersIDs.Child].Method,
			URL:                fmt.Sprintf("%s://%s%s", config.Protocol, config.Target, child.Name),
			WantStatusCode:     *config.Children[workersIDs.Child].WantStatusCode,
			GracePeriodPercent: config.GracePeriodPercent,
		}

		for i := uint32(0); i < numWorkers; i++ {
			targetReqs := baseReqsPerWorker
			if i < remainderReqs {
				targetReqs++
			}

			// Delay calculation per request to fill BaseDurationSecs properly
			var reqDelay time.Duration
			if targetReqs > 0 {
				reqDelay = time.Duration((float64(child.BaseDurationSecs) / float64(targetReqs)) * float64(time.Second))
			}

			worker := NewWorker(
				workersIDs,
				req,
				targetReqs,
				reqDelay,
				time.Duration(child.BaseDurationSecs),
				httpClient,
				counts,
				logger,
			)

			globalCounts.Workers.Add(1)
			wg.Add(1)

			workerID := domain.ID(globalCounts.Workers.Add(1))

			go executeWorker(
				logChan,
				wg,
				&globalCounts,
				IDs{
					Worker: workerID,
					Child:  domain.ID(childID + 1),
				},
				req,
				targetReqs,
				child.BaseDurationSecs,
				reqDelay,
			)
		}
	}

	wg.Wait()
	testDuration := time.Since(testStartTime)

	// TODO: Move streaming logs to controller.

	return domain.PostTestMetrics{
		TestID:             testID,
		TestDuration:       time.Duration(testDuration),
		GracePeriodPercent: config.GracePeriodPercent,
	}, &globalCounts
}
