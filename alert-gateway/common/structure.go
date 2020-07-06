package common

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"
)

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

type SingleAlert struct {
	Id       int64             `json:"id"`
	Count    int               `json:"count"`
	Value    float64           `json:"value"`
	Summary  string            `json:"summary"`
	Hostname string            `json:"hostname"`
	Labels   map[string]string `json:"labels"`
}

type Ready2Send struct {
	RuleId int64
	Start  int64
	User   []string
	Alerts []SingleAlert
}

type ValidUserGroup struct {
	User      string
	Group     string
	DutyGroup string
}

var Lock sync.Mutex
var Rw sync.RWMutex

var Maintain map[string]bool

var RuleCount map[[2]int64]int64

var ErrHttpRequest = errors.New("create HTTP request failed")

/*
 Check if UserGroup is valid.
*/
func (u UserGroup) IsValid() bool {
	return u.User != "" || u.DutyGroup != "" || u.Group != ""
}

/*
 IsOnDuty return if current UserGroup is on duty or not by StartTime & EndTime.
 If the UserGroup is not on duty, alerts should not be sent to them.
*/
func (u UserGroup) IsOnDuty() bool {
	now := time.Now().Format("15:04")

	return (u.StartTime <= u.EndTime && u.StartTime <= now && u.EndTime >= now) || // 不跨 00:00
		(u.StartTime > u.EndTime && (u.StartTime <= now || now <= u.EndTime)) // // 跨 00:00
}

func HttpPost(url string, params map[string]string, headers map[string]string, body []byte) (*http.Response, error) {
	//new request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, ErrHttpRequest
	}
	//add params
	q := req.URL.Query()
	if params != nil {
		for key, val := range params {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
	}
	//add headers
	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}
	//http client
	client := &http.Client{Timeout: 5 * time.Second} //Add the timeout,the reason is that the default client has no timeout set; if the remote server is unresponsive, you're going to have a bad day.
	return client.Do(req)
}

func HttpGet(url string, params map[string]string, headers map[string]string) (*http.Response, error) {
	//new request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
		return nil, ErrHttpRequest
	}
	//add params
	q := req.URL.Query()
	if params != nil {
		for key, val := range params {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
	}
	//add headers
	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}
	//http client
	client := &http.Client{Timeout: 5 * time.Second} //Add the timeout,the reason is that the default client has no timeout set; if the remote server is unresponsive, you're going to have a bad day.
	return client.Do(req)
}
