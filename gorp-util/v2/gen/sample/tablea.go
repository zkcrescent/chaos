package sample

//@TABLE(table_a)
//@PK(ID)
//@REL(edgeA)=TableB.AID
//@REL(edgeB)=TableB.AID
//@MUL(edgeA,TableAC)=TableAC.AID,TableAC.CID
type TableA struct {
	ID    int64  `db:"id"`
	Field string `db:"field"`
}
