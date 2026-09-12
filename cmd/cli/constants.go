package main

var expectedStatus uint16 = 404
var pExpectedStatus = &expectedStatus
var clientTestData = ClientTestData{
	TestID: "12CHARACTERS",
	Config: PunchConfig{
		Protocol: "http",
		Target:   "localhost:3000/api/lilify/v1",
		Children: []Child{
			{
				Name:           "/urls?alias=gAFYyK",
				Method:         "GET",
				TotalRequests:  5,
				DurationSecs:   2,
				ExpectedStatus: nil,
			},
			{
				Name:           "/urls?alias=gAFYyK",
				Method:         "GET",
				TotalRequests:  10,
				DurationSecs:   2,
				ExpectedStatus: pExpectedStatus,
			},
		},
	},
}
