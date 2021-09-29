package gorpUtil

import (
	"fmt"
	"strings"
)

const (
	Connector_AND = "and"
	Connector_OR  = "or"
)

type Condition struct {
	key        *string
	operator   string
	value      interface{}
	connector  string
	conditions []*Condition
	rel        bool
}

func Case(key, operator string, value interface{}) *Condition {
	return &Condition{
		key:       &key,
		operator:  operator,
		value:     value,
		connector: Connector_AND,
	}
}

func OrCase(key, operator string, value interface{}, cs ...*Condition) *Condition {
	return &Condition{
		key:        &key,
		operator:   operator,
		value:      value,
		connector:  Connector_OR,
		conditions: cs,
	}
}

func caseCombine(connector string, cs ...*Condition) *Condition {
	if len(cs) == 1 {
		cs[0].connector = connector
		return cs[0]
	}
	return &Condition{
		connector:  connector,
		conditions: cs,
	}
}

func Or(cs ...*Condition) *Condition {
	return caseCombine(Connector_OR, cs...)
}

func And(cs ...*Condition) *Condition {
	return caseCombine(Connector_AND, cs...)
}

func (c *Condition) ValString(index int) (string, []interface{}, error) {
	var (
		v    string = "?"
		vals []interface{}
		err  error
	)
	if c.value != nil {
		query, isQuery := c.value.(*Query)
		fp, isFieldPtr := c.value.(*Field)
		arr, isArray := c.value.([]interface{})
		if isQuery {
			v, vals, err = query.valQuery()
			if err != nil {
				return "", nil, err
			}
			v = "(" + v + ")"
		} else if isFieldPtr {
			v = fp.String()
		} else if isArray {
			v = fmt.Sprintf("(%v)", strings.Repeat("?,", len(arr))[0:len(arr)*2-1])
			vals = append(vals, arr...)
		} else {
			vals = append(vals, c.value)
		}
	}

	var s string
	if c.key != nil {
		s = fmt.Sprintf("%v %v %v", *c.key, c.operator, v)
	}
	if c.conditions != nil && len(c.conditions) > 0 {
		for i, cc := range c.conditions {
			sep := i
			if len(s) > 0 {
				sep += 1
			}
			// if i == 0 && len(s) > 0 {
			// 	sep = 1
			// }
			_s, _vals, err := cc.ValString(sep)
			if err != nil {
				return "", nil, err
			}
			s = fmt.Sprintf("%v%v", s, _s)
			vals = append(vals, _vals...)
		}
		if len(c.conditions) > 1 {
			s = fmt.Sprintf("(%v)", s)
		}
	}
	if index > 0 {
		s = fmt.Sprintf(" %v %v", c.connector, s)
	}
	return s, vals, nil
}

func (c *Condition) String(index int) (string, error) {
	var (
		v   string
		err error
	)
	if c.value != nil {
		v = fmt.Sprintf("%v", c.value)
		_, isInt := c.value.(int)
		query, isQuery := c.value.(*Query)
		if !(isInt || isQuery || c.rel || c.value == "null" || c.operator == IN || c.operator == NOTIN) {
			v = fmt.Sprintf("'%v'", c.value)
		}
		if isQuery {
			v, err = query.Query()
			if err != nil {
				return "", err
			}
			v = "(" + v + ")"
		}
	}

	var s string
	if c.key != nil {
		s = fmt.Sprintf("%v %v %v", *c.key, c.operator, v)
	}
	if c.conditions != nil && len(c.conditions) > 0 {
		for i, cc := range c.conditions {
			sep := i
			if len(s) > 0 {
				sep += 1
			}
			// if i == 0 && len(s) > 0 {
			// 	sep = 1
			// }
			_s, err := cc.String(sep)
			if err != nil {
				return "", err
			}
			s = fmt.Sprintf("%v%v", s, _s)
		}
		if len(c.conditions) > 1 {
			s = fmt.Sprintf("(%v)", s)
		}
	}
	if index > 0 {
		s = fmt.Sprintf(" %v %v", c.connector, s)
	}
	return s, nil
}

func (c *Condition) Relation() *Condition {
	c.rel = true
	return c
}

func (c *Condition) Sub(cs ...*Condition) *Condition {
	if c.conditions == nil {
		c.conditions = []*Condition{}
	}
	c.conditions = append(c.conditions, cs...)
	return c
}
