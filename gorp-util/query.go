package gorpUtil

import (
	"fmt"
	"strings"

	"github.com/juju/errors"
	"gopkg.in/gorp.v2"
)

var (
	ErrNilSqlExecutor  = errors.New("nil sqlExecutor")
	ErrNilModel        = errors.New("nil model")
	ErrEmptyConditions = errors.New("empty conditions")
	ErrEmptySets       = errors.New("empty sets")
	ErrRelNotFound     = errors.New("found no rel table")
	ErrNilFields       = errors.New("Need fields when without model")
)

const (
	EQ    = "="
	NE    = "<>"
	GT    = ">"
	GTE   = ">="
	LT    = "<"
	LTE   = "<="
	LIKE  = "like"
	IN    = "in"
	NOTIN = "not in"
	IS    = "is"
	ISNOT = "is not"
)

const (
	ASC  = "asc"
	DESC = "desc"
)

type QueryCol map[*Field]interface{}

type queryCol struct {
	field  *Field
	holder interface{}
}

type Query struct {
	onlyFields         []*queryCol // 自选映射字段
	whereFields        []*queryCol // 自选映射时自配WHERE
	model              Model
	rels               []*Relation
	conditions         []*Condition
	fields             []string
	orderBy            []string
	groupBy            []string
	withoutModelFields bool
	forUpdate          bool
	limit              int
	offset             int
}

func Get(m Model) *Query {
	return &Query{
		model: m,
	}
}

func Rel(m Model, edge string) *Query {
	return &Query{
		rels: []*Relation{
			{model: m, edge: edge},
		},
	}
}

func (q *Query) Fork() *Query {
	nq := *q
	return &nq
}

// Only for query given field, args must be Find, hoder,Field, hoder...
func (q *Query) Only(cols QueryCol, where ...QueryCol) *Query {
	var w QueryCol
	if len(where) > 0 {
		w = where[0]
	}
	for k, v := range cols {
		q.onlyFields = append(q.onlyFields, &queryCol{
			field:  k,
			holder: v,
		})
	}
	for k, v := range w {
		q.whereFields = append(q.whereFields, &queryCol{
			field:  k,
			holder: v,
		})
	}

	return q
}

func (q *Query) Pagination(page, size int) *Query {
	q.Offset((page - 1) * size).Limit(size)
	return q
}

func (q *Query) WithoutModelFields() *Query {
	q.withoutModelFields = true
	return q
}

func (q *Query) OrderBy(field *Field, typ string) *Query {
	return q.OrderByString(field.String(), typ)
}

func (q *Query) OrderByString(field, typ string) *Query {
	if q.orderBy == nil {
		q.orderBy = []string{}
	}
	q.orderBy = append(q.orderBy, fmt.Sprintf("%v %v", field, typ))
	return q
}

func (q *Query) GroupBy(fs ...*Field) *Query {
	for _, f := range fs {
		q.GroupByString(f.String())
	}
	return q
}

func (q *Query) GroupByString(fs ...string) *Query {
	if q.groupBy == nil {
		q.groupBy = []string{}
	}
	q.groupBy = append(q.groupBy, fs...)
	return q
}

func (q *Query) ForUpdate() *Query {
	q.forUpdate = true
	return q
}

func (q *Query) Limit(l int) *Query {
	q.limit = l
	return q
}

func (q *Query) Offset(o int) *Query {
	q.offset = o
	return q
}

func (q *Query) Get(m Model) *Query {
	q.model = m
	return q
}

func (q *Query) Rel(m Model, edge string) *Query {
	if q.rels == nil {
		q.rels = []*Relation{}
	}
	q.rels = append(q.rels, R(m, edge))
	return q
}

func (q *Query) RelWith(m Model, c *Condition) *Query {
	if q.rels == nil {
		q.rels = []*Relation{}
	}
	q.rels = append(q.rels, RWith(m, c))
	return q
}

