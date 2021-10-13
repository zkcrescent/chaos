package config

import (
	"net/http"

	gorp "gopkg.in/gorp.v2"
)

// package config defines common config for most usage at backend

// configStore for some common usage in global config
type configStore struct {
	HTTP *http.Client
	DB   *gorp.DbMap
	// Config for set custom config
	Config interface{}
}

var Global = &configStore{}

func (c *configStore) SetConfig(in interface{}) {
	c.Config = in
}
