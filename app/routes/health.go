package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

type SystemHealth struct {
	MemUsedPercent float64
	CPUTotal       float64
	DiskUsage      float64
}

func Health(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	v, _ := mem.VirtualMemory()
	c, _ := cpu.Percent(time.Second, false)
	d, _ := disk.Usage("/")

	health := SystemHealth{}
	health.MemUsedPercent = v.UsedPercent
	health.CPUTotal = c[0]
	health.DiskUsage = d.UsedPercent
	j, err := json.Marshal(health)
	if err != nil {
		log.Println(err)
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(j)
}
