package gorpUtil

import (
	"fmt"
	"strings"
)

const (
	NULL        = "null"
	Empty       = ""
	EmptyString = `""`
	SUM         = "sum"
	AVG         = "avg"
	MAX         = "max"
	MIN         = "min"
	COUNT       = "count"
	DISTINCT    = "distinct"
	IfNull      = "ifnull"
)

type Field struct {
	key    string
	table  string
	full   string
	as     *string
	method *FieldMethod
}

func ShardTableField(table string, key string, sharding int64) func(int64) *Field {
	return func(shardkey int64) *Field {
		f := &Field{
			key:   key,
			table: table,
		}
		if table != "" {
			f.full = fmt.Sprintf("%v_%v.%v", table, shardkey%shardkey, key)
		} else {
			f.full = key
		}

		return f
	}
}

func TableField(table, key string) *Field {
	f := &Field{
		key:   key,
		table: table,
	}
	if table != "" {
		f.full = fmt.Sprintf("%v.%v", table, key)
	} else {
		f.full = key
	}
	return f
}

func (f *Field) String() string {
	full := f.full
	if f.method != nil {
		full = f.method.Format(full)
	}
	if f.as != nil {
		full = fmt.Sprintf("%v as %v", full, *f.as)
	}
	return full
}

func (f *Field) Fork() *Field {
	nf := *f
	return &nf
}

func (f *Field) Key() string {
	return f.key
}

func (f *Field) Table() string {
	return f.table
}

func (f *Field) Count() *Field {
	return f.Method(COUNT)
}

func (f *Field) Sum() *Field {
	return f.Method(SUM)
}

func (f *Field) Avg() *Field {
	return f.Method(AVG)
}

func (f *Field) Max() *Field {
	return f.Method(MAX)
}

func (f *Field) Min() *Field {
	return f.Method(MIN)
}

func (f *Field) Distinct() *Field {
	return f.Method(DISTINCT)
}

func (f *Field) IfNull(val interface{}) *Field {
	return f.Method(IfNull, val)
}

func (f *Field) Method(m string, args ...interface{}) *Field {
	nf := f.Fork()
	nf.method = &FieldMethod{
		Method: m,
		Args:   args,
	}
	return nf
}

func (f *Field) AS(an string) *Field {
	nf := f.Fork()
	nf.as = &an
	return nf
}

func (f *Field) EQ(value interface{}) *Condition {
	if b, ok := value.(bool); ok {
		value = 0
		if b {
			value = 1
		}
	}
	return Case(f.full, EQ, value)
}

func (f *Field) NE(value interface{}) *Condition {
	return Case(f.full, NE, value)
}

func (f *Field) GT(value interface{}) *Condition {
	return Case(f.full, GT, value)

}
func (f *Field) GTE(value interface{}) *Condition {
	return Case(f.full, GTE, value)
}

func (f *Field) LT(value interface{}) *Condition {
	return Case(f.full, LT, value)
}

func (f *Field) LTE(value interface{}) *Condition {
	return Case(f.full, LTE, value)
}

func (f *Field) Like(value interface{}) *Condition {
	return Case(f.full, LIKE, value)
}

func (f *Field) arrayString(vals ...interface{}) string {
	allstr := true
	for _, val := range vals {
		if _, ok := val.(string); !ok {
			allstr = false
		}
	}
	tmp := make([]string, len(vals))
	for i, v := range vals {
		if allstr {
			tmp[i] = fmt.Sprintf("'%v'", v)
		} else {
			tmp[i] = fmt.Sprintf("%v", v)
		}
	}
	return fmt.Sprintf("(%s)", strings.Join(tmp, ","))
}

func (f *Field) IN(vals ...interface{}) *Condition {
	// return f.INStr(f.arrayString(vals...))
	return Case(f.full, IN, vals)
}

func (f *Field) NOTIN(vals ...interface{}) *Condition {
	// return f.NOTINStr(f.arrayString(vals...))
	return Case(f.full, NOTIN, vals)
}

// func (f *Field) INStr(value string) *Condition {
// 	return Case(f.full, IN, value)
// }
//
// func (f *Field) NOTINStr(value string) *Condition {
// 	return Case(f.full, NOTIN, value)
// }

func (f *Field) subQueryString(value string) string {
	return fmt.Sprintf("(%s)", value)
}

func (f *Field) INSubQuery(value *Query) *Condition {
	return Case(f.full, IN, value)
}

func (f *Field) NOTINSubQuery(value *Query) *Condition {
	return Case(f.full, NOTIN, value)
}

func (f *Field) IS(value interface{}) *Condition {
	return Case(f.full, IS, value)
}

func (f *Field) ISEmpty() *Condition {
	return Case(f.full, EQ, Empty)
}

func (f *Field) ISNOTEmpty() *Condition {
	return Case(f.full, NE, Empty)
}

func (f *Field) ISNULL() *Condition {
	return f.IS(NULL)
}
func (f *Field) ISNOTNULL() *Condition {
	return f.ISNOT(NULL)
}

func (f *Field) ISNOT(value interface{}) *Condition {
	return Case(f.full, ISNOT, value)
}

type FieldMethod struct {
	Method string
	Args   []interface{}
}

func (m *FieldMethod) Format(full string) string {
	args := ""
	if m.Args != nil && len(m.Args) > 0 {
		args = fmt.Sprintf(""+strings.Repeat(", %v", len(m.Args)), m.Args...)
	}
	return fmt.Sprintf("%v(%v%v)", m.Method, full, args)
}