func (q *Query) Rels(rs ...*Relation) *Query {
	if q.rels == nil {
		q.rels = rs
	} else {
		q.rels = append(q.rels, rs...)
	}
	return q
}

func (q *Query) AppendWhereIf(ok bool, cs ...*Condition) *Query {
	if ok {
		q.Where(cs...)
	}
	return q
}

func (q *Query) Where(cs ...*Condition) *Query {
	if q.conditions == nil {
		q.conditions = []*Condition{}
	}
	q.conditions = append(q.conditions, cs...)
	return q
}

func (q *Query) Fields(fs ...string) *Query {
	if q.fields == nil {
		q.fields = []string{}
	}
	q.fields = append(q.fields, fs...)
	return q
}

func (q *Query) ifAdd(query string, v []string, s string) string {
	if v != nil && len(v) > 0 {
		query = fmt.Sprintf("%v %v %v", query, s, strings.Join(v, ","))
	}
	return query
}

func (q *Query) modelFields(prefix ...string) []string {
	fs := q.model.Fields()
	if len(prefix) > 0 {
		for i, v := range fs {
			fs[i] = fmt.Sprintf("%v%v", prefix[0], v)
		}
	}
	return fs
}

func (q *Query) QueryFields() ([]string, error) {
	var fields []string

	if len(q.onlyFields) != 0 {
		for _, v := range q.onlyFields {
			fields = append(fields, v.field.String())
		}
		return fields, nil
	}

	if !q.withoutModelFields {
		fields = q.model.Fields()
	}
	if q.fields != nil && len(q.fields) > 0 {
		fields = append(fields, q.fields...)
	} else if q.withoutModelFields {
		return nil, ErrNilFields
	}
	return fields, nil
}

func (q *Query) Query() (string, error) {
	fields, err := q.QueryFields()
	if err != nil {
		return "", err
	}
	v, _, err := q.fieldQuery(true, fields...)
	return v, err
}

func (q *Query) ValQuery() (string, []interface{}, error) {
	return q.valQuery()
}

func (q *Query) valQuery() (string, []interface{}, error) {
	fields, err := q.QueryFields()
	if err != nil {
		return "", nil, err
	}
	return q.fieldQuery(false, fields...)
}

