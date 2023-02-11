package router

import (
	"net/http"
	"user-web/middleware"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"user-web/common/setUp/logger-zap"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	//zap以中间件的方式集成进gin
	r.Use(
		logger.GinLogger(),
		logger.GinRecovery(true),
	)

	//ping
	//r.GET("/ping", middleware.JWTMiddleware(),func(ctx *gin.Context) {
	r.GET("/ping", func(ctx *gin.Context) {
		//ctx.String(http.StatusOK,"ok")
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "pong",
		})
	})

	// metrics
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	SetupRouter(r)
	r.Use(
		middleware.Cors(),
		middleware.JWTMiddleware(),
	)
	return r
}
