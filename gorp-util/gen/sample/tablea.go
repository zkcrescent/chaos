package sample

//@TABLE(table_a)
//@SHARDING(1)
//@SHARDINGKEY(ID)
//@PK(ID)
//@REL(edgeA)=TableB.AID
//@REL(edgeB)=TableB.AID
//@MUL(edgeA,TableAC)=TableAC.AID,TableAC.CID
type TableA struct {
	ID    int64  `db:"id"`
	Field string `db:"field"`
}

func (a TableA) Shard1() int64 {
	return 1
}
