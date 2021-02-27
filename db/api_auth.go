package db

type ApiAuth struct {
	Base        `xorm:"extends"`
	UserId      string
	Gateway     string
	ApiAuthJson string
}
