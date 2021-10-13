package config

import (
	"net"
	"net/http"

	"github.com/zkcrescent/chaos/utils"
)

const (
	DefaultHTTPTimeout = 10.0
)

type HTTPClient struct {
	Timeout             float64 `json:"timeout" yaml:"timeout" toml:"timeout"`
	DialTimeout         float64 `json:"dial_timeout" yaml:"dial_timeout" toml:"dial_timeout"`
	TLSHandshakeTimeout float64 `json:"tls_handshake_timeout" yaml:"tls_handshake_timeout" toml:"tls_handshake_timeout"`
	DisableKeepAlives   bool    `json:"disable_keep_alives" yaml:"disable_keep_alives" toml:"disable_keep_alives"`
}

func (c *HTTPClient) Validate() error {
	if c.Timeout == 0 {
		c.Timeout = DefaultHTTPTimeout
	}
	if c.DialTimeout == 0 {
		c.DialTimeout = DefaultHTTPTimeout
	}
	if c.TLSHandshakeTimeout == 0 {
		c.TLSHandshakeTimeout = DefaultHTTPTimeout
	}
	return nil
}

// Client will init one http client instance and SET IN GLOBAL
func (c *HTTPClient) Client() *http.Client {

	Global.HTTP = &http.Client{
		Timeout: utils.Duration(c.Timeout),
		Transport: &http.Transport{
			DisableKeepAlives: c.DisableKeepAlives,
			Dial: (&net.Dialer{
				Timeout: utils.Duration(c.DialTimeout),
			}).Dial,
			TLSHandshakeTimeout: utils.Duration(c.TLSHandshakeTimeout),
		},
	}
	return Global.HTTP
}
