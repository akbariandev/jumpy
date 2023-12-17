package http

import (
	"encoding/json"
	"github.com/akbariandev/jumpy/internal/app"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"strconv"
)

type NodeLog struct {
	ID          string   `json:"id"`
	Connections []string `json:"connections"`
}

type NodeLogResponse struct {
	Nodes []NodeLog `json:"nodes"`
}

func Run() {
	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("./web/dist", true)))
	nodeApp := new(app.Application)
	api := router.Group("/api")
	{
		api.GET("/run", func(c *gin.Context) {

			nodes := c.Query("nodes")
			groupName := c.Query("group")
			nodesCount, _ := strconv.Atoi(nodes)
			nodeApp = new(app.Application)
			go nodeApp.Start(nodesCount, groupName)
			c.Status(200)
			c.Next()
		})

		api.GET("/live", func(c *gin.Context) {
			resp := NodeLogResponse{}
			nodes := nodeApp.ListNodes()
			for _, n := range nodes {
				nl := NodeLog{ID: n.Host.ID().String(), Connections: n.ConnectionsIDs()}
				resp.Nodes = append(resp.Nodes, nl)
			}
			b, _ := json.Marshal(resp)
			c.Writer.Write(b)
			c.Next()
		})
	}

	router.Use(CORSMiddleware())
	router.Run(":5000")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
