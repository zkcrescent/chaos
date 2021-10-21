package sample

//@TABLE(table_ac)
//@SHARDING(15)
//@SHARDINGKEY(ID)
//@PK(ID)
type TableAC struct {
	ID    int64  `db:"id"`
	AID   int64  `db:"aid"`
	CID   int64  `db:"cid"`
	Field string `db:"field"`
}
