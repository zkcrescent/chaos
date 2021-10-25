package gorpUtil

import (
	gorp "gopkg.in/gorp.v2"
)

var Tables = Models{}

type TableCheck func(db *gorp.DbMap, t *TableMap, tableName string) error

type Table interface {
	TableName() string
	Fields() []string
	VersionField() string
	Version() int64
	PK() (*Field, interface{})
}

type ShardingTable interface {
	Sharding() int64
}

type Model interface {
	Table
	Relation(string) (*Field, bool)                 // fk name, exists
	Multiple(string) (string, *Field, *Field, bool) // tablename, sk name, fk name, exists
}

type Models map[string]Table

func (ms Models) Add(mm ...Table) {
	for _, m := range mm {
		ms[m.TableName()] = m
	}
}

func (ms Models) Elem() map[string]Table {
	return (map[string]Table)(ms)
}

func (ms Models) checkTable(db *gorp.DbMap, fs ...TableCheck) ([]*TableMap, error) {
	tms := []*TableMap{}
	var f TableCheck
	if len(fs) != 0 {
		f = fs[0]
	}
	for _, t := range ms {
		pk, _ := t.PK()
		tm := &TableMap{
			TableMap: db.AddTableWithName(t, t.TableName()).SetKeys(true, pk.Key()),
		}
		if v := t.VersionField(); v != "" {
			tm.SetVersionCol(v)
		}
		if f != nil {
			if err := f(db, tm, t.TableName()); err != nil {
				return nil, err
			}
		}
		tms = append(tms, tm)
	}
	return tms, nil
}

func (ms Models) CreateTableQueries(db *gorp.DbMap, skipExists bool, fs ...TableCheck) ([]string, error) {
	tms, err := ms.checkTable(db, fs...)
	if err != nil {
		return nil, err
	}
	queries := []string{}
	for _, tm := range tms {
		queries = append(queries, tm.SqlForCreate(skipExists))
		queries = append(queries, tm.AllCreateIndexQueries(db)...)
	}
	return queries, nil
}

func (ms Models) AddTables(db *gorp.DbMap, fs ...TableCheck) error {
	_, err := ms.checkTable(db, fs...)
	return err
}

func (ms Models) CheckTable(db *gorp.DbMap, fs ...TableCheck) error {
	_, err := ms.checkTable(db, fs...)
	if err != nil {
		return err
	}
	return db.CreateTablesIfNotExists()
}
