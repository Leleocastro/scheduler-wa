package domain

type Plan struct {
	Name        string
	WebSocket   bool
	LimitPerDay int
	Route       string
	Group       string
}
