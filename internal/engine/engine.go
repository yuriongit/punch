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
	"github.com/yuriongit/punch/internal/stream"
)

// RunTest starts a load test
func RunTest(dir string) error {
	cnf, err := config.GetConfigFile(dir)
	if err != nil {
		return err
	}

	testID := config.CreateTestID()
	
	var wg sync.WaitGroup
	var logWg sync.WaitGroup

	logWg.Add(1)

	// Buffer channel to prevent blocking worker goroutines
	logChanLen := uint32(50)
	for _, v := range cnf.Children {
		logChanLen += v.TotalRequests
	}
	logChan := make(chan string, logChanLen)

	stream.PreTestLogs(logChan, &testID)
	
	go stream.TestLogs(logChan, &logWg)
	RunWorkers(&wg, logChan, cnf, &testID)
	logWg.Wait()

	// stream.StreamPostTestLogs(
	// 	logChan,
	// 	&postTestMetrics,
	// 	&globalCounts,
	// )
	
	close(logChan)
	return nil
}
