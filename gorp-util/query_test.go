package gorpUtil

import (
	"fmt"
	"testing"
)

var (
	TableA_ID    = TableField("table_a", "id")
	TableA_Field = TableField("table_a", "field")

	TableB_ID    = TableField("table_b", "id")
	TableB_Field = TableField("table_b", "field")
	TableB_AID   = TableField("table_b", "a_id")

	TableC_ID    = TableField("table_c", "id")
	TableC_Field = TableField("table_c", "field")

	TableD_ID    = TableField("table_d", "id")
	TableD_Field = TableField("table_d", "field")

	TableE_ID    = TableField("table_e", "id")
	TableE_Field = TableField("table_e", "field")
	TableE_CID   = TableField("table_e", "c_id")

	AC_AID = TableField("a_c", "a_id")
	AC_CID = TableField("a_c", "c_id")

	CD_CID = TableField("c_d", "c_id")
	CD_DID = TableField("c_d", "d_id")
)

type TableA struct {
	ID    int64
	Field string
}

func (t *TableA) TableName() string {
	return "table_a"
}

func (t *TableA) VersionField() string {
	return ""
}

func (t *TableA) Version() int64 {
	return 0
}

func (t *TableA) Fields() []string {
	return []string{
		"table_a.id",
		"table_a.field",
	}
}

func (t *TableA) PK() (*Field, interface{}) {
	return TableA_ID, t.ID
}

func (t *TableA) Relation(edge string) (*Field, bool) {
	fm := map[string]*Field{
		"a_b": TableB_AID,
	}
	fk, ok := fm[edge]
	return fk, ok
}

func (t *TableA) Multiple(edge string) (string, *Field, *Field, bool) {
	switch edge {
	case "a_c":
		return "a_c", AC_AID, AC_CID, true
	}
	return "", nil, nil, false
}

type TableB struct {
	ID    int64
	Field string
}

func (t *TableB) TableName() string {
	return "table_b"
}

func (t *TableB) VersionField() string {
	return ""
}

func (t *TableB) Fields() []string {
	return []string{
		"table_b.id", "table_b.field",
	}
}

func (t *TableB) PK() (*Field, interface{}) {
	return TableB_ID, t.ID
}

func (t *TableB) Relation(edge string) (*Field, bool) {
	return nil, false
}

func (t *TableB) Multiple(edge string) (string, *Field, *Field, bool) {
	return "", nil, nil, false
}

type TableC struct {
	ID    int64
	Field string
}

func (t *TableC) TableName() string {
	return "table_c"
}

func (t *TableC) VersionField() string {
	return ""
}

func (t *TableC) Fields() []string {
	return []string{
		"table_c.id", "table_c.field",
	}
}

func (t *TableC) PK() (*Field, interface{}) {
	return TableC_ID, t.ID
}

func (t *TableC) Relation(edge string) (*Field, bool) {
	fm := map[string]*Field{
		"c_e": TableE_CID,
	}
	fk, ok := fm[edge]
	return fk, ok
}

func (t *TableC) Multiple(edge string) (string, *Field, *Field, bool) {
	switch edge {
	// case "a_c":
	// 	return "a_c", AC_CID, AC_AID, true
	case "c_d":
		return "c_d", CD_CID, CD_DID, true
	}
	return "", nil, nil, false
}

type TableD struct {
	ID    int64
	Field string
}

func (t *TableD) TableName() string {
	return "table_d"
}

func (t *TableD) VersionField() string {
	return ""
}

func (t *TableD) Fields() []string {
	return []string{
		"table_d.id", "table_d.field",
	}
}

func (t *TableD) PK() (*Field, interface{}) {
	return TableD_ID, t.ID
}

func (t *TableD) Relation(edge string) (*Field, bool) {
	return nil, false
}

func (t *TableD) Multiple(edge string) (string, *Field, *Field, bool) {
	switch edge {
	case "c_d":
		return "c_d", CD_DID, CD_CID, true
	}
	return "", nil, nil, false
}

type TableE struct {
	ID    int64
	Field string
}

func (t *TableE) TableName() string {
	return "table_e"
}

func (t *TableE) VersionField() string {
	return ""
}

func (t *TableE) Fields() []string {
	return []string{
		"table_e.id", "table_e.field",
	}
}

func (t *TableE) PK() (*Field, interface{}) {
	return TableE_ID, t.ID
}

