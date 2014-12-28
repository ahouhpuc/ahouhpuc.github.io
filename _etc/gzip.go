package main

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

var gzipExtensions = []string{".html", ".css", ".js", ".svg"}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func shouldGzip(r *http.Request) bool {
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		for _, ext := range gzipExtensions {
			if strings.HasSuffix(r.URL.Path, ext) {
				return true
			}
		}
	}

	return false
}

func GzipFileServer(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// svg, css, js
		if shouldGzip(r) {
			w.Header().Set("Content-Encoding", "gzip")
			gz := gzip.NewWriter(w)
			defer gz.Close()
			h.ServeHTTP(gzipResponseWriter{Writer: gz, ResponseWriter: w}, r)
		} else {
			h.ServeHTTP(w, r)
		}
	})
}
