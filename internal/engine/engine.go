/*
Package controller controls workers and is responsible for
ensuring worker logs are properly streamed and outputted.
Punch's CLI-tool as it implements the core functionality it
relies on.
*/
package engine

import (
	"sync"

	"github.com/yuriongit/punch/internal/config"
	"github.com/yuriongit/punch/internal/streamer"
)

// RunTest starts a load test
func RunTest(dir string) error {
	cnf, err := config.LoadConfig(dir)
	if err != nil {
		return err
	}
	testID := config.CreateTestID()

	var wg sync.WaitGroup
	var logWg sync.WaitGroup

	// Buffer channel to prevent blocking worker goroutines
	logChanLen := uint32(50)
	for _, v := range cnf.Children {
		logChanLen += v.TotalRequests
	}
	logChan := make(chan string, logChanLen)

	logWg.Add(1)

	streamer.PreTestLogs(logChan, &testID)
	go streamer.TestData(logChan, &logWg)
	metrics, counts := RunWorkers(&wg, logChan, cnf, &testID)
	streamer.PostTestMetrics(logChan, &metrics, counts)

	close(logChan)
	logWg.Wait()

	return nil
}
