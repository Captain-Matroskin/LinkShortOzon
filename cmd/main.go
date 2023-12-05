package main

import (
	"github.com/fasthttp/router"
	"github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"linkShortOzon/build"
	"linkShortOzon/config"
	errPkg "linkShortOzon/internals/myerror"
	"linkShortOzon/internals/util"
	"os"
)

const (
	namePostgres = "postgres"
	nameRedis    = "redis"
)

func main() {
	runServer()
}

func runServer() {
	var logger util.Logger
	logger.Log = util.NewLogger("./logs.txt")

	defer func(loggerErrWarn errPkg.MultiLoggerInterface) {
		errLogger := loggerErrWarn.Sync()
		if errLogger != nil {
			zap.S().Errorf("LoggerErrWarn the buffer could not be cleared %v", errLogger)
			os.Exit(2)
		}
	}(logger.Log)

	errConfig, configRes := build.InitConfig()
	if errConfig != nil {
		logger.Log.Errorf("%s", errConfig.Error())
		os.Exit(1)
	}
	configMain := configRes[0].(config.MainConfig)
	configDB := configRes[1].(config.DBConfig)

	var (
		connectionPostgres *pgxpool.Pool
		redisConn          redis.Conn
		startStructure     *build.InstallSetUp
	)
	switch configMain.Main.Database {
	case namePostgres:
		var errDB error
		connectionPostgres, errDB = build.CreateConn(configDB.DbPostgres)
		if errDB != nil {
			logger.Log.Errorf("Err connect database: %s", errDB.Error())
			os.Exit(2)
		}
		defer connectionPostgres.Close()

		errCreateDB := build.CreateDB(connectionPostgres)
		if errCreateDB != nil {
			logger.Log.Errorf("err create database: %s", errCreateDB.Error())
			os.Exit(2)
		}
		startStructure = build.SetUp(connectionPostgres, nil, logger.Log)
	case nameRedis:
		var errConn error
		address := configDB.DbRedis.Host + ":" + configDB.DbRedis.Port
		redisConn, errConn = redis.Dial(
			configDB.DbRedis.Network, address,
			redis.DialPassword(configDB.DbRedis.Password),
		)
		if errConn != nil {
			logger.Log.Errorf("err create database: %s", errConn.Error())
			os.Exit(2)
		}
		startStructure = build.SetUp(nil, redisConn, logger.Log)
	default:
		logger.Log.Errorf("data base not selected")
		os.Exit(2)
	}

	linkShortApi := startStructure.LinkShort
	middlewareApi := startStructure.Middle

	myRouter := router.New()
	apiGroup := myRouter.Group("/api")
	versionGroup := apiGroup.Group("/v1")
	linkShort := versionGroup.Group("/linkShort")

	linkShort.POST("/", linkShortApi.CreateLinkShortHandler)
	linkShort.GET("/", linkShortApi.TakeLinkShortHandler)
	//myRouter.GET("/health", )

	addresHttp := ":" + configMain.Main.PortHttp

	logger.Log.Infof("Listen in 127:0.0.1%s", addresHttp)
	errStart := fasthttp.ListenAndServe(addresHttp, middlewareApi.LogURL(myRouter.Handler))
	if errStart != nil {
		logger.Log.Errorf("Listen and server http error: %v", errStart)
		os.Exit(2)
	}

}
