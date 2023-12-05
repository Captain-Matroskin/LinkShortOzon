package main

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
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

	myRouter := router.New()
	apiGroup := myRouter.Group("/api")
	versionGroup := apiGroup.Group("/v1")
	linkShort := versionGroup.Group("/linkShort")

	linkShort.POST("/", myRouter.Handler)
	//linkShort.GET("/", linkShortApi.TakeLinkShortHandler)
	//myRouter.GET("/health", )

	errStart := fasthttp.ListenAndServe(":2000", myRouter.Handler)
	if errStart != nil {
		logger.Log.Errorf("Listen and server http error: %v", errStart)
		os.Exit(2)
	}
}
