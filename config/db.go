package config

import (
	"database/sql"
	"encoding/json"
	"time"

	_ "github.com/go-sql-driver/mysql"
	gorp "gopkg.in/gorp.v2"
)

// Database for basic innoDB mysql
type Database struct {
	DSN         string        `json:"dsn" yaml:"dsn" valid:"required" toml:"dsn" xml:"dsn"`
	MaxConnLife time.Duration `json:"max_conn_life" yaml:"max_conn_life" valid:"required" toml:"max_conn_life" xml:"max_conn_life"`
	MaxConn     int           `json:"max_conn" yaml:"max_conn" valid:"required" toml:"max_conn" xml:"max_conn"`
	MaxIdle     int           `json:"max_idle" yaml:"max_idle" valid:"required" toml:"max_idle" xml:"max_idle"`
}

func (c *Database) UnmarshalJSON(b []byte) error {
	d := &struct {
		DSN         string      `json:"dsn" yaml:"dsn" valid:"required" toml:"dsn" xml:"dsn"`
		MaxConnLife interface{} `json:"max_conn_life" yaml:"max_conn_life" valid:"required" toml:"max_conn_life" xml:"max_conn_life"`
		MaxConn     int         `json:"max_conn" yaml:"max_conn" valid:"required" toml:"max_conn" xml:"max_conn"`
		MaxIdle     int         `json:"max_idle" yaml:"max_idle" valid:"required" toml:"max_idle" xml:"max_idle"`
	}{}

	if err := json.Unmarshal(b, d); err != nil {
		return err
	}

	c.DSN = d.DSN
	c.MaxConn = d.MaxConn
	c.MaxIdle = d.MaxIdle

	switch d.MaxConnLife.(type) {
	case int:
		c.MaxConnLife = time.Duration(d.MaxConnLife.(int64))
	case string:
		dur, err := time.ParseDuration(d.MaxConnLife.(string))
		if err != nil {
			return err
		}
		c.MaxConnLife = dur
	}
	return nil
}

func (c Database) DB() (*gorp.DbMap, error) {
	db, err := sql.Open("mysql", c.DSN)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(c.MaxConnLife)
	db.SetMaxOpenConns(c.MaxConn)
	db.SetMaxIdleConns(c.MaxIdle)

	return &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8"}}, nil
}

func (c Database) GDB() (*gorp.DbMap, error) {
	if db, err := c.DB(); err != nil {
		return nil, err
	} else {
		Global.DB = db
	}
	return Global.DB, nil
}