func (t *TableE) Relation(edge string) (*Field, bool) {
	return nil, false
}

func (t *TableE) Multiple(edge string) (string, *Field, *Field, bool) {
	return "", nil, nil, false
}

func Test_Query(t *testing.T) {
	qs := Get(&TableA{}).Where(TableA_Field.EQ("a"), Or(TableA_Field.IS("null")))
	fmt.Println(
		qs.Offset(10).Limit(10).GroupBy(TableA_Field, TableA_ID).OrderBy(TableA_Field, ASC).OrderBy(TableA_ID, DESC).Query(),
	)
	fmt.Println(
		qs.CountQuery(),
	)
	fmt.Println(
		Get(&TableA{}).Where(Or(TableA_Field.EQ("b"), TableA_Field.EQ("c"))).Query(),
	)
	fmt.Println(
		Get(&TableA{}).Where(TableA_Field.EQ("a"), Or(TableA_Field.EQ("b"), TableA_Field.EQ("c"))).Query(),
	)
	fmt.Println(
		Get(&TableA{}).Where(TableA_Field.EQ("a"), OrCase(TableA_Field.String(), "=", "b", TableA_Field.EQ("c"))).Query(),
	)
	fmt.Println(
		Get(&TableB{}).Rel(&TableA{}, "a_b").Where(TableB_Field.EQ("b"), TableB_Field.EQ("testapp")).Query(),
	)
	fmt.Println(
		Rel(&TableA{ID: 1}, "a_c").Get(&TableC{}).Where(TableC_Field.EQ("c")).Fields(TableA_Field.Count().AS("a_field").String()).Query(),
	)
	fmt.Println(
		Rel(&TableC{ID: 1}, "a_c").Get(&TableA{}).Where(TableA_Field.EQ("a")).Fields(TableC_Field.Count().AS("c_field").String()).Query(),
	)
	fmt.Println(
		Get(&TableA{}).Rels(
			R(&TableC{}, "a_c",
				R(&TableD{}, "c_d"),
				R(&TableE{}, "c_e"),
			),
		).Where(TableC_Field.EQ("c")).Query(),
	)
	fmt.Println(
		Get(&TableA{}).Rels(
			R(&TableC{}, "a_c",
				R(&TableD{}, "c_d"),
				R(&TableE{}, "c_e").Join(Join_Left),
			).Join(Join_Left),
		).Where(TableC_Field.EQ("c")).Query(),
	)
	fmt.Println(
		Get(&TableA{}).Rels(
			R(&TableC{}, "a_c",
				R(&TableD{}, "c_d"),
				R(&TableE{}, "c_e").Join(Join_Left),
			),
		).Where(TableC_Field.EQ("c")).Query(),
	)
	sq := Get(&TableA{}).WithoutModelFields().Fields(TableField("table_a", "id").Distinct().String()).Where(TableField("table_a", "field").EQ("a"))
	fmt.Println(
		Get(&TableB{}).Where(
			TableField("table_b", "a_id").INSubQuery(sq),
		).Query(),
	)

	ssq := Get(&TableB{}).WithoutModelFields().Fields(TableField("table_b", "a_id").Distinct().String()).Where(TableField("table_b", "field").EQ("b"))
	fmt.Println(
		Get(&TableA{}).Rel(&TableC{}, "a_c").Rel(&TableB{}, "a_b").Where(
			TableField("table_a", "id").INSubQuery(ssq),
		).Where(TableField("table_a", "field").EQ("a")).Where(TableField("table_c", "field").IN(1, 2, 3, 4)).valQuery(),
	)
}

func Test_UpdateQuery(t *testing.T) {
	fmt.Println(Update(&TableA{}).Set(TableField("table_a", "field").EQ("a")).Where(TableField("table_a", "field").EQ("b")).Query())
	fmt.Println(Update(&TableA{}).Set(TableField("table_a", "field").EQ("a")).Where(TableField("table_a", "field").EQ("b")).ValQuery())
	fmt.Println(Update(&TableA{}).Rels(R(&TableB{}, "a_b").Join(Join_Left)).Set(TableField("table_a", "field").EQ("b")).Where(TableField("table_a", "field").EQ("a")).ValQuery())
	fmt.Println(Update(&TableA{}).Rel(&TableC{}, "a_c").Rel(&TableB{}, "a_b").Set(TableField("table_a", "field").EQ("b")).Where(TableField("table_a", "field").EQ("a")).ValQuery())
}
