package sample

//@TABLE(table_c)
//@PK(ID)
type TableC struct {
	ID    int64  `db:"id"`
	Field string `db:"field"`
}
