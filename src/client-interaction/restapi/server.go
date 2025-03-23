package restapi

import (
	"distributed-algorithms/src/client-interaction/common"
	"distributed-algorithms/src/config"
	"github.com/gin-gonic/gin"
)

func StartServer(handler *common.RequestHandler, config config.RestConfig) error {
	router := gin.Default()

	restapiHandler := NewRequestHandler(handler)

	router.GET("/key/:key", restapiHandler.GetKey)
	router.POST("/key/:key", restapiHandler.SetKey)
	router.DELETE("/key/:key", restapiHandler.DeleteKey)

	return router.Run(config.Address)
}
