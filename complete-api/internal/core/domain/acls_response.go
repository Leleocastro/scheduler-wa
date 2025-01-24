package domain

type ACLsResponse struct {
	Data []ACL `json:"data"`
}

type ACL struct {
	ID    string `json:"id"`
	Group string `json:"group"`
	Tags  string `json:"tags"`
}
