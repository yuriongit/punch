// Package main is the entry point of Punch's CLI tool.
package main

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/yuriongit/punch/internal/cli"
)

func main() {
	// Create temporary pointers for config and testID constants
	config := &cli.ClientData.Config
	testID := &cli.ClientData.TestID

	// Create 2 WaitGroups, one for main and one for StreamLogs
	var wg sync.WaitGroup
	var logWg sync.WaitGroup

	logWg.Add(1)

	// Channel for logs
	logChan := make(chan string, 125)

	// Spawn StreamLogs goroutine that runs concurrently
	go StreamLogs(logChan, &logWg)
	RunTestWorkers(&wg, logChan, config, testID)

	// When RunTestWorkers returns, close logChan
	close(logChan)
	logWg.Wait()
}

// StreamLogs ...
func StreamLogs(
	logChan chan string,
	logWg *sync.WaitGroup,
) {
	defer logWg.Done()
	for msg := range logChan {
		fmt.Println(msg)
	}
}

// RunTestWorkers ...
func RunTestWorkers(
	wg *sync.WaitGroup,
	logChan chan<- string,
	config *cli.PunchConfig,
	testID *string,
) {
  // Output initializing logs with line break
  logChan<-"Punch------------------------------------------------------------------"
	logChan <- fmt.Sprintf("[INIT] Initializing load test for TEST-%s", *testID)
	logChan<-"------------------------------------------------------------------Punch"
	
	var wkrCount atomic.Uint32
	var globalReqCount atomic.Uint32
	testStartTime := time.Now()
	
	// Concurrently run requests against each child
	for _, v := range config.Children {
		rps := float32(v.TotalRequests) / float32(v.BaseDurationSecs)
		numWkrs := uint32(rps)
		if numWkrs < 1 {
			numWkrs = 1
		}

		reqsPerWkr := v.TotalRequests / uint32(numWkrs)
		
		// 3. Calculate delay per worker to stretch requests over full DurationSecs
		// E.g, child #1: 10s - 10 reqs/worker = request/1s
		wkrReqDelayMs := (float64(v.BaseDurationSecs) / float64(reqsPerWkr)) * 1000
		
		var globalWkrCount atomic.Uint32

		// Handle optional field
		if v.WantStatusCode == nil {
			childExpectedStatus := uint16(200)
			v.WantStatusCode = &childExpectedStatus
		}

		URL := fmt.Sprintf("%s://%s%s", config.Protocol, config.Target, v.Name)

		for j := uint32(0); j < numWkrs; j++ {
		  wkrCount.Add(1)
    
			wg.Add(1)

			go executeWorker(
				logChan,
				wg,
				&globalReqCount,
				globalWkrCount.Add(1),
				cli.TestReqInfo{
					Method:         v.Method,
					URL:            URL,
					ChildName:      v.Name,
					WantStatusCode: *v.WantStatusCode,
				},
				reqsPerWkr,
				v.BaseDurationSecs,
				wkrReqDelayMs,
			)
		}
	}

	wg.Wait()
	
	// Output test completion with line break
	logChan<-"Punch------------------------------------------------------------------"
	logChan<-"[SUCCESS] Completed load test successfully!"
	logChan<-fmt.Sprintf("[LOG-MET] Total test duration: %.4fs", time.Since(testStartTime).Seconds())
	logChan<-fmt.Sprintf("[LOG-MET] Total workers: %d", wkrCount.Load())
	logChan<-fmt.Sprintf("[LOG-MET] Fulfilled requests: %d", globalReqCount.Load())
	logChan<-fmt.Sprintf("[LOG-MET] Fulfilled requests: %d", globalReqCount.Load())
	logChan<-"------------------------------------------------------------------Punch\n"
}

