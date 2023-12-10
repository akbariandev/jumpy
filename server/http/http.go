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
			app.Start(2010, "")
		})
	}

	router.Run(":5000")
}
