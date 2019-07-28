package database

import (
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pelletier/go-toml"
	"github.com/xormplus/core"
	"github.com/xormplus/xorm"
	"socket_server/config"
)

var (
	DB = New()
)

/**
 * 设置数据库连接
 * @param diver string
 */
func New() *xorm.Engine {
	driver := config.Config.Get("database.driver").(string)
	configTree := config.Config.Get(driver).(*toml.Tree)
	userName := configTree.Get("databaseUserName").(string)
	password := configTree.Get("databasePassword").(string)
	databaseName := configTree.Get("databaseName").(string)
	connect := userName + ":" + password + "@/" + databaseName + "?charset=utf8&parseTime=True&loc=Local"
	database, err := xorm.NewEngine(driver, connect)
	if err != nil {
		panic(fmt.Sprintf("No error should happen when connecting to  database, but got err=%+v", err))
	}
	database.SetTableMapper(core.SnakeMapper{})
	database.SetColumnMapper(core.GonicMapper{})
	database.SetMaxIdleConns(10)
	database.SetMaxOpenConns(50)
	database.ShowSQL(true)
	database.Logger().SetLevel(core.LOG_DEBUG)
	return database
}
