package gorpUtil

import (
	"encoding/json"
	"fmt"
	"time"

	null "gopkg.in/nullbio/null.v6"
)

// Abstract
type Base struct {
	ID          int64  `db:"id" json:"id"`
	CreatedTime Time   `db:"created_time" json:"created_time"`
	UpdatedTime Time   `db:"updated_time" json:"updated_time"`
	UpdatedSeq  int64  `db:"updated_seq" json:"updated_seq"` // 乐观锁标记
	Region      string `db:"region" json:"region"`
}

// Abstract
type EnableBase struct {
	Base
	RemovedTime Time `db:"removed_time" json:"removed_time"`
	Removed     bool `db:"removed" json:"removed"` // 软删标记
}

// Abstract
type TaskEnableBase struct {
	EnableBase
	TaskBase
}

type TaskBase struct {
	TaskStatus  string `db:"task_status" json:"task_status"`
	TaskRetries int64  `db:"task_retries" json:"task_retries"`
	Operator    string `db:"operator" json:"operator"`
}

type FullTaskBase struct {
	TaskStatus   string `db:"task_status" json:"task_status"`
	TaskRetries  int64  `db:"task_retries" json:"task_retries"`
	Operator     string `db:"operator" json:"operator"`
	OperatorName string `db:"operator_name" json:"operator_name"`
}

func TimeFrom(t time.Time) Time {
	return Time{null.TimeFrom(t)}
}

func Now() Time {
	return Time{null.TimeFrom(time.Now())}
}

func NowInSecond() Time {
	now := time.Now()
	return Time{null.TimeFrom(time.Unix(now.Unix(), 0))}
}

const format = "2006-01-02 15:04:05"

const formatMs = "2006-01-02 15:04:05.999"

type Time struct {
	null.Time
}

func (t Time) String() string {
	if !t.Valid {
		return string(null.NullBytes)
	}
	return t.Time.Time.Format(format)
}
func (t Time) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return null.NullBytes, nil
	}
	return json.Marshal(t.Time.Time.Format(format))
}
func (t *Time) UnmarshalJSON(b []byte) error {
	if b == nil || len(b) == 0 || string(b) == string(null.NullBytes) {
		t.Valid = false
		return nil
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	loc, _ := time.LoadLocation("Local")
	val, err := time.ParseInLocation(format, s, loc)
	if err != nil {
		val, err = time.Parse(time.RFC3339, s)
		if err != nil {
			return err
		}
	}
	t.SetValid(val)
	return nil
}

func (t *Time) SetValid(v time.Time) {
	if v.Unix() <= 0 {
		return
	}
	t.Time.SetValid(v)
}

func (t *Time) Scan(value interface{}) error {
	var err error
	switch x := value.(type) {
	case time.Time:
		t.SetValid(x)
	case nil:
		t.Valid = false
		return nil
	default:
		err = fmt.Errorf("null: cannot scan type %T into null.Time: %v", value, value)
	}
	return err
}

type TimeMs struct {
	Time
}

func (t TimeMs) String() string {
	if !t.Valid {
		return string(null.NullBytes)
	}
	return t.Time.Time.Time.Format(formatMs)
}
func (t TimeMs) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return null.NullBytes, nil
	}
	return json.Marshal(t.Time.Time.Time.Format(formatMs))
}
func (t *TimeMs) UnmarshalJSON(b []byte) error {
	if b == nil || len(b) == 0 || string(b) == string(null.NullBytes) {
		t.Valid = false
		return nil
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	loc, _ := time.LoadLocation("Local")
	val, err := time.ParseInLocation(formatMs, s, loc)
	if err != nil {
		val, err = time.Parse(time.RFC3339, s)
		if err != nil {
			return err
		}
	}
	t.SetValid(val)
	return nil
}
