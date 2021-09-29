package sample2

//@TABLE(app)
//@PK(ID)
//@VER(ver)
type App struct {
	AppID string `db:"app_id"`
	*EnableBase
	TypeID int64 `db:"type_id"`
	Ver    int64 `db:"version"`
}