func (q *Query) whereQuery(withVal bool) (string, []interface{}, error) {
	var rc []*Condition
	if len(q.conditions) > 0 {
		rc = append(rc, q.conditions...)
	} else {
		return "", nil, nil
	}
	var query string
	var vals []interface{}
	for i, c := range rc {
		var (
			cs    string
			_vals []interface{}
			err   error
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

		query = fmt.Sprintf("%v%v", query, cs)
	}
	return query, vals, nil
}

func (q *Query) fieldQuery(withVal bool, fields ...string) (string, []interface{}, error) {
	if q.model == nil {
		return "", nil, ErrNilModel
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
	if len(q.conditions) > 0 {
		rc = append(rc, q.conditions...)
	}
	if len(rc) > 0 {
		query = fmt.Sprintf("%v where ", query)
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

		query = fmt.Sprintf("%v%v", query, cs)
	}
	query = q.ifAdd(query, q.groupBy, "group by")
	query = q.ifAdd(query, q.orderBy, "order by")
	if q.limit > 0 {
		query = fmt.Sprintf("%v limit %v", query, q.limit)
		if q.offset > 0 {
			query = fmt.Sprintf("%v offset %v", query, q.offset)
		}
	}
	query = fmt.Sprintf("select %v from %v%v", strings.Join(fields, ","), tn, query)
	if q.forUpdate {
		query += " for update"
	}
	return query, vals, nil
}

func (q *Query) Fetch(db gorp.SqlExecutor, placeholders ...interface{}) error {
	if db == nil {
		return ErrNilSqlExecutor
	}
	var holder interface{}
	if len(placeholders) > 0 {
		holder = placeholders[0]
	} else {
		holder = q.model
	}

	query, vals, err := q.valQuery()
	if err != nil {
		return err
	}
	//fmt.Println("fetch", query)
	if len(q.onlyFields) != 0 {
		err = q.queryRow(db, query, vals)
	} else {
		err = db.SelectOne(holder, query, vals...)
	}

	if err != nil {
		return q.QueryValError(err, query, vals)
	}

	return err
}

func (q *Query) Update(db gorp.SqlExecutor) error {
	if len(q.onlyFields) == 0 {
		_, err := db.Update(q.model)
		return err
	}
	var fields []string
	var args []interface{}
	var query []string
	for _, v := range q.onlyFields {
		args = append(args, v.holder)
		fields = append(fields, fmt.Sprintf("`%s`= ?", v.field.String()))
	}

	for _, v := range q.whereFields {
		args = append(args, v.holder)
		query = append(query, fmt.Sprintf("`%s`= ?", v.field.String()))
	}
	where, holders, err := q.whereQuery(false)
	if err != nil {
		return err
	}
	args = append(holders, args...)
	var sql string
	if where != "" {
		sql = fmt.Sprintf("UPDATE %s SET %s WHERE %s AND %s", q.model.TableName(), strings.Join(fields, ","), where, strings.Join(query, "AND"))
	} else {
		sql = fmt.Sprintf("UPDATE %s SET %s WHERE %s", q.model.TableName(), strings.Join(fields, ","), strings.Join(query, "AND"))
	}
	if q.model.VersionField() != "" {
		sql = fmt.Sprintf("%s AND `%s`=%v", sql, q.model.VersionField(), q.model.Version())
	}

	if _, err := db.Exec(sql, args...); err != nil {
		return q.QueryError(err, sql)
	}

	return nil

}

func (q *Query) queryRow(db gorp.SqlExecutor, query string, args []interface{}) error {
	row := db.QueryRow(query, args...)
	var res []interface{}
	for _, v := range q.onlyFields {
		res = append(res, v.holder)
	}
	return row.Scan(res...)
}

func (q *Query) QueryError(err error, qs string) error {
	return errors.Annotatef(err, "query: %v", qs)
}

func (q *Query) QueryValError(err error, qs string, vals []interface{}) error {
	return errors.Annotatef(err, "query: %v, vals: %v", qs, vals)
}

func (q *Query) CountQuery(fields ...string) (string, error) {
	query, _, err := q.countQuery(true, fields...)
	return query, err
}

func (q *Query) countQuery(withVal bool, fields ...string) (string, []interface{}, error) {
	pk, _ := q.model.PK()
	f := pk
	if len(fields) > 0 {
		f = TableField("", fields[0])
	}
	return q.fieldQuery(withVal, f.Count().String())
}

func (q *Query) Count(db gorp.SqlExecutor, fields ...string) (int64, error) {
	if db == nil {
		return 0, ErrNilSqlExecutor
	}
	query, vals, err := q.countQuery(false, fields...)
	if err != nil {
		return 0, errors.Trace(err)
	}
	//fmt.Println("count", query)
	v, err := db.SelectInt(query, vals...)
	if err != nil {
		return v, q.QueryValError(err, query, vals)
	}
	return v, err
}

func (q *Query) FetchAll(db gorp.SqlExecutor, placeholders ...interface{}) ([]interface{}, error) {
	if db == nil {
		return nil, ErrNilSqlExecutor
	}
	var holder interface{} = q.model
	if len(placeholders) > 0 {
		holder = placeholders[0]
	}

	fields, err := q.QueryFields()
	if err != nil {
		return nil, err
	}
	query, vals, err := q.fieldQuery(false, fields...)
	if err != nil {
		return nil, err
	}

	//fmt.Println("fetchall", query)
	v, err := db.Select(holder, query, vals...)
	if err != nil {
		return v, q.QueryValError(err, query, vals)
	}

	return v, err
}
