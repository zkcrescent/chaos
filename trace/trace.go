package trace

import (
	"runtime"
	"strings"
)

const (
	TRACE_ID = "trace_id"
)

func spanName() string {
	pc := make([]uintptr, 1) // at least 1 entry needed
	n := runtime.Callers(3, pc)
	if n > 0 {
		name := runtime.FuncForPC(pc[0]).Name()
		strs := strings.Split(name, "/")
		l := len(strs)
		if l >= 2 {
			return strs[l-2] + "/" + strs[l-1]
		} else {
			return name
		}
	}
	return ""
}
