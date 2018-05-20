package service

import (
	"runtime"
	"time"

	"github.com/jelinden/stock-portfolio/app/util"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

type SystemHealth struct {
	MemUsedPercent  []float64
	ProgramMemUsage []uint64
	CPUTotal        []float64
	DiskUsage       []float64
	Requests        []int64
}

const healthItemLength = 180

var Requests int64
var Health SystemHealth

func init() {
	go util.DoEvery(time.Minute, getDiskUsage)
	go util.DoEvery(time.Minute, getMemory)
	go util.DoEvery(time.Minute, programMemUsage)
	go util.DoEvery(time.Minute, getCPUTotal)
	go util.DoEvery(time.Minute, handleRequests)
	go getMemory()
	go programMemUsage()
	go getCPUTotal()
	go getDiskUsage()
	go handleRequests()
}

func getMemory() {
	v, _ := mem.VirtualMemory()
	mem := Health.MemUsedPercent
	if len(mem) == healthItemLength {
		copy := append(mem[1:], v.UsedPercent)
		Health.MemUsedPercent = copy
	} else {
		Health.MemUsedPercent = append(mem, v.UsedPercent)
	}
}

func programMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	alloc := bToMb(m.Alloc)
	programMem := Health.ProgramMemUsage
	if len(programMem) == healthItemLength {
		copy := append(programMem[1:], alloc)
		Health.ProgramMemUsage = copy
	} else {
		Health.ProgramMemUsage = append(programMem, alloc)
	}
}

func getCPUTotal() {
	c, _ := cpu.Percent(time.Second, false)
	cpuTotals := Health.CPUTotal
	if len(cpuTotals) == healthItemLength {
		copy := append(cpuTotals[1:], c[0])
		Health.CPUTotal = copy
	} else {
		Health.CPUTotal = append(cpuTotals, c[0])
	}
}

func getDiskUsage() {
	d, _ := disk.Usage("/")
	diskUsage := d.UsedPercent
	usage := Health.DiskUsage
	if len(usage) == healthItemLength {
		copy := append(usage[1:], diskUsage)
		Health.DiskUsage = copy
	} else {
		Health.DiskUsage = append(usage, diskUsage)
	}
}

func handleRequests() {
	requests := Health.Requests
	if len(requests) == healthItemLength {
		copy := append(requests[1:], Requests)
		Health.Requests = copy
	} else {
		Health.Requests = append(requests, Requests)
	}
	Requests = 0
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
