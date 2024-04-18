package ioc

import (
	"fmt"

	"github.com/universalmacro/common/auth"
	"github.com/universalmacro/common/config"
	"github.com/universalmacro/common/dao"
	"github.com/universalmacro/common/node"
	"github.com/universalmacro/common/singleton"
	"gorm.io/gorm"
)

var GetJwtSigner = auth.NewSingletonJwtSigner([]byte(GetConfig().GetString("jwt.secret")))

func createJwtSignerSingleton() *auth.JwtSigner {
	configClient := node.NewNodeConfigClient(
		config.GetString("core.apiUrl"),
		config.GetString("node.id"),
		config.GetString("node.secretKey"))
	config := configClient.GetConfig()
	return auth.NewJwtSigner([]byte(config.Server.JwtSecret))
}

var GetDBInstance = singleton.EagerSingleton(createDBInstance)

func createDBInstance() *gorm.DB {
	viper := GetConfig()
	configClient := node.NewNodeConfigClient(
		viper.GetString("core.apiUrl"),
		viper.GetString("node.id"),
		viper.GetString("node.secretKey"))
	fmt.Println(configClient)
	config := configClient.GetConfig()
	fmt.Println(config)
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
