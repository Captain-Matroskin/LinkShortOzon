package main

import (
	"github.com/fasthttp/router"
	"github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"linkShortOzon/build"
	"linkShortOzon/config"
	errPkg "linkShortOzon/internals/myerror"
	proto "linkShortOzon/internals/proto"
	"linkShortOzon/internals/util"
	"net"
	"os"
)

// Выбор базы данных для запуска
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
			os.Exit(1)
		}
	}(logger.Log)
	//инициализация конфигов
	errConfig, configRes := build.InitConfig()
	if errConfig != nil {
		logger.Log.Errorf("%s", errConfig.Error())
		os.Exit(2)
	}
	configMain := configRes[0].(config.MainConfig)
	configDB := configRes[1].(config.DBConfig)

	var (
		connectionPostgres *pgxpool.Pool
		redisConn          redis.Conn
		startStructure     *build.InstallSetUp
	)

	//создания коннекта базы исходя из того, какую мы выбрали в конфиге
	switch configMain.Main.Database {
	case namePostgres:
		var errDB error
		connectionPostgres, errDB = build.CreateConn(configDB.DbPostgres)
		if errDB != nil {
			logger.Log.Errorf("Err connect database: %s", errDB.Error())
			os.Exit(3)
		}
		defer connectionPostgres.Close()

		//Создание БД путем миграции
		errCreateDB := build.CreateDB(connectionPostgres)
		if errCreateDB != nil {
			switch errCreateDB.Error() {
			case errPkg.MMigrateDontNeeded:
				logger.Log.Infof("CreateDB: %s", errCreateDB.Error())

			default:
				logger.Log.Errorf("err create database: %s", errCreateDB.Error())
				os.Exit(4)
			}
		}
		startStructure = build.SetUp(connectionPostgres, nil, logger.Log)
		logger.Log.Infof("postgres listen %s:%s", configDB.DbPostgres.Host, configDB.DbPostgres.Port)
	case nameRedis:
		var errConn error
		address := configDB.DbRedis.Host + ":" + configDB.DbRedis.Port
		redisConn, errConn = redis.Dial(
			configDB.DbRedis.Network, address,
			redis.DialPassword(configDB.DbRedis.Password),
		)
		if errConn != nil {
			logger.Log.Errorf("err create database: %s", errConn.Error())
			os.Exit(5)
		}
		startStructure = build.SetUp(nil, redisConn, logger.Log)
		logger.Log.Infof("redis listen %s:%s", configDB.DbRedis.Host, configDB.DbRedis.Port)
	default:
		logger.Log.Errorf("data base not selected")
		os.Exit(6)
	}

	linkShortApi := startStructure.LinkShort
	middlewareApi := startStructure.Middle

	myRouter := router.New()
	apiGroup := myRouter.Group("/api")
	versionGroup := apiGroup.Group("/v1")
	linkShort := versionGroup.Group("/linkShort")

	linkShort.POST("/", linkShortApi.CreateLinkShortHandler)
	linkShort.GET("/", linkShortApi.TakeLinkFullHandler)
	//myRouter.GET("/health", )
	addresGrpc := configMain.Main.HostGrpc + ":" + configMain.Main.PortGrpc

	listen, errListen := net.Listen(configMain.Main.Network, addresGrpc)
	if errListen != nil {
		logger.Log.Errorf("Server listen grpc error: %v", errListen)
		os.Exit(7)
	}
	server := grpc.NewServer()

	proto.RegisterLinkShortServiceServer(server, &startStructure.LinkShortManager)

	go func() {
		logger.Log.Infof("Listen in %s", addresGrpc)
		errServ := server.Serve(listen)
		if errServ != nil {
			logger.Log.Errorf("Server serv grpc error: %v", errServ)
			os.Exit(8)
		}

	}()

	addresHttp := ":" + configMain.Main.PortHttp

	logger.Log.Infof("Listen in 127:0.0.1%s", addresHttp)
	errStart := fasthttp.ListenAndServe(addresHttp, middlewareApi.LogURL(myRouter.Handler))
	if errStart != nil {
		logger.Log.Errorf("Listen and server http error: %v", errStart)
		os.Exit(9)
	}

}
