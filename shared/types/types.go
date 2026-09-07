package types

import "encoding/json"

type Child struct {
	Body   json.RawMessage `json:"body,omitempty"`
	Name   string          `json:"name" binding:"required"`
	Method string          `json:"method" binding:"required,oneof=GET POST"`

	// optional
	ConcurrencyRate *uint8 `json:"concurrency_rate,omitempty" binding:"gte=1,lte=5"`
	DurationSecs    *uint8 `json:"duration_secs" binding:"required,gte=10,lte=30"`
	// optional
	KnockoutAttempt *bool `json:"knockout_attempt,omitempty"`
	// optional
	ExpectedStatus *uint16 `json:"expected_status" binding:"omitempty,oneof=100 101 102 103 200 201 202 203 204 205 206 207 208 226 300 301 302 303 304 305 306 307 308 400 401 402 403 404 405 406 407 408 409 410 411 412 413 414 415 416 417 418 421 422 423 424 425 426 428 429 431 451 500 501 502 503 504 505 506 507 508 510 511"`
}

type PunchConfig struct {
	Protocol string  `json:"protocol" binding:"required,oneof=http https"`
	Target   string  `json:"target" binding:"required,url"`
	Children []Child `json:"children" binding:"required,dive"`
}
