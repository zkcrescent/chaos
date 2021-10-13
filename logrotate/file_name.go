package logrotate

import (
	"strconv"
	"strings"
)

type fileName struct {
	base      string
	extension string
}

func toFileName(n string) *fileName {
	if n == "" {
		return nil
	}

	fn := &fileName{}

	seps := strings.Split(n, ".")
	if sl := len(seps); sl > 1 && !(sl == 2 && seps[0] == "") {
		fn.extension = seps[sl-1]
		fn.base = strings.Join(seps[0:sl-1], ".")
	} else {
		fn.base = n
	}

	return fn
}

func (f *fileName) String() string {
	return f.appendExtension(f.base)
}

func (f *fileName) StringInNumber(i int) string {
	return f.appendExtension(f.base + "." + strconv.Itoa(i))
}

func (f *fileName) appendExtension(r string) string {
	if f.extension != "" {
		r += "." + f.extension
	}
	return r
}
