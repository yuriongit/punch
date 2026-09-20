package domain

import "sync/atomic"

/*
* GlobalCounts carries the statuses of the delegated requests
* a load-test is set to fulfill. These counts are incremented
* by workers themselves.
 */
type GlobalCounts struct {
	Current    atomic.Uint32
	FatalErr   atomic.Uint32
	RegularErr atomic.Uint32
	Success    atomic.Uint32
	Workers    atomic.Uint32
}
