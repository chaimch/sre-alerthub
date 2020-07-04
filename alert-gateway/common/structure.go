package common

import "time"

type Alerts []Alert

type Alert struct {
	ActiveAt    time.Time `json:"active_at"`
	Annotations struct {
		Description string `json:"description"`
		Summary     string `json:"summary"`
		RuleId      string `json:"rule_id"`
	} `json:"annotations"`
	FiredAt    time.Time         `json:"fired_at"`
	Labels     map[string]string `json:"labels"`
	LastSentAt time.Time         `json:"last_sent_at"`
	ResolvedAt time.Time         `json:"resolved_at"`
	State      int               `json:"state"`
	ValidUntil time.Time         `json:"valid_until"`
	Value      float64           `json:"value"`
}

type UserGroup struct {
	Id                    int64
	StartTime             string
	EndTime               string
	Start                 int
	Period                int
	ReversePolishNotation string
	User                  string
	Group                 string
	DutyGroup             string
	Method                string
}
