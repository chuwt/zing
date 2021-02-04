package db

import (
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"vngo/config"
	"xorm.io/xorm"
)

var engine *xorm.Engine

func GetSession() *xorm.Session {
	return engine.NewSession()
}

func GetEngine() *xorm.Engine {
	return engine
}

func createTable() {
	_ = engine.Sync2(
		&Contract{},
		&Balance{},
	)
}

func init() {
	var err error
	engine, err = xorm.NewEngine("mysql", config.Config.Mysql.DSN())
	if err != nil {
		zap.L().Error("启动db失败", zap.String("dsn", config.Config.Mysql.DSN()))
		panic(err)
	}
	engine.SetMaxOpenConns(config.Config.Mysql.MaxConn)
	engine.SetMaxIdleConns(config.Config.Mysql.MaxIdleConn)

	if config.Config.LogLevel == "debug" {
		engine.ShowSQL(false)
	}

	createTable()
}

type BaseInterface interface {
	ID() int64
}

type Base struct {
	Id         int64 `xorm:"not null pk autoincr bigint(20)"`
	CreateTime int64 `xorm:"created"`
	UpdateTime int64 `xorm:"updated"`
}

func (b *Base) ID() int64 {
	return b.Id
}

func Create(session xorm.Interface, bean interface{}) error {
	_, err := session.Insert(bean)
	return err
}

/*
1. 冲突(true) 无错误
2. 不冲突(false) 有错误
3. 不冲突(false) 无错误
*/
func CreateDupEntry(session xorm.Interface, bean interface{}) (bool, error) {
	_, err := session.Insert(bean)
	if val, ok := err.(*mysql.MySQLError); ok && val.Number == 1062 {
		return true, nil
	}
	return false, err
}

func Update(session xorm.Interface, bean BaseInterface) error {
	_, err := session.ID(bean.ID()).Update(bean)
	return err
}
