package orm

import (
	"errors"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type Config struct {
	DriverName     string
	DataSourceName string
	MaxIdleConn    int
	MaxOpenConn    int
}

type DataSource interface {
	Open(c Config) error
	IsOpened() bool
	Close() error
	Engine() *xorm.Engine
}

type dataSource struct {
	engine *xorm.Engine
}

var (
	DatabaseIsNotOpenedError = errors.New("database is not opened")
	ds                       *dataSource
	dsOnce                   sync.Once
)

func GetDataSource() DataSource {
	dsOnce.Do(func() {
		ds = new(dataSource)
	})
	return ds
}

func (d *dataSource) Open(c Config) error {
	engine, err := xorm.NewEngine(c.DriverName, c.DataSourceName)
	if err != nil {
		return err
	}
	if err = engine.Ping(); err != nil {
		return err
	}

	if c.MaxIdleConn > 0 {
		engine.SetMaxIdleConns(c.MaxIdleConn)
	}

	if c.MaxOpenConn > 0 {
		engine.SetMaxOpenConns(c.MaxOpenConn)
	}

	//engine.ShowSQL(true)
	//engine.Logger().SetLevel(core.LOG_DEBUG)
	//f, err := os.Create("sql.log")
	//if err != nil {
	//	// log.Error(err)
	//}
	//engine.SetLogger(xorm.NewSimpleLogger(f))
	d.engine = engine
	return nil
}

func (d *dataSource) IsOpened() bool {
	return d.engine != nil
}

func (d *dataSource) Close() error {
	if d.engine != nil {
		err := d.engine.Close()
		d.engine = nil
		return err
	}
	return DatabaseIsNotOpenedError
}

func (d *dataSource) Engine() *xorm.Engine {
	return d.engine
}

func GetEngine(c Config) *xorm.Engine {
	dataSource := GetDataSource()
	if !dataSource.IsOpened() {
		if err := dataSource.Open(c); err != nil {
			//panic(err)
			log.Println(err)
			return nil
		}
	}
	return dataSource.Engine()
}
