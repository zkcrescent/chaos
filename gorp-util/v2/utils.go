package gorpUtil

import (
	"fmt"

	"github.com/juju/errors"
	gorp "gopkg.in/gorp.v2"
)

func tryToDumpSource(v interface{}) {
	if m, ok := v.(ModelV2); ok && m.GenVersion() == GenVersion_V2 {
		m.DumpSource()
	}
}

func ColumnMapFilter(fields map[string]interface{}, versionField string) gorp.ColumnFilter {
	return func(m *gorp.ColumnMap) bool {
		if versionField == m.ColumnName && versionField != "" {
			return true
		}

		_, ok := fields[m.ColumnName]
		return ok
	}
}

func UpdateColumns(db gorp.SqlExecutor, m ModelV2) (int64, error) {
	valueMap := m.UpdateColumnsFields()
	if len(valueMap) == 0 {
		return 0, nil
	}
	query, vals := updateColumnsQuery(m, valueMap)
	r, err := db.Exec(query, vals...)
	if err != nil {
		return 0, errors.Annotatef(err, "query: %s, vals: %v", query, vals)
	}
	return r.RowsAffected()
}

func UpdateColumnsQuery(m ModelV2) string {
	valueMap := m.UpdateColumnsFields()
	q, _ := updateColumnsQuery(m, valueMap)
	return q
}

func updateColumnsQuery(m ModelV2, valueMap map[string]interface{}) (string, []interface{}) {
	pkf, pk := m.PK()
	vf, v := m.VersionK()
	vals := []interface{}{}
	tableName := m.TableName()

	query := fmt.Sprintf("UPDATE `%v` SET", tableName)
	x := 0
	for k, v := range valueMap {
		if x == 0 {
			query += " "
		} else {
			query += ", "
		}
		query += fmt.Sprintf("`%v`=?", k)
		vals = append(vals, formatValue(v))
		x++
	}
	query += fmt.Sprintf(" WHERE `%v`=? AND `%v`=?", pkf.Key(), vf.Key())
	vals = append(vals, pk, v)
	return query, vals
}

func formatValue(v interface{}) string {
	needQuote := true
	if _, ok := v.(int); ok {
		needQuote = ok
	}
	if _, ok := v.(bool); ok {
		needQuote = ok
	}
	if needQuote {
		return fmt.Sprintf("'%v'", v)
	}
	return fmt.Sprintf("%v", v)
}
