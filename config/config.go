package config

import (
	"fmt"
	"net/http"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/go-redis/redis/v8"
	gorp "gopkg.in/gorp.v2"
)

// package config defines common config for most usage at backend

// configStore for some common usage in global config
type configStore struct {
	HTTP      *http.Client
	DB        *gorp.DbMap
	FileStore FileStore
	Redis     *redis.Client
	// Config for set custom config
	Config interface{}
}

var Global = &configStore{}

func (c *configStore) SetConfig(in interface{}) {
	c.Config = in
}

// GetOssBuckt returns FileStore to oss file store or panic
func (c *configStore) GetOssBuckt() *oss.Bucket {
	store, ok := c.FileStore.(*OSSStore)
	if !ok {
		panic(fmt.Sprintf(
			"cont convert file store to oss store: current: %v",
			c.FileStore,
		))
	}
	return store.bucket
}
