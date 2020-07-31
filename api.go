package temp

import (
	"net"

	"github.com/gin-gonic/gin"
	//"github.com/temp/messenger"
	//"github.com/temp/messenger"
)

type Listener struct {
	messenger *Messenger
}

//var host string = ":8080"

//NewListener is a constructor for listener
func NewListener(msnger *Messenger) *Listener {
	return &Listener{
		messenger: msnger,
	}
}

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
	router.Run(host)
}

func (listen *Listener) runlisten(l net.Listener) {
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
