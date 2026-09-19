/*
Package controller controls workers and is responsible for
ensuring worker logs are properly streamed and outputted.
Punch's CLI-tool as it implements the core functionality it
relies on.
*/
package controller

import (
	"sync"

	"github.com/yuriongit/punch/internal/config"
	"github.com/yuriongit/punch/internal/loadtest/worker"
)

// RunTest starts a load test
func RunTest(directory string) error {
	clientConfig, err := config.ParseConfigFile(directory)
	if err != nil {
		return err
	}

	testID := config.GenerateTestID()

	var wg sync.WaitGroup
	var logWg sync.WaitGroup

	logWg.Add(1)

	// Buffer channel to prevent blocking worker goroutines
	logChanLen := uint32(50)
	for _, v := range clientConfig.Children {
		logChanLen += v.TotalRequests
	}
	logChan := make(chan string, logChanLen)

	go worker.StreamWorkerLogs(logChan, &logWg)
	worker.RunWorkers(&wg, logChan, clientConfig, &testID)

	close(logChan)
	logWg.Wait()
	return nil
}
