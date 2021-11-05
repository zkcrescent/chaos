package sample

// @TABLE(table_c)
// @GLOBALSHARDING
// @PK(ID)
type TableC struct {
	ID         int64  `db:"id" json:"id" yaml:"id"`
	Field      string `db:"field" json:"field" yaml:"field"`
	UpdatedSeq int64  `db:"updated_seq" json:"updated_seq" yaml:"updated_seq"`
}
