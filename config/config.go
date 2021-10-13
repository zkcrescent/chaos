package config

import "net/http"

// package config defines common config for most usage at backend

// configStore for some common usage in global config
type configStore struct {
	HTTP *http.Client
	// Config for set custom config
	Config interface{}
}

var Global = &configStore{}

func (c *configStore) SetConfig(in interface{}) {
	c.Config = in
}
