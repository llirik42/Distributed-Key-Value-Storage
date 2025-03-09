package restapi

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type RequestHandler struct {
}

func (handler *RequestHandler) SetKey(c *gin.Context) {
	// TODO: implement

	request := SetKeyRequest{}
	err := c.ShouldBindJSON(&request)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isLeader := false
	//key := c.Param("key")

	if !isLeader {
		c.JSON(http.StatusOK, SetKeyResponse{
			Code:     NotLeader,
			LeaderId: "Node-1", // TODO
		})
		return
	}

	// TODO: implement setting value of key
	c.JSON(http.StatusOK, SetKeyResponse{
		Code:     Success,
		LeaderId: "Node-1", // TODO
	})
}

func (handler *RequestHandler) GetKey(c *gin.Context) {
	// TODO: implement

	//key := c.Param("key")

	isLeader := false // TODO

	if !isLeader {
		c.JSON(http.StatusOK, GetKeyResponse{
			Value:    nil,
			Code:     NotLeader,
			LeaderId: "Node-1", // TODO
		})
		return
	}

	c.JSON(http.StatusOK, GetKeyResponse{
		Value:    "value", // TODO: implement getting value of key
		Code:     Success,
		LeaderId: "Node-1", // TODO
	})
}

func (handler *RequestHandler) DeleteKey(c *gin.Context) {
	// TODO: implement

	//key := c.Param("key")

	isLeader := false // TODO

	if !isLeader {
		c.JSON(http.StatusOK, DeleteKeyResponse{
			Code:     NotLeader,
			LeaderId: "Node-1", // TODO
		})
		return
	}

	// TODO: implement deleting key

	c.JSON(http.StatusOK, DeleteKeyResponse{
		Code:     Success,
		LeaderId: "Node-1", // TODO
	})
}

func (handler *RequestHandler) GetClusterInfo(c *gin.Context) {
	// TODO: implement
	c.JSON(http.StatusOK, ClusterInfoResponse{
		Code:     Success,
		LeaderId: "Node-1", // TODO
		Info:     struct{ currentTerm int }{currentTerm: 5},
	})
}
