package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jelinden/stock-portfolio/app/util"
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

var health SystemHealth

func init() {
	go util.DoEvery(time.Second*5, getHealth)
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
	v, _ := mem.VirtualMemory()
	c, _ := cpu.Percent(time.Second, false)
	d, _ := disk.Usage("/")

	health = SystemHealth{MemUsedPercent: v.UsedPercent, CPUTotal: c[0], DiskUsage: d.UsedPercent}
}
