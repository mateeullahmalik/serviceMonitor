package metadata

//ServiceAverageResponse is a Helper struct to fetch Average Response Time report from DB
type ServiceAverageResponse struct {
	Name             string  `json:"service"`
	AvrgResponseTime float64 `json:"average_response_time"`
}

//ServiceHits is a Helper struct to fetch Availability report from DB
type ServiceHits struct {
	Name             string `json:"service"`
	TotalHits        int    `json:"total_hits"`
	SuccessfulHits   int    `json:"successful_hits"`
	UnsuccessfulHits int    `json:"unsuccessful_hits"`
}
