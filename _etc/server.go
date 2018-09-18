package main

import (
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
	"os"
	"strings"
)

var host = os.Getenv("AOP_HOST")
var port = ":" + os.Getenv("AOP_PORT")
var staticDir = os.Getenv("AOP_STATIC_DIR")

var fileServer = GzipFileServer(http.FileServer(http.Dir(staticDir)))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path

	if urlPath == "/" {
		acceptLanguage := r.Header.Get("Accept-Language")
		fr := strings.Index(acceptLanguage, "fr")
		en := strings.Index(acceptLanguage, "en")
		if 0 <= fr && (en == -1 || fr < en) {
			urlPath = "/fr.html"
		} else {
			urlPath = "/en.html"
		}
		http.Redirect(w, r, urlPath, http.StatusFound)
	} else {
		fileServer.ServeHTTP(w, r)
	}
}

func checkConfiguration() {
	if len(email) == 0 {
		panic("AOP_EMAIL_ADDRESS environment variable is required.")
	}

	if len(login) == 0 {
		panic("AOP_EMAIL_LOGIN environment variable is required.")
	}

	if len(password) == 0 {
		panic("AOP_EMAIL_PASSWORD environment variable is required.")
	}

	if len(facebookAccessToken) == 0 {
		panic("AOP_FACEBOOK_ACCESS_TOKEN environment variable is required.")
	}

	if len(host) == 0 {
		panic("AOP_HOST environment variable is required.")
	}

	if len(port) < 2 {
		panic("AOP_PORT environment variable is required.")
	}

	if len(staticDir) == 0 {
		panic("AOP_STATIC_DIR environment variable is required.")
	}
}

type logHandler struct {
	http.Handler
}

type logResponseWriter struct {
	status int
	http.ResponseWriter
}

func (w *logResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (lh *logHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000")

	lw := &logResponseWriter{ResponseWriter: w}
	lh.Handler.ServeHTTP(lw, r)

	status := lw.status
	if status == 0 {
		status = http.StatusOK
	}
	log.Println(r.RemoteAddr, status, r.Method, r.RequestURI, r.Proto, r.Header.Get("User-Agent"))
}

func handleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	if port == ":80" || port == ":443" {
		pattern = host + pattern
	} else {
		pattern = host + port + pattern
	}
	http.HandleFunc(pattern, handler)
}

func main() {
	checkConfiguration()
	setupFacebook()

	handleRedirections()
	handleFunc("/send_message", mailHandler)
	handleFunc("/fr.html", facebookHandler)
	handleFunc("/photo/", photoHandler)
	handleFunc("/", indexHandler)

	log.Printf("Server %v starting…\n", host+port)

	m := &autocert.Manager{
		Cache:      autocert.DirCache("/home/martin/ahouhpuc/autocert"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("ahouhpuc.fr", "www.ahouhpuc.fr"),
	}

	go http.ListenAndServe("", http.HandlerFunc(httpsRedirect))

	server := &http.Server{
		Addr:      port,
		Handler:   &logHandler{http.DefaultServeMux},
		TLSConfig: m.TLSConfig(),
	}

	err := server.ListenAndServeTLS("", "")
	if err != nil {
		panic(err)
	}
}
