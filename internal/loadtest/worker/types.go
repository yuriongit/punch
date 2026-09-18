/*
Package worker is responsible for carrying out the functionality
of sending HTTP requests to the configured 'target' and 'children'
in the Punch configuration file.
relies on. Workers are able to make HTTP requests both in parallel
and with concurrency.
*/
package worker

import (
	"sync/atomic"
	"time"
)

/*
	ReqInfo defines what a worker needs to actually start

sending requests to the child the worker is assigned to.
Additionally, it allows the worker to determine if the
request's response resulted in the desired output specified
in the configuration.
*/
type ReqInfo struct {
	Method             string
	URL                string
	ChildName          string
	WantStatusCode     uint16
	GracePeriodPercent uint8
}

/*
	GlobalCounts carries the statuses of the delegated requests

a load-test is set to fulfill. These counts are incremented
by workers themselves.
*/
type GlobalCounts struct {
	Curr    atomic.Uint32
	FatErr  atomic.Uint32
	RegErr  atomic.Uint32
	Suc     atomic.Uint32
	Workers atomic.Uint32
}

/*
	ChildCounts carries the statuses of its own delegated

requests the worker is set out to fulfill. These counts are
incremented by workers themselves.
*/
type ChildCounts struct {
	Curr   uint32
	FatErr uint32
	RegErr uint32
	Suc    uint32
}

/*
CreateWorkerResLogCounts is what allows for helper CreateLog
to accurately represent the requests that have been
fulfilled by each worker.
*/
type CreateWorkerResLogCounts struct {
	GlobalReqCounter uint32
	Curr             uint32
}

// LogLvl represents the log severity level for worker execution logs.
type LogLvl int

// Load returns the string output of the enum
func (l LogLvl) Load() string {
	return [...]string{"LOG-SUC", "LOG-FAT", "LOG-FIN", "LOG-ERR"}[l-1]
}

// EnumIdx returns the index of the LogLvl enum.
func (l LogLvl) EnumIdx() int {
	return int(l)
}

// LogLvl enums initialization
const (
	LogSuc LogLvl = iota + 1
	LogFat
	LogFin
	LogErr
)

/*
Log defines the structure of the test logs presented to
the client.
*/
type Log struct {
	Lvl            LogLvl
	WorkerID       uint32
	ChildID        uint32
	ReqMethod      string
	GotStatusCode  uint16
	WantStatusCode uint16
	Error          string
}

// PostTestMetrics ...
type PostTestMetrics struct {
	TestID             *string
	TestDur            time.Duration
	GracePeriodPercent uint8
}
