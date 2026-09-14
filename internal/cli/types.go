// Package cli implements the functionality of Punch's CLI tool. 
package cli

// In this file, types.go defines the structs and types needed
// for the CLI-tool Punch offers. Additionally, some of the
// types in ./cli/types.go might seem redundant, although,
// package server and package cli hold completely different
// types.

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

	TotalRequests    uint32 `json:"total_requests" binding:"required,gte=5,lte=50"`
	BaseDurationSecs uint16 `json:"duration_secs" binding:"required,gte=10,lte=30"`

	// optional
	WantStatusCode *uint16 `json:"want_status_code,omitempty" binding:"oneof=100 101 102 103 200 201 202 203 204 205 206 207 208 226 300 301 302 303 304 305 306 307 308 400 401 402 403 404 405 406 407 408 409 410 411 412 413 414 415 416 417 418 421 422 423 424 425 426 428 429 431 451 500 501 502 503 504 505 506 507 508 510 511"`
}

// PunchConfig defines the configuration for a load
// test.
type PunchConfig struct {
	Protocol string  `json:"protocol" binding:"required,oneof=http https"`
	Target   string  `json:"target" binding:"required,url"`
	Children []Child `json:"children" binding:"required,dive"`
}

// TestReqInfo defines what a worker needs to actually
// send a request.
type TestReqInfo struct {
	Method         string
	URL            string
	ChildName      string
	WantStatusCode uint16
}

type ClientTestData struct {
	TestID string
	Config PunchConfig
}

// GlobalCounts carries the statuses of the delegated requests 
// a load-test is set to fulfill. These counts are incremented
// by workers themselves.
type GlobalCounts struct {
  Curr  atomic.Uint32
	FatErr atomic.Uint32
	RegErr   atomic.Uint32
	Suc  atomic.Uint32
}

// ChildCounts carries the statuses of its own delegated
// requests the worker is set out to fulfill. These counts are
// incremented by workers themselves.
type ChildCounts struct {
	Curr  uint32
	FatErr uint32
	RegErr   uint32
	Suc  uint32
}

// CreateLogCounts is what allows for helper CreateLog
// to accurately represent the requests that have been
// fulfilled by each worker.
type CreateLogCounts struct {
  GlobalCurrReqAmt *atomic.Uint32
  WkrCurrReqAmt uint32
  Lvls ChildCounts
}


type LogLvl int

const (
  LogSuc LogLvl = iota
  LogFat
  LogFin
  LogErr
)

func (l LogLvl) Load() string {
	return [...]string{"LOG-SUC", "LOG-FAT", "LOG-FIN", "LOG-ERR"}[l-1]
}
func (l LogLvl) EnumIdx() int {
  return int(l)
}

// Log messages
type Log struct {
  Lvl LogLvl
  WkrID uint32
  ReqMethod string
  ChildName string
  GotStatusCode uint16
  WantStatusCode uint16
}
