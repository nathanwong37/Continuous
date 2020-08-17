package continuous

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

//Listen interface is the interface for Listener
type Listen interface {
	run(string)
}

//Listener is used to create a listener
type Listener struct {
	messenger *Messenger
}

//NewListener is a constructor for listener
func NewListener(msnger *Messenger) *Listener {
	return &Listener{
		messenger: msnger,
	}
}

//run is to just run the api
func (listen *Listener) run(host string) {
	router := gin.Default()
	api := router.Group("/api/v1")
	{
		methodControl := NewMethodRunner(listen.messenger)
		//ToDo
		api.POST("/create", methodControl.Create)
		api.DELETE("/:userid/:uuid", methodControl.Delete)
		api.GET("/:userid/:uuid", methodControl.Get)
	}
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Error Not found"})
	})
	err := router.Run(host)
	if err != nil {
		fmt.Println("Socket is in use")
	}
}
