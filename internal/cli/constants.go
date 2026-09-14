// Package cli is the most important package for Punch's
// CLI-tool as it implements the core functionality it
// relies on. Disclaimer: Some of the types in
// ./cli/types.go can seem to hold duplication. Despite,
// package server and package cli sharing type and struct
// names, they're completely different.
package cli

var wantStatusCode uint16 = 404
var pWantStatusCode = &wantStatusCode
var ClientData = ClientTestData{
	TestID: GenerateTestID(),
	Config: PunchConfig{
		Protocol: "http",
		Target:   "localhost:3000/api/lilify/v1",
		Children: []Child{
			{
				Name:             "/urls?alias=123456",
				Method:           "GET",
				TotalRequests:    28,
				BaseDurationSecs: 4,
				WantStatusCode:   pWantStatusCode,
			},
			{
				Name:             "/urls?alias=bBKq8z",
				Method:           "GET",
				TotalRequests:    10,
				BaseDurationSecs: 2,
			},
		},
	},
}
