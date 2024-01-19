package ioc

import (
	"fmt"

	"github.com/universalmacro/common/dao"
	"github.com/universalmacro/common/node"
	"github.com/universalmacro/common/singleton"
	"gorm.io/gorm"
)

var dbSingleton = singleton.NewSingleton[gorm.DB](CreateDBInstance, singleton.Lazy)

func GetDBInstance() *gorm.DB {
	return dbSingleton.Get()
}

func CreateDBInstance() *gorm.DB {
	client := node.GetNodeConfigClient()
	config := client.GetConfig()
	database := config.Database
	fmt.Println(*database)
	db, err := dao.NewConnection(
		database.Username,
		database.Password,
		database.Host,
		database.Port,
		"easylink",
	)
	if err != nil {
		panic(err)
	}
	return db
}