// CreateLog outputs
func CreateLog(logChan chan<- string, l cli.Log, c *cli.CreateLogCounts) {
  switch l.Lvl.Load() {
    case "LOG-SUC":
    logChan <- fmt.Sprintf(
      "%s [%s] %d %s %s | Global-Req #%d / Wkr-Req #%d | Got %d, want %d", 
      time.Now().Format(time.RFC3339), 
      l.Lvl.Load(), 
      l.WkrID,
      l.ReqMethod,
      l.ChildName,
      c.GlobalCurrReqAmt.Load(),
      c.WkrCurrReqAmt,
      l.GotStatusCode,
      l.WantStatusCode,
    )
  }
  // switch l.Lvl.Load() {
  //   case "LOG-FIN":
  //   logChan<-fmt.Sprintf(
  //     "[%s] WKR-%d - %s %s | Completed %d requests - Fatal %d / Error: %d / Success: %d | My duration: %.4f seconds", time.Now().Format(time.RFC3339Nano), l.WkrID,
  //     )
  // }
}

// executeWorker ...
func executeWorker(
	logChan chan<- string,
	wg *sync.WaitGroup,
	globalReqCount *atomic.Uint32,
	wkrID uint32,
	req cli.TestReqInfo,
	reqsToDo uint32,
	baseDur uint16,
	wkrReqDelayMs float64,
) {
	defer wg.Done()

	switch req.Method {
	case "POST":
		panic("not yet implemented")
	case "GET":
		bufferDuration := float32(baseDur) * 0.25
		deadline := time.Now().Add((time.Duration(baseDur) * time.Second) + time.Duration(bufferDuration))

		reqCounts := cli.ChildCounts{}

		resStart := time.Now()
		for time.Now().Before(deadline) {
  		var resEnd time.Duration
			resp, err := http.Get(req.URL)
			
			
			time.Sleep(time.Duration(wkrReqDelayMs) * time.Millisecond)

			reqCounts.Curr++
			globalReqCount.Add(1)
			
			switch {
			case reqCounts.Curr == reqsToDo:
				time.Sleep((time.Duration(bufferDuration) * time.Second ) / 2)
				
				// logWithTime(logChan, "LOG-FIN", fmt.Sprintf("WKR-%d - %s %s | Completed %d requests - Fatal %d / Error: %d / Success: %d | My duration: %.4f seconds", 
				// wkrID,
				// req.Method,
				// req.ChildName,
				// localCounts.Curr,
				// localCounts.FatErr,
				// localCounts.RegErr,
				// localCounts.Success,
				// testEnd.Seconds(),
				// ))
				
				return
			case err != nil:
				reqCounts.FatErr++
				time.Since(resStart)
				
				h, m, s := time.Now().Clock()
				
				logChan <- fmt.Sprintf(
					"%s [LOG-FAT] WRK-%d | Global-Req #%d / Wkr-Req #%d - %s %s | Res-Time: %d | Response: '%s'",
					fmt.Sprintf("%d:%d:%d", h,m,s),
					wkrID,
					globalReqCount.Load(),
					reqCounts.Curr,
					req.Method,
					req.ChildName,
					err.Error(),
					resEnd,
				)
			case resp.StatusCode == int(req.WantStatusCode):
				reqCounts.Suc++
				
				h, m, s := time.Now().Clock()

				// successful request
				logChan <- fmt.Sprintf(
					"%s [LOG-SUC] WRK-%d | Global-Req #%d / Wkr-Req #%d - %s %s | got %d, want %d",
					fmt.Sprintf("%d:%d:%d", h,m,s),
					wkrID,
					globalReqCount.Load(),
					reqCounts.Curr,
					req.Method,
					req.ChildName,
					resp.StatusCode,
					req.WantStatusCode,
				)

				_ = resp.Body.Close()
			default:
				reqCounts.RegErr++

				logChan <- fmt.Sprintf(
					"[LOG-ERR] WRK-%d | Global-Req #%d / Wkr-Req #%d - %s %s | got %d, want %d",
					wkrID,
					globalReqCount.Load(),
					reqCounts.Curr,
					req.Method,
					req.ChildName,
					resp.StatusCode,
					req.WantStatusCode,
				)

				_ = resp.Body.Close()
			}
		}
	}
}
