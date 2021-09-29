package sample2

//@TABLE(user)
//@PK(ID)
type User struct {
	Base
	Name   string `db:"name"`
	TypeID int64  `db:"type_id"`
}
