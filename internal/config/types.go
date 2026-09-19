/*
Package config offers the functionality of Punch configuration files
and creating unique test IDs. Alongside with the types and structs
that the rest of Punch uses throughout the entire application.
*/
package config

import (
	"encoding/json"

	"github.com/yuriongit/punch/internal/shared"
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
	WantStatusCode *shared.StatusCode `json:"want_status_code,omitempty" binding:"oneof=100 101 102 103 200 201 202 203 204 205 206 207 208 226 300 301 302 303 304 305 306 307 308 400 401 402 403 404 405 406 407 408 409 410 411 412 413 414 415 416 417 418 421 422 423 424 425 426 428 429 431 451 500 501 502 503 504 505 506 507 508 510 511"`
}

// File defines the configuration for a load
// test.
type File struct {
	Protocol           string  `json:"protocol" binding:"required,oneof=http https"`
	Target             string  `json:"target" binding:"required,url"`
	GracePeriodPercent uint8   `json:"grace_period_percent" binding:"required,gte=0,lte=50"`
	Children           []Child `json:"children" binding:"required,dive"`
}

type TestID string

/*
ClientTestData defines the initial data required to
initialize the load-test.
*/
type ClientTestData struct {
	TestID TestID
	Config File
}
