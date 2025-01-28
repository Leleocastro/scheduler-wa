package domain

type UsageResponse struct {
	Status string    `json:"status"`
	Data   UsageData `json:"data"`
}

type UsageData struct {
	ResultType string   `json:"resultType"`
	Result     []Result `json:"result"`
}

type Result struct {
	Metric Metric          `json:"metric"`
	Values [][]interface{} `json:"values"`
}

type Metric struct {
	Consumer string `json:"consumer"`
	Service  string `json:"service"`
}
