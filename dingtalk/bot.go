package dingtalk

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/juju/errors"
)

// Bot is a dingtalk bot client
type Bot struct {
	Webhook string `json:"webhook" yaml:"webhook"`
}

// Response is the Webhook response
type Response struct {
	ErrorCode int    `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
}

func (r *Response) IsSuccess() bool {
	return r.ErrorCode == 0
}

// Send sends message to hook
func (d *Bot) Send(message Message) (*Response, error) {
	r, err := http.Post(d.Webhook, "application/json; charset=utf-8", strings.NewReader(message.Message()))
	if err != nil {
		return nil, errors.Trace(err)
	}
	if r == nil {
		return nil, nil
	}
	defer r.Body.Close()

	resp := new(Response)
	if err := json.NewDecoder(r.Body).Decode(resp); err != nil {
		return nil, errors.Trace(err)
	}

	return resp, err
}
