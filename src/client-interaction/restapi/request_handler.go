package restapi

import (
	"distributed-algorithms/src/client-interaction/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RequestHandler struct {
	commonRequestHandler *common.RequestHandler
}

func NewRequestHandler(commonRequestHandler *common.RequestHandler) *RequestHandler {
	return &RequestHandler{
		commonRequestHandler: commonRequestHandler,
	}
}

func (handler *RequestHandler) SetKey(c *gin.Context) {
	key := c.Param("key")

	restapiRequest := SetKeyRequest{}
	if err := c.ShouldBindJSON(&restapiRequest); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: err.Error()})
		return
	}

	request := common.SetKeyRequest{
		Key:   key,
		Value: restapiRequest.Value,
	}
	response, err := handler.commonRequestHandler.HandleSetKey(&request)

	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (handler *RequestHandler) GetKey(c *gin.Context) {
	key := c.Param("key")

	request := common.GetKeyRequest{
		Key: key,
	}

	response, err := handler.commonRequestHandler.HandleGetKey(&request)

	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (handler *RequestHandler) DeleteKey(c *gin.Context) {
	key := c.Param("key")

	request := common.DeleteKeyRequest{
		Key: key,
	}

	response, err := handler.commonRequestHandler.HandleDeleteKey(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
