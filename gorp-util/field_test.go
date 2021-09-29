package gorpUtil

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestArrayString(t *testing.T) {
	f := TableField("table", "field")
	fmt.Println(f.IN(1, 2, 3).String(0))
	fmt.Println(f.NOTIN(1, 2, 3).String(0))
	fmt.Println(f.Method("ifNull", `""`))
	fmt.Println(f.Max())
}

func TestGuMs(t *testing.T) {
	var tt = struct {
		T TimeMs `json:"t"`
	}{
		T: TimeMs{
			Now(),
		},
	}

	bts, _ := json.Marshal(tt)

	fmt.Println("json:", string(bts))

	var ttt = struct {
		T TimeMs `json:"t"`
	}{}

	err := json.Unmarshal(bts, &ttt)

	fmt.Println("json unmarshal:", err, ttt.T.Time.Time.Time)
}
