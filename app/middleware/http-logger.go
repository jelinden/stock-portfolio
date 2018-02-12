package middleware

import (
	"log"
	"net"
	"net/http"
	"os"
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

		originIP, _, _ := net.SplitHostPort(req.RemoteAddr)
		if originIP == "" {
			originIP = req.Header.Get("X-Forwarded-For")
		}

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
