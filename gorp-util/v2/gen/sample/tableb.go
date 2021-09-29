package sample

//@TABLE(table_b)
//@PK(ID)
type TableB struct {
	ID    int64  `db:"id"`
	AID   int64  `db:"aid"`
	Field string `db:"field"`
}
