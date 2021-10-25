package config

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	Network            string        `json:"network" yaml:"network" valid:"required" toml:"network"`
	Addr               string        `json:"addr" yaml:"addr" valid:"required" toml:"addr"`
	Password           string        `json:"password" yaml:"password" valid:"required" toml:"password"`
	DB                 int           `json:"db" yaml:"db" valid:"required" toml:"db"`
	MaxRetries         int           `json:"maxretries" yaml:"maxretries" toml:"max_retries"`
	MinRetryBackoff    time.Duration `json:"minretrybackoff" yaml:"minretrybackoff" toml:"min_retry_backoff"`
	MaxRetryBackoff    time.Duration `json:"maxretrybackoff" yaml:"maxretrybackoff" toml:"max_retry_backoff"`
	DialTimeout        time.Duration `json:"dialtimeout" yaml:"dialtimeout" toml:"dial_timeout"`
	ReadTimeout        time.Duration `json:"readtimeout" yaml:"readtimeout" toml:"read_timeout"`
	WriteTimeout       time.Duration `json:"writetimeout" yaml:"writetimeout" toml:"write_timeout"`
	PoolSize           int           `json:"poolsize" yaml:"poolsize" toml:"pool_size"`
	MinIdleConns       int           `json:"minidleconns" yaml:"minidleconns" toml:"min_idle_conns"`
	MaxConnAge         time.Duration `json:"maxconnage" yaml:"maxconnage" toml:"max_conn_age"`
	PoolTimeout        time.Duration `json:"pooltimeout" yaml:"pooltimeout" toml:"pool_timeout"`
	IdleTimeout        time.Duration `json:"idletimeout" yaml:"idletimeout" toml:"idle_timeout"`
	IdleCheckFrequency time.Duration `json:"idlecheckfrequency" yaml:"idlecheckfrequency" toml:"idle_check_frequency"`
}

func (c Redis) Client() (*redis.Client, error) {
	cli := redis.NewClient(&redis.Options{
		Network:            c.Network,
		Addr:               c.Addr,
		Password:           c.Password,
		DB:                 c.DB,
		MaxRetries:         c.MaxRetries,
		MinRetryBackoff:    c.MinRetryBackoff,
		MaxRetryBackoff:    c.MaxRetryBackoff,
		DialTimeout:        c.DialTimeout,
		ReadTimeout:        c.ReadTimeout,
		WriteTimeout:       c.WriteTimeout,
		PoolSize:           c.PoolSize,
		MinIdleConns:       c.MinIdleConns,
		MaxConnAge:         c.MaxConnAge,
		PoolTimeout:        c.PoolTimeout,
		IdleTimeout:        c.IdleTimeout,
		IdleCheckFrequency: c.IdleCheckFrequency,
	})
	_, err := cli.Ping(context.Background()).Result()
	return cli, err
}

func (c Redis) GClient() (*redis.Client, error) {
	cli := redis.NewClient(&redis.Options{
		Network:            c.Network,
		Addr:               c.Addr,
		Password:           c.Password,
		DB:                 c.DB,
		MaxRetries:         c.MaxRetries,
		MinRetryBackoff:    c.MinRetryBackoff,
		MaxRetryBackoff:    c.MaxRetryBackoff,
		DialTimeout:        c.DialTimeout,
		ReadTimeout:        c.ReadTimeout,
		WriteTimeout:       c.WriteTimeout,
		PoolSize:           c.PoolSize,
		MinIdleConns:       c.MinIdleConns,
		MaxConnAge:         c.MaxConnAge,
		PoolTimeout:        c.PoolTimeout,
		IdleTimeout:        c.IdleTimeout,
		IdleCheckFrequency: c.IdleCheckFrequency,
	})
	_, err := cli.Ping(context.Background()).Result()
	if err == nil {
		Global.Redis = cli
	}
	return cli, err
}
