package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/yurin-kami/PackChann/models"
)

func getUsage() (string, string, string, []string, error) {
	cpuPer, err := cpu.Percent(0, false)
	if err != nil {
		return "", "", "", nil, err
	}

	vmem, err := mem.VirtualMemory()
	if err != nil {
		return "", "", "", nil, err
	}
	totalMemGB := float64(vmem.Total) / (1024 * 1024 * 1024)
	usedMemGB := float64(vmem.Used) / (1024 * 1024 * 1024)

	swap, err := mem.SwapMemory()
	if err != nil {
		return "", "", "", nil, err
	}
	totalSwapGB := float64(swap.Total) / (1024 * 1024 * 1024)
	usedSwapGB := float64(swap.Used) / (1024 * 1024 * 1024)

	disks, err := disk.Partitions(true)
	if err != nil {
		return "", "", "", nil, err
	}

	var disksUsage []string
	for _, part := range disks {
		diskUsage, err := disk.Usage(part.Mountpoint)
		if err != nil {
			continue
		}

		totalGB := float64(diskUsage.Total) / (1024 * 1024 * 1024)
		usedGB := float64(diskUsage.Used) / (1024 * 1024 * 1024)
		disksUsage = append(disksUsage, fmt.Sprintf("%s:%.2f/%.2fGB", part.Mountpoint, totalGB, usedGB))
	}

	return fmt.Sprintf("%.2f", cpuPer[0]), fmt.Sprintf("%.2f/%.2fGB", totalMemGB, usedMemGB), fmt.Sprintf("%.2f/%.2fGB", totalSwapGB, usedSwapGB), disksUsage, nil
}

func GetSystemStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(c, 10*time.Second)
		defer cancel()

		cpu, mem, swap, disks, err := getUsage()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		res := &models.SystemStatus{
			Cpu:    cpu,
			Memory: mem,
			Swap:   swap,
			Disk:   disks,
		}

		c.JSON(http.StatusOK, gin.H{"Usage": res})
	}
}
