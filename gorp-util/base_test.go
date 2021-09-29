package gorpUtil

import (
	"encoding/json"
	"log"
	"testing"
)

func TestTime(t *testing.T) {
	var a = Time{}
	err := json.Unmarshal([]byte(`"0000-00-00 00:00:00"`), &a)
	if err != nil {
		panic(err)
	}
	log.Println("valid", a.Valid, "unix", a.Time.Time.UTC().Unix())

}
