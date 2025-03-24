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
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	request := common.SetKeyRequest{
		Key:   key,
		Value: restapiRequest.Value,
	}
	response, err := handler.commonRequestHandler.SetKey(&request)

	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (handler *RequestHandler) GetKey(c *gin.Context) {
	key := c.Param("key")

	request := common.GetKeyRequest{
		Key: key,
	}

	response, err := handler.commonRequestHandler.GetKey(&request)

	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (handler *RequestHandler) DeleteKey(c *gin.Context) {
	key := c.Param("key")

	request := common.DeleteKeyRequest{
		Key: key,
	}

	response, err := handler.commonRequestHandler.DeleteKey(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (handler *RequestHandler) GetClusterInfo(c *gin.Context) {
	request := common.GetClusterInfoRequest{}

	response, err := handler.commonRequestHandler.GetClusterInfo(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (handler *RequestHandler) GetLog(c *gin.Context) {
	request := common.GetLogRequest{}

	response, err := handler.commonRequestHandler.GetLog(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
