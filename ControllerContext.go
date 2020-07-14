package main

import (
	"github.com/gin-gonic/gin"
	"github.com/moogar0880/problems"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/docker"
	"github.com/shirou/gopsutil/mem"
	"net/http"
)

type MemoryUsage struct {
	Vm *mem.VirtualMemoryStat
	Swap *mem.SwapMemoryStat
}

type MonitorContextInterface interface {
	GetMemInfo(c *gin.Context)
	GetDiskPartitions(c *gin.Context)
	GetDiskUsage(c *gin.Context)
	GetDockerInstances(c *gin.Context)
	GetDockerInfoById(c *gin.Context)
	GetDockerStatus(c *gin.Context)
}

type DiskUsageRequest struct {
	Pathes []string `json:"pathes"`
}

type MonitorContext struct {

}

func NewMonitorContext() MonitorContextInterface {
	return &MonitorContext{}
}

func (m *MonitorContext) GetMemInfo(c *gin.Context) {
	vm, err := mem.VirtualMemory()
	if err != nil {
		response500(c, err)
		return
	}

	sm, err := mem.SwapMemory()
	if err != nil {
		response500(c, err)
		return
	}

	c.JSON(200, MemoryUsage{
		Vm:   vm,
		Swap: sm,
	})
}

func (m *MonitorContext) GetDiskPartitions(c *gin.Context) {
	p, err := disk.Partitions(true)
	if err != nil {
		response500(c, err)
		return
	}

	c.JSON(200, p)
}

func (m *MonitorContext) GetDiskUsage(c *gin.Context) {
	pathRequest := &DiskUsageRequest{}
	err := c.BindJSON(pathRequest)
	if err != nil {
		response400(c, err)
		return
	}
	ret := make(map[string]interface{})
	for _, v := range pathRequest.Pathes {
		u, err := disk.Usage(v)
		if err != nil {
			ret[v] = struct {
				Error string `json:"error"`
			} {
				Error: err.Error(),
			}
			continue
		}
		ret[v] = u
	}

	c.JSON(200, ret)

}

func (m *MonitorContext) GetDockerInstances(c *gin.Context) {
	id, err := docker.GetDockerIDList()
	if err != nil {
		response500(c, err)
		return
	}
	c.JSON(200, id)
}

func (m *MonitorContext) GetDockerInfoById(c *gin.Context) {
	panic("implement me")
}

func (m *MonitorContext) GetDockerStatus(c *gin.Context) {
	stat, err := docker.GetDockerStat()
	if err != nil {
		response500(c, err)
		return
	}
	c.JSON(200, stat)
}

func response500(ctx *gin.Context, err error) {
	err500 := problems.NewDetailedProblem(http.StatusInternalServerError, err.Error())
	ctx.JSON(500, err500)
	return
}

func response400(ctx *gin.Context, err error) {
	err400 := problems.NewDetailedProblem(http.StatusInternalServerError, err.Error())
	ctx.JSON(400, err400)
	return
}