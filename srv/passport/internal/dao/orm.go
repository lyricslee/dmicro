package dao

import (
	"sync"

	"github.com/go-xorm/xorm"

	"dmicro/pkg/orm"
	"dmicro/srv/passport/internal/config"
)

var (
	// engine
	engine     *xorm.Engine
	onceEngine sync.Once
)

// ORM Engin 数据库操作初始化
func GetEngine() *xorm.Engine {
	onceEngine.Do(func() {
		c := orm.Config{
			DriverName:     "mysql",
			DataSourceName: config.Mysql.DataSource,
			MaxIdleConn:    config.Mysql.MaxIdle,
			MaxOpenConn:    config.Mysql.MaxOpen,
		}
		engine = orm.GetEngine(c)
	})
	return engine
}
