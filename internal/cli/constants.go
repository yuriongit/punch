// Package cli implements the functionality of Punch's CLI tool. 
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
