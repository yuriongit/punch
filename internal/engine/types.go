/*
Package worker is responsible for carrying out the functionality
of sending HTTP requests to the configured 'target' and 'children'
in the Punch configuration file.
relies on. Workers are able to make HTTP requests both in parallel
and with concurrency.
*/
package engine

import (
	"time"

	"github.com/yuriongit/punch/internal/domain"
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
	WantStatusCode     domain.StatusCode
	GracePeriodPercent uint8
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

type ID uint32

type IDs struct {
	Worker ID
	Child  ID
}

type StatusCodes struct {
	Got  domain.StatusCode
	Want domain.StatusCode
}

/*
Log defines the structure of the test logs presented to
the client.
*/
type Log struct {
	Lvl         domain.LogLevel
	IDs         IDs
	ReqMethod   string
	StatusCodes StatusCodes
	Latency     time.Duration
}

