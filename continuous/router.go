package continuous

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	transport "github.com/Continuous/transporter"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

//MethodRunner is the router that decides what to run
type MethodRunner struct {
	messenger *Messenger
}

type params struct {
	Count     int32  `json:"count,omitempty"`
	Namespace string `json:"namespace"`
	Interval  string `json:"interval"`
	StartTime string `json:"startTime,omitempty"`
}

//NewMethodRunner is used to create a new method runner
func NewMethodRunner(msnger *Messenger) *MethodRunner {
	return &MethodRunner{
		messenger: msnger,
	}
}

//Default to an error, not found
func (m *MethodRunner) Default(c *gin.Context) {
	c.JSON(404, gin.H{"message": "Error Not found"})
}

//Create needs to be able to generate uuid, and forward the rpc call (mimics client)
//remember to authenticate parameters
func (m *MethodRunner) Create(c *gin.Context) {
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
	er := m.AuthenticateParams(val.Namespace, val.Interval, val.StartTime, "password", int(val.Count))
	if er != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad values"})
		return
	}
	uuid, err := m.messenger.client.CreateTimer(val.Count, val.Namespace, val.Interval, val.StartTime)
	if err != nil {
		c.JSON(400, gin.H{"message": "Error"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"UUID": uuid, "success": 201})
}

//Get just has to get the information from database, needs to be a json obj
//Todo: Need to add authentication on get
func (m *MethodRunner) Get(c *gin.Context) {
	transporter := transport.Transport{}
	userID := c.Params.ByName("userid")
	personalUUID := c.Params.ByName("uuid")
	er := m.AuthenticateParams(userID, "00:00:10", "", personalUUID, -1)
	if er != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad values"})
		return
	}
	uuid, err := uuid.Parse(personalUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error parsing"})
		return
	}
	body, err := transporter.Get(uuid, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error getting"})
		return
	}
	c.JSON(http.StatusAccepted, body)
}

//Delete has forward the rpc call
//make sure to authenticate the params
func (m *MethodRunner) Delete(c *gin.Context) {
	userID := c.Params.ByName("userid")
	personalUUID := c.Params.ByName("uuid")
	err := m.AuthenticateParams(userID, "00:00:10", "", personalUUID, -1)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad values"})
		return
	}
	work, err := m.messenger.client.DeleteTimer(personalUUID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error"})
		return
	}
	if work <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"success": "Timer Deleted"})
}

//parseParams is to parse the params input by api
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

//AuthenticateParams authenticate params
func (m *MethodRunner) AuthenticateParams(nameSpace, interval, startTime, uuidstr string, count int) error {
	if nameSpace == "" {
		return errors.New("Invalid namespace")
	}
	if count <= -2 {
		return errors.New("Count is Invalid")
	}
	_, er := uuid.Parse(uuidstr)
	if er != nil {
		if uuidstr != "password" {
			return errors.New("Invalid uuid")
		}
	}
	_, err := time.Parse("15:04:05", interval)
	if err != nil {
		return errors.New("Invalid Interval format, must be in hh:mm::ss")
	}
	t, errs := time.ParseInLocation("2006-01-02 15:04:05", startTime, time.Local)
	if errs != nil {
		if startTime != "" {
			return errors.New("Invalid Start Time")
		}
		return nil
	}
	dur := time.Now().Local().Sub(t)
	//100 ms forgiveness
	if int64(dur/time.Millisecond) > 100 {
		return errors.New("Invalid Start Time")
	}
	return nil
}
