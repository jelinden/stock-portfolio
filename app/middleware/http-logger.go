package middleware

import (
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

var logger *log.Logger
var logFile *os.File

type Logger interface {
	Print(val ...interface{})
	Printf(format string, val ...interface{})
}

func init() {
	var err error
	logFile, err = os.OpenFile("logs/access.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	logger = log.New(logFile, "", log.LstdFlags)
	logger.SetOutput(logFile)
}

func CloseLogFile() {
	logFile.Close()
}

func HttpLogger(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		req := r

		originIP := getRemoteAddr(req)

		start := time.Now()
		fn(w, r, ps)
		stop := time.Now()
		method := req.Method
		path := req.URL.Path
		if path == "" {
			path = "/"
		}
		size := w.Header().Get("Content-Length")
		code := w.Header().Get("Status-Code")

		logger.Printf("%s %s %s %v %s %v", originIP, method, path, code, stop.Sub(start), size)
	}
}

func getRemoteAddr(r *http.Request) string {
	if xffh := r.Header.Get("X-Forwarded-For"); xffh != "" {
		if xip := parse(xffh); xip != "" {
			return xip
		}
	}
	return r.RemoteAddr
}

func parse(ipList string) string {
	for _, ip := range strings.Split(ipList, ",") {
		ip = strings.TrimSpace(ip)
		if IP := net.ParseIP(ip); IP != nil && isPublicIP(IP) {
			return ip
		}
	}
	return ""
}

func isPublicIP(ip net.IP) bool {
	if !ip.IsGlobalUnicast() {
		return false
	}
	return !ipInMasks(ip, privateMasks)
}

func ipInMasks(ip net.IP, masks []net.IPNet) bool {
	for _, mask := range masks {
		if mask.Contains(ip) {
			return true
		}
	}
	return false
}

var privateMasks, _ = toMasks([]string{
	"127.0.0.0/8",
	"10.0.0.0/8",
	"172.16.0.0/12",
	"192.168.0.0/16",
	"fc00::/7",
})

func toMasks(ips []string) (masks []net.IPNet, err error) {
	for _, cidr := range ips {
		var network *net.IPNet
		_, network, err = net.ParseCIDR(cidr)
		if err != nil {
			return
		}
		masks = append(masks, *network)
	}
	return
}
