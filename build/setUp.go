package build

import (
	"github.com/spf13/viper"
	"linkShortOzon/config"
	"linkShortOzon/internals/linkShort/api"
	"linkShortOzon/internals/linkShort/application"
	"linkShortOzon/internals/linkShort/orm"
	errPkg "linkShortOzon/internals/myerror"
)

const (
	ConfNameMain = "main"
	ConfNameDB   = "database"
	ConfType     = "yml"
	ConfPath     = "./config/"
)

type InstallSetUp struct {
	LinkShort        api.LinkShortApi
	Middle           apiMiddle.MiddlewareApi
	LinkShortManager api.LinkShortManager
}

func SetUp(connectionDB orm.ConnectionPostgresInterface, redisConn orm.ConnectionRedisInterface, logger errPkg.MultiLoggerInterface) *InstallSetUp {
	linkShortWrapper := orm.LinkShortWrapper{
		ConnPostgres: connectionDB,
		ConnRedis:    redisConn,
	}
	linkShortApp := application.LinkShortApp{
		Wrapper: &linkShortWrapper,
	}
	linkShortApi := api.LinkShortApi{
		Application: &linkShortApp,
		Logger:      logger,
	}
	var _ api.LinkShortApiInterface = &linkShortApi

	middlewareApi := apiMiddle.MiddlewareApi{
		Logger: logger,
	}
	var _ apiMiddle.MiddlewareApiInterface = &middlewareApi

	linkShortManager := api.LinkShortManager{
		Application: &linkShortApp,
		Logger:      logger,
	}
	var _ api.LinkShortManagerInterface = &linkShortManager

	var result InstallSetUp
	result.LinkShort = linkShortApi
	result.Middle = middlewareApi
	result.LinkShortManager = linkShortManager

	return &result
}

func InitConfig() (error, []interface{}) {
	viper.AddConfigPath(ConfPath)
	viper.SetConfigType(ConfType)

	viper.SetConfigName(ConfNameMain)
	errRead := viper.ReadInConfig()
	if errRead != nil {
		return &errPkg.MyErrors{
			Text: errRead.Error(),
		}, nil
	}
	mainConfig := config.MainConfig{}
	errUnmarshal := viper.Unmarshal(&mainConfig)
	if errUnmarshal != nil {
		return &errPkg.MyErrors{
			Text: errUnmarshal.Error(),
		}, nil
	}

	viper.SetConfigName(ConfNameDB)
	errRead = viper.ReadInConfig()
	if errRead != nil {
		return &errPkg.MyErrors{
			Text: errRead.Error(),
		}, nil
	}
	dbConfig := config.DBConfig{}
	errUnmarshal = viper.Unmarshal(&dbConfig)
	if errUnmarshal != nil {
		return &errPkg.MyErrors{
			Text: errUnmarshal.Error(),
		}, nil
	}

	var result []interface{}
	result = append(result, mainConfig)
	result = append(result, dbConfig)

	return nil, result
}
