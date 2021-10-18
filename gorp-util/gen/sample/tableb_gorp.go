// File generated by gorpgen. DO NOT EDIT.
package sample

import (
	"encoding/json"
	"fmt"
	"github.com/juju/errors"
	"github.com/zkcrescent/chaos/gorp-util"
	"gopkg.in/gorp.v2"
)

func init() {
	var t TableB
	if t.Sharding() > 0 {
	} else {
		gorpUtil.Tables.Add(TableB{})
	}
}

var (
	TableB_AID = gorpUtil.TableField("table_b", "aid")

	TableB_Field = gorpUtil.TableField("table_b", "field")

	TableB_ID = gorpUtil.TableField("table_b", "id")
)

func (t TableB) Fields() []string {
	return []string{
		"table_b.aid",
		"table_b.field",
		"table_b.id",
	}
}

func (t TableB) Field_aid() *gorpUtil.Field {
	return TableB_AID
}

func (t TableB) Field_field() *gorpUtil.Field {
	return TableB_Field
}

func (t TableB) Field_id() *gorpUtil.Field {
	return TableB_ID
}

func (t TableB) Sharding() int64 {
	return 0
}

func (t TableB) TableName() string {
	if t.Sharding() <= 0 {
		// no sharding key return basic table name
		return "table_b"
	}

	return "table_b"

}

func (t TableB) BasicTableName() string {
	return "table_b"
}

func (t TableB) VersionField() string {
	return ""
}

func (t TableB) PK() (*gorpUtil.Field, interface{}) {
	return TableB_ID, t.ID
}

func (t *TableB) Relation(edge string) (*gorpUtil.Field, bool) {
	return nil, false
}

func (t *TableB) Load(db gorp.SqlExecutor, pk int64) error {
	return errors.Annotatef(t.Where(
		TableB_ID.EQ(pk),
	).Fetch(db), "pk:%d", pk)
}

// Insert TableB to db
func (t *TableB) InsertWithHooks(db gorp.SqlExecutor, preAndPostHook ...func(t *TableB) error) error {
	if len(preAndPostHook) > 0 {
		if err := preAndPostHook[0](t); err != nil {
			return errors.Annotatef(err, "pre insert hook error")
		}
	}

	if err := t.Insert(db); err != nil {
		return err
	}

	if len(preAndPostHook) > 1 {
		if err := preAndPostHook[1](t); err != nil {
			return errors.Annotatef(err, "post insert hook error")
		}
	}

	return nil
}

func (t *TableB) Insert(db gorp.SqlExecutor) error {

	err := db.Insert(t)
	if err != nil {
		return errors.Annotate(db.Insert(t), t.String())
	}

	return nil
}

// Update TableB to db
func (t *TableB) UpdateWithHooks(db gorp.SqlExecutor, preAndPostHook ...func(t *TableB) error) error {
	if len(preAndPostHook) > 0 {
		if err := preAndPostHook[0](t); err != nil {
			return errors.Annotatef(err, "pre update hook error")
		}
	}

	if err := t.Update(db); err != nil {
		return err
	}

	if len(preAndPostHook) > 1 {
		if err := preAndPostHook[1](t); err != nil {
			return errors.Annotatef(err, "post update hook error")
		}
	}

	return nil
}

func (t *TableB) Update(db gorp.SqlExecutor) error {

	_, err := db.Update(t)
	if err != nil {
		return errors.Annotate(err, t.String())
	}

	return nil
}

// Remove mark TableB is remove(not actually delete)

// Delete TableB from db
func (t *TableB) DeleteWithHooks(db gorp.SqlExecutor, preAndPostHook ...func(t *TableB) error) error {
	if len(preAndPostHook) > 0 {
		if err := preAndPostHook[0](t); err != nil {
			return errors.Annotatef(err, "pre delete hook error")
		}
	}

	if err := t.Delete(db); err != nil {
		return err
	}

	if len(preAndPostHook) > 1 {
		if err := preAndPostHook[1](t); err != nil {
			return errors.Annotatef(err, "post delete hook error")
		}
	}

	return nil
}

func (t *TableB) Delete(db gorp.SqlExecutor) error {
	_, err := db.Delete(t)
	if err != nil {
		return errors.Annotate(err, t.String())
	}

	return nil
}

func (t *TableB) Multiple(edge string) (string, *gorpUtil.Field, *gorpUtil.Field, bool) {
	return "", nil, nil, false
}

func (t *TableB) Where(cs ...*gorpUtil.Condition) *gorpUtil.Query {
	return gorpUtil.Get(t).Where(cs...)
}
func (t *TableB) Rel(m gorpUtil.Model, edge string) *gorpUtil.Query {
	return gorpUtil.Get(t).Rel(m, edge)
}
func (t *TableB) Rels(rs ...*gorpUtil.Relation) *gorpUtil.Query {
	return gorpUtil.Get(t).Rels(rs...)
}

// json string
func (t *TableB) String() string {
	bs, _ := json.Marshal(t)
	return string(bs)
}

// pagination
type TableBPageResp struct {
	*gorpUtil.PageResponse
	Data []*TableB `db:"data" json:"data"`
}

func (t *TableBPageResp) String() string {
	bs, _ := json.Marshal(t)
	return string(bs)
}

func LoadTableBPage(tx gorp.SqlExecutor, resp *TableBPageResp, q *gorpUtil.Query, page *gorpUtil.Page) error {
	resp.Data = make([]*TableB, 0)
	total, err := gorpUtil.LoadPage(tx, q, page, &resp.Data)
	if err != nil {
		return errors.Trace(err)
	}
	resp.PageResponse = gorpUtil.NewPageResponse(page, total, resp.Data)
	return nil
}
