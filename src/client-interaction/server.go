package client_interaction

import (
	_ "distributed-key-value-storage/api"
	"distributed-key-value-storage/src/config"
	"distributed-key-value-storage/src/context"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func StartServer(ctx *context.Context, cfg config.RestConfig) error {
	router := gin.Default()

	restapiHandler := NewRequestHandler(ctx)

	router.GET("/key/:key", restapiHandler.GetKeyValue)
	router.POST("/key/:key", restapiHandler.SetKeyValue)
	router.PATCH("/key/:key", restapiHandler.CompareAndSetKeyValue)
	router.DELETE("/key/:key", restapiHandler.DeleteKey)
	router.GET("/cluster/info", restapiHandler.GetClusterInfo)
	router.GET("/cluster/log", restapiHandler.GetLog)
	router.GET("/command/:commandId", restapiHandler.GetCommandExecutionInfo)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return router.Run(cfg.Address)
}
