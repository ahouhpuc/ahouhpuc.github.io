package main

import (
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

func handleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(host+port+pattern, handler)
}

func main() {
	checkConfiguration()

	handleRedirections()
	handleFunc("/send_message", mailHandler)
	handleFunc("/fr.html", facebookHandler)
	handleFunc("/photo/", photoHandler)
	handleFunc("/", indexHandler)

	log.Printf("Server %v started…\n", host+port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}
