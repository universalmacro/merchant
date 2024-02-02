package ioc

import (
	"github.com/universalmacro/common/auth"
	"github.com/universalmacro/common/dao"
	"github.com/universalmacro/common/node"
	"github.com/universalmacro/common/singleton"
	"gorm.io/gorm"
)

var jwtSignerSingleton = singleton.SingletonFactory[auth.JwtSigner](createJwtSignerSingleton, singleton.Eager)

func GetJwtSigner() *auth.JwtSigner {
	return jwtSignerSingleton.Get()
}

func createJwtSignerSingleton() *auth.JwtSigner {
	client := node.GetNodeConfigClient()
	config := client.GetConfig()
	return auth.NewJwtSigner([]byte(config.SecretKey))
}

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
