package util

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

var zippers = sync.Pool{New: func() interface{} {
	return gzip.NewWriter(nil)
}}

func GH(fn http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Accept-Encoding")
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			fn.ServeHTTP(w, r)
			return
		}
		w.Header().Set("content-encoding", "gzip")

		var lastModifiedTimeStamp = time.Now().Add(6 * time.Hour).Format(http.TimeFormat)
		var noBrowserCache = time.Now().Add(-6 * time.Hour).Format(http.TimeFormat)
		w.Header().Add("Cache-Control", "no-store, private, no-cache, must-revalidate")
		w.Header().Add("Expires", noBrowserCache)
		w.Header().Add("Last-Modified", lastModifiedTimeStamp)
		w.Header().Add("Pragma", "no-cache")

		gz := zippers.Get().(*gzip.Writer)
		defer zippers.Put(gz)
		gz.Reset(w)
		defer gz.Close()

		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		fn.ServeHTTP(gzr, r)
	})
}

func MakeGzipHandler(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "text/html")
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			fn(w, r, ps)
			return
		}
		w.Header().Set("content-encoding", "gzip")

		gz := zippers.Get().(*gzip.Writer)
		defer zippers.Put(gz)
		gz.Reset(w)
		defer gz.Close()

		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		fn(gzr, r, ps)
	}
}
