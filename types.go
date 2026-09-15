// Package cli is the most important package for Punch's
// CLI-tool as it implements the core functionality it
// relies on. Disclaimer: Some of the types in
// ./cli/types.go can seem to hold duplication. Despite,
// package server and package cli sharing type and struct
// names, they're completely different.
package cli

import (
	"encoding/json"
	"sync/atomic"
)

// Child defines the structure of the test, it's what allows
// the client to truly configure their load-test to their
// liking. Currently limited to HTTP with GET and POST methods.
// Eventually, Child will continue to expand to types such as
// HTTPChild and GraphQLChild.
type Child struct {
	Body   json.RawMessage `json:"body,omitempty"`
	Name   string          `json:"name" binding:"required"`
	Method string          `json:"method" binding:"required,oneof=GET POST"`

	TotalRequests    uint32 `json:"total_requests" binding:"required,gte=1,lte=4294967295"`
	BaseDurationSecs uint16 `json:"base_duration_secs" binding:"required,gte=1,lte=65535"`

	// optional
	WantStatusCode *uint16 `json:"want_status_code,omitempty" binding:"oneof=100 101 102 103 200 201 202 203 204 205 206 207 208 226 300 301 302 303 304 305 306 307 308 400 401 402 403 404 405 406 407 408 409 410 411 412 413 414 415 416 417 418 421 422 423 424 425 426 428 429 431 451 500 501 502 503 504 505 506 507 508 510 511"`
}

// PunchConfig defines the configuration for a load
// test.
type PunchConfig struct {
	Protocol           string  `json:"protocol" binding:"required,oneof=http https"`
	Target             string  `json:"target" binding:"required,url"`
	GracePeriodPercent uint8   `json:"grace_period_percent" binding:"required,gte=0,lte=50"`
	Children           []Child `json:"children" binding:"required,dive"`
}

// WorkerReqInfo defines what a worker needs to actually start
// sending requests to the child the worker is assigned to.
// Additionally, it allows the worker to determine if the
// request's response resulted in the desired output specified
// in the configuration.
type WorkerReqInfo struct {
	Method             string
	URL                string
	ChildName          string
	WantStatusCode     uint16
	GracePeriodPercent uint8
}

// ClientTestData defines the initial data required to
// initialize the load-test.
type ClientTestData struct {
	TestID string
	Config PunchConfig
}

// GlobalCounts carries the statuses of the delegated requests
// a load-test is set to fulfill. These counts are incremented
// by workers themselves.
type GlobalCounts struct {
	Curr   atomic.Uint32
	FatErr atomic.Uint32
	RegErr atomic.Uint32
	Suc    atomic.Uint32
}

// ChildCounts carries the statuses of its own delegated
// requests the worker is set out to fulfill. These counts are
// incremented by workers themselves.
type ChildCounts struct {
	Curr   uint32
	FatErr uint32
	RegErr uint32
	Suc    uint32
}

// CreateWorkerResLogCounts is what allows for helper CreateLog
// to accurately represent the requests that have been
// fulfilled by each worker.
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

// Log defines the structure of the test logs presented to
// the client.
type Log struct {
	Lvl            LogLvl
	WorkerID       uint32
	ChildID        uint32
	ReqMethod      string
	GotStatusCode  uint16
	WantStatusCode uint16
	Error          string
}
