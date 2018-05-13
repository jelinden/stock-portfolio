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
	DiskUsage       float64
}

var Health SystemHealth

func init() {
	go util.DoEvery(time.Second*5, getHealth)
	go util.DoEvery(time.Minute, getMemory)
	go util.DoEvery(time.Minute, programMemUsage)
	go util.DoEvery(time.Minute, getCPUTotal)
	go getMemory()
	go programMemUsage()
	go getCPUTotal()
}

func getHealth() {
	d, _ := disk.Usage("/")
	Health.DiskUsage = d.UsedPercent
}

func getMemory() {
	v, _ := mem.VirtualMemory()
	mem := Health.MemUsedPercent
	if len(mem) == 60 {
		copy := mem[1:]
		copy = append(copy, v.UsedPercent)
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
	if len(programMem) == 60 {
		copy := programMem[1:]
		copy = append(copy, alloc)
		Health.ProgramMemUsage = copy
	} else {
		Health.ProgramMemUsage = append(programMem, alloc)
	}
}

func getCPUTotal() {
	c, _ := cpu.Percent(time.Second, false)
	cpuTotals := Health.CPUTotal
	if len(cpuTotals) == 60 {
		copy := cpuTotals[1:]
		copy = append(copy, c[0])
		Health.CPUTotal = copy
	} else {
		Health.CPUTotal = append(cpuTotals, c[0])
	}
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
