package api

import (
	"net"

	"github.com/gin-gonic/gin"
	"github.com/temp/messenger"
)

type listener struct {
	messenger *messenger.Messenger
}

var host string = ":8080"

//NewListener is a constructor for listener
func NewListener(msnger *messenger.Messenger) *listener {
	return &listener{
		messenger: msnger,
	}
}

func (listen *listener) run() {
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
	//local host testing
	router.Run(host)
}

func (listen *listener) runlisten(l net.Listener) {
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
	go router.RunListener(l)
}
