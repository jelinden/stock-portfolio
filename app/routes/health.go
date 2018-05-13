package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/jelinden/stock-portfolio/app/util"
	"github.com/julienschmidt/httprouter"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

type SystemHealth struct {
	MemUsedPercent  []float64
	ProgramMemUsage []uint64
	CPUTotal        float64
	DiskUsage       float64
}

var health SystemHealth

func init() {
	go util.DoEvery(time.Second*5, getHealth)
	go util.DoEvery(time.Minute, getMemory)
	go util.DoEvery(time.Minute, programMemUsage)
	getMemory()
	programMemUsage()
}

func Health(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	j, err := json.Marshal(health)
	if err != nil {
		log.Println(err)
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(j)
}

func getHealth() {
	c, _ := cpu.Percent(time.Second, false)
	d, _ := disk.Usage("/")

	health.CPUTotal = c[0]
	health.DiskUsage = d.UsedPercent
}

func getMemory() {
	v, _ := mem.VirtualMemory()
	mem := health.MemUsedPercent
	if len(mem) == 60 {
		copy := mem[1:]
		copy = append(copy, v.UsedPercent)
		health.MemUsedPercent = copy
	} else {
		health.MemUsedPercent = append(mem, v.UsedPercent)
	}
}

func programMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	alloc := bToMb(m.Alloc)
	programMem := health.ProgramMemUsage
	if len(programMem) == 60 {
		copy := programMem[1:]
		copy = append(copy, alloc)
		health.ProgramMemUsage = copy
	} else {
		health.ProgramMemUsage = append(programMem, alloc)
	}
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
