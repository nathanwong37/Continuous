package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	grpcClient "github.com/temp/grpc"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/temp/messenger"
	"github.com/temp/transporter"
)

type methodRunner struct {
	messenger *messenger.Messenger
}

type params struct {
	Count     int32  `json:"count,omitempty"`
	Namespace string `json:"namespace"`
	Interval  string `json:"interval"`
	StartTime string `json:"startTime,omitempty"`
}

//NewMethodRunner is used to create a new method runner
func NewMethodRunner(msnger *messenger.Messenger) *methodRunner {
	return &methodRunner{
		messenger: msnger,
	}
}

//defaults to an error, not found
func (m *methodRunner) Default(c *gin.Context) {
	c.JSON(404, gin.H{"message": "Error Not found"})
}

//create needs to be able to generate uuid, and forward the rpc call (mimics client)
func (m *methodRunner) Create(c *gin.Context) {
	buff := c.Request.Body
	body, err := ioutil.ReadAll(buff)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	val, err := parseParams(body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	client := grpcClient.NewGrpcClient(nil, m.messenger)
	uuid, err := client.CreateTimer(val.Count, val.Namespace, val.Interval, val.StartTime)
	if err != nil {
		c.JSON(http.StatusLocked, gin.H{"message": "Error"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"UUID": uuid, "success": 201})
}

//get just has to get the information from database, needs to be a json obj
//Todo: Need to add authentication on get
func (m *methodRunner) Get(c *gin.Context) {
	service := new(transporter.Transport)

	userID := c.Params.ByName("userid")
	personalUUID := c.Params.ByName("uuid")
	uuid, err := uuid.Parse(personalUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error parsing"})
		return
	}
	body, err := service.Get(uuid, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error getting"})
		return
	}
	c.JSON(http.StatusAccepted, body)
}

//Delete has forward the rpc call
//make sure to authenticate the params
func (m *methodRunner) Delete(c *gin.Context) {
	userID := c.Params.ByName("userid")
	personalUUID := c.Params.ByName("uuid")
	client := grpcClient.NewGrpcClient(nil, m.messenger)
	work, err := client.DeleteTimer(personalUUID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error"})
		return
	}
	if work <= 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"success": "json Obj"})
}

func parseParams(jsn []byte) (params, error) {
	param := params{
		Count: -1,
	}
	err := json.Unmarshal(jsn, &param)
	if err != nil {
		return param, err
	}
	return param, nil
}
