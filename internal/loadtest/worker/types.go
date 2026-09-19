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

	"github.com/yuriongit/punch/internal/config"
	"github.com/yuriongit/punch/internal/shared"
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
	WantStatusCode     shared.StatusCode
	GracePeriodPercent uint8
}

/*
	GlobalCounts carries the statuses of the delegated requests

a load-test is set to fulfill. These counts are incremented
by workers themselves.
*/
type GlobalCounts struct {
	Current    atomic.Uint32
	FatalErr   atomic.Uint32
	RegularErr atomic.Uint32
	Success    atomic.Uint32
	Workers    atomic.Uint32
}

/*
	ChildCounts carries the statuses of its own delegated

requests the worker is set out to fulfill. These counts are
incremented by workers themselves.
*/
type ChildCounts struct {
	Current    uint32
	FatalErr   uint32
	RegularErr uint32
	Success    uint32
}

/*
CurrentRequestCounts is what allows for helper CreateLog
to accurately represent the requests that have been
fulfilled by each worker.
*/
type CurrentRequestCounts struct {
	GlobalCurrent uint32
	WorkerCurrent uint32
}

// LogLvl represents the log severity level for worker execution logs.
type LogLvl int

// Load returns the string output of the enum
func (l LogLvl) Load() string {
	return [...]string{"LOG-SUCC", "LOG-FATA", "LOG-FINI", "LOG-ERRO"}[l-1]
}

// EnumIdx returns the index of the LogLvl enum.
func (l LogLvl) EnumIdx() int {
	return int(l)
}

// LogLvl enums initialization
const (
	LogSucc LogLvl = iota + 1
	LogFata
	LogFini
	LogErro
)

type ID uint32

type IDs struct {
	Worker ID
	Child  ID
}

type StatusCodes struct {
	Got  shared.StatusCode
	Want shared.StatusCode
}

/*
Log defines the structure of the test logs presented to
the client.
*/
type Log struct {
	Lvl         LogLvl
	IDs         IDs
	ReqMethod   string
	StatusCodes StatusCodes
	Latency     time.Duration
}

// PostTestMetrics ...
type PostTestMetrics struct {
	TestID             *config.TestID
	TestDuration       time.Duration
	GracePeriodPercent uint8
}
