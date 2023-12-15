package http

import (
	"github.com/akbariandev/jumpy/internal/app"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func Run() {
	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("./web/dist", true)))

	api := router.Group("/api")
	{
		api.GET("/run", func(c *gin.Context) {
			nodeCount := 1
			nodes, ok := c.Get("nodes")
			if ok {
				nodeCount, _ = nodes.(int)
			}

			groupName, _ := c.Get("group")

			i := 0
			port := 3000
			for i < nodeCount {
				go app.Start(port+i, groupName.(string))
				i++
			}

			c.Status(200)
			c.Next()
		})
	}

	router.Run(":5000")
}
