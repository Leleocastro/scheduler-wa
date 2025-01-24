package domain

type PluginsResponse struct {
	Data []Plugin `json:"data"`
}

type Plugin struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Route Route  `json:"route"`
}

type Route struct {
	ID   string  `json:"id"`
	Name *string `json:"name"`
}
