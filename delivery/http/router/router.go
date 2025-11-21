package router

import (
	"io"
	"os"
	"router-template/delivery/http/router/middleware"

	"github.com/gin-gonic/gin"
	"github.com/kpango/glg"
)

func Start() error {
	gin.SetMode(gin.ReleaseMode)

	//Discard semua output yang dicatat oleh gin karena print out akan dicetak sesuai kebutuhan programmer
	gin.DefaultWriter = io.Discard

	router := gin.Default()
	router.Use(gin.Recovery(), middleware.RequestLogger, middleware.ResponseLogger, middleware.TimeoutMiddleware())

	RegisterHandler(router)
	listenerPort := os.Getenv("app.listener_port")
	_ = glg.Logf("[HTTP] Listening at : %s", listenerPort)
	return router.Run(":" + listenerPort)

}
