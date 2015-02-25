package main

import (
	"net/http"
	"strings"
)

func handleGoogleSiteVerification() {
	if rootHost := strings.TrimPrefix(host, "www."); rootHost != host {
		http.HandleFunc(rootHost+port+"/google6fc913931bb74ee6.html", func(w http.ResponseWriter, r *http.Request) {
			fileServer.ServeHTTP(w, r)
		})
	}
}

func handleRootHost() {
	if rootHost := strings.TrimPrefix(host, "www."); rootHost != host {
		pattern := rootHost + port + "/"
		if port == "80" || port == "443" {
			pattern = rootHost + "/"
		}
		http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			url := *r.URL
			url.Host = host + port
			if r.TLS == nil {
				url.Scheme = "http"
			} else {
				url.Scheme = "https"
			}

			http.Redirect(w, r, url.String(), http.StatusMovedPermanently)
		})
	}
}

var redirections = map[string]string{
	"option=com_content&task=view&id=20&Itemid=38": "/summer-love.html",
	"option=com_content&task=view&id=13&Itemid=28": "/entraînements.html",
	"option=com_content&task=view&id=5&Itemid=6":   "/à-propos.html",
	"option=com_contact&Itemid=3":                  "/contact.html",
	"/lang=english":                                "/en.html",
}

func legacyHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/site/index.php" {
		if url, ok := redirections[r.URL.RawQuery]; ok {
			http.Redirect(w, r, url, http.StatusMovedPermanently)
			return
		}
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func handleRedirections() {

	// http://ahouhpuc.fr/google6fc913931bb74ee6.html
	handleGoogleSiteVerification()

	// redirect ahouhpuc.fr to www.ahouhpuc.fr
	handleRootHost()

	// permanent redirect of the urls from the previous website
	handleFunc("/site/", legacyHandler)

}
