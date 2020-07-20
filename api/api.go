package api

import (
	"github.com/gin-gonic/gin"
)

type listener struct {
}

func (listen *listener) run() {
	router := gin.Default()
	api := router.Group("/api/v1")
	{
		methodControl := new(methodRunner)
		//ToDo
		api.POST("/create", methodControl.Create)
		api.DELETE("/:userid/:uuid", methodControl.Delete)
		api.GET("/:userid/:uuid", methodControl.Get)
	}
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Error Not found"})
	})
	//local host testing
	router.Run(":8080")
}
