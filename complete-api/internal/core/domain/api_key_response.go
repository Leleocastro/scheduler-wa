package domain

type ApiKeyResponse struct {
	Data []DataItem `json:"data"`
}

type DataItem struct {
	Key       string      `json:"key"`
	Tags      interface{} `json:"tags"`
	CreatedAt int64       `json:"created_at"`
	ID        string      `json:"id"`
	TTL       interface{} `json:"ttl"`
}
