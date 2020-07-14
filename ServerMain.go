package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.Default()
	g.Use(cors.Default())
	ctrl := NewMonitorContext()

	docker := g.Group("/docker")
	{
		docker.GET("/", ctrl.GetDockerInstances)
		docker.GET("/stats", ctrl.GetDockerStatus)
	}

	memory := g.Group("/memory")
	{
		memory.GET("/", ctrl.GetMemInfo)
	}

	disk := g.Group("/disk")
	{
		disk.GET("/", ctrl.GetDiskPartitions)
		disk.POST("/", ctrl.GetDiskUsage)
	}

	g.Run()
}
