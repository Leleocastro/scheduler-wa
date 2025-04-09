package domain

type ScheduleMessage struct {
	ID       string `json:"id"`
	Phone    string `json:"phone"`
	Text     string `json:"text"`
	SendAt   int64  `json:"sendAt"`
	Channel  string `json:"channel"`
	CronExpr string `json:"cronExpr,omitempty"`
	Repeats  int    `json:"repeats,omitempty"`
	Until    int64  `json:"until,omitempty"`
}
