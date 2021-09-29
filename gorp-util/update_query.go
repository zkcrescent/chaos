package gorpUtil

import (
	"fmt"
	"strings"

	"github.com/juju/errors"
)

type UpdateQuery struct {
	model      Model
	rels       []*Relation
	conditions []*Condition
	sets       []*Condition
}

func Update(m Model) *UpdateQuery {
	return &UpdateQuery{
		model: m,
	}
}

func (q *UpdateQuery) Rel(m Model, edge string) *UpdateQuery {
	if q.rels == nil {
		q.rels = []*Relation{}
	}
	q.rels = append(q.rels, R(m, edge))
	return q
}

func (q *UpdateQuery) RelWith(m Model, c *Condition) *UpdateQuery {
	if q.rels == nil {
		q.rels = []*Relation{}
	}
	q.rels = append(q.rels, RWith(m, c))
	return q
}

func (q *UpdateQuery) Rels(rs ...*Relation) *UpdateQuery {
	if q.rels == nil {
		q.rels = rs
	} else {
		q.rels = append(q.rels, rs...)
	}
	return q
}

func (q *UpdateQuery) AppendWhereIf(ok bool, cs ...*Condition) *UpdateQuery {
	if ok {
		q.Where(cs...)
	}
	return q
}

func (q *UpdateQuery) Where(cs ...*Condition) *UpdateQuery {
	if q.conditions == nil {
		q.conditions = []*Condition{}
	}
	q.conditions = append(q.conditions, cs...)
	return q
}

func (q *UpdateQuery) Set(cs ...*Condition) *UpdateQuery {
	if q.sets == nil {
		q.sets = []*Condition{}
	}
	q.sets = append(q.sets, cs...)
	return q
}

func (q *UpdateQuery) ValQuery() (string, []interface{}, error) {
	return q.fieldQuery(false)
}

func (q *UpdateQuery) Query() (string, error) {
	query, _, err := q.fieldQuery(true)
	return query, err
}

func (q *UpdateQuery) fieldQuery(withVal bool) (string, []interface{}, error) {
	if q.model == nil {
		return "", nil, ErrNilModel
	}
	if q.conditions == nil || len(q.conditions) == 0 {
		if q.rels == nil || len(q.rels) == 0 {
			return "", nil, ErrEmptyConditions
		}
	}
	if q.sets == nil || len(q.sets) == 0 {
		return "", nil, ErrEmptySets
	}
	var (
		query  string
		rc     []*Condition
		vals   []interface{}
		tables []string
		joins  []string
		tn     = q.model.TableName()

		err error
	)
	if q.rels != nil {
		if err := Relations(q.rels).Conditions(q.model, &tables, &joins, &rc); err != nil {
			return "", nil, errors.Trace(err)
		}
		if len(tables) > 0 {
			query = fmt.Sprintf(",%v", strings.Join(tables, ","))
		}
		if len(joins) > 0 {
			tn = fmt.Sprintf("%v %v", tn, strings.Join(joins, " "))
		}
	}
	if q.conditions != nil && len(q.conditions) > 0 {
		rc = append(rc, q.conditions...)
	}

	fields := ""
	for i, c := range q.sets {
		var (
			cs    string
			_vals []interface{}
		)
		_c := c.Fork().SetConnector(",")
		if withVal {
			cs, err = _c.String(i)
		} else {
			cs, _vals, err = _c.ValString(i)
			vals = append(vals, _vals...)
		}
		if err != nil {
			return "", nil, err
		}

		fields = fmt.Sprintf("%v%v", fields, cs)
	}

	condis := ""
	if len(rc) > 0 {
		condis = " where "
	}
	for i, c := range rc {
		var (
			cs    string
			_vals []interface{}
		)
		if withVal {
			cs, err = c.String(i)
		} else {
			cs, _vals, err = c.ValString(i)
			vals = append(vals, _vals...)
		}
		if err != nil {
			return "", nil, err
		}

		condis = fmt.Sprintf("%v%v", condis, cs)
	}

	query = fmt.Sprintf("update %v%v set %v%v", tn, query, fields, condis)
	return query, vals, nil
}
