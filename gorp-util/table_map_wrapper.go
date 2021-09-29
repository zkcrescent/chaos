package gorpUtil

import (
	"bytes"
	"fmt"
	"reflect"

	gorp "gopkg.in/gorp.v2"
)

const (
	MysqlIndexTypeBtree = "Btree"
	MysqlIndexTypeHash  = "Hash"
)

type TableMap struct {
	*gorp.TableMap

	indexs []*IndexMap
}

func (tm *TableMap) AddIndex(name string, idxtype string, columns []string) *IndexMap {
	im := &IndexMap{
		IndexMap: tm.TableMap.AddIndex(name, idxtype, columns),
		columns:  columns,
	}
	tm.indexs = append(tm.indexs, im)
	return im
}

func (tm *TableMap) AllCreateIndexQueries(m *gorp.DbMap) []string {
	queries := []string{}
	for _, im := range tm.indexs {
		queries = append(queries, tm.CreateIndexQuery(m, im))
	}
	return queries
}

func (tm *TableMap) CreateIndexQuery(m *gorp.DbMap, index *IndexMap) string {
	dialect := reflect.TypeOf(m.Dialect)
	s := bytes.Buffer{}
	s.WriteString("create")
	if index.Unique {
		s.WriteString(" unique")
	}
	s.WriteString(" index")
	s.WriteString(fmt.Sprintf(" %s on %s", index.IndexName, tm.TableName))
	if dname := dialect.Name(); dname == "PostgresDialect" && index.IndexType != "" {
		s.WriteString(fmt.Sprintf(" %s %s", m.Dialect.CreateIndexSuffix(), index.IndexType))
	}
	s.WriteString(" (")
	for x, col := range index.columns {
		if x > 0 {
			s.WriteString(", ")
		}
		s.WriteString(m.Dialect.QuoteField(col))
	}
	s.WriteString(")")

	if dname := dialect.Name(); dname == "MySQLDialect" && index.IndexType != "" {
		s.WriteString(fmt.Sprintf(" %s %s", m.Dialect.CreateIndexSuffix(), index.IndexType))
	}
	s.WriteString(";")
	return s.String()
}

type IndexMap struct {
	*gorp.IndexMap

	columns []string
}
