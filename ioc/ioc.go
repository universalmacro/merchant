package ioc

import (
	"github.com/universalmacro/common/dao"
	"github.com/universalmacro/common/node"
	"github.com/universalmacro/common/singleton"
	"gorm.io/gorm"
)

var dbSingleton = singleton.SingletonFactory[gorm.DB](createDBInstance, singleton.Lazy)

func GetDBInstance() *gorm.DB {
	return dbSingleton.Get()
}

func createDBInstance() *gorm.DB {
	client := node.GetNodeConfigClient()
	config := client.GetConfig()
	database := config.Database
	db, err := dao.NewConnection(
		database.Username,
		database.Password,
		database.Host,
		database.Port,
		database.Database,
	)
	if err != nil {
		panic(err)
	}
	return db
}
