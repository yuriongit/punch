/*
Package worker is responsible for carrying out the functionality
of sending HTTP requests to the configured 'target' and 'children'
in the Punch configuration file.
relies on. Workers are able to make HTTP requests both in parallel
and with concurrency.
*/
package worker

// Miscellaneous
const (
	// Formats
	timeFormat = "15:04:05"
	dateFormat = "2006-01-02"
)
