package domain

type Usage struct {
	Usage []UsageService `json:"usage"`
}

type UsageService struct {
	Service string      `json:"service"`
	Values  []UsageItem `json:"values"`
}

type UsageItem struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}
