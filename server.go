package main

import (
	"log"
	"net/http"
	"os"

	"golang.org/x/net/webdav"
)

type logResponseWrite struct {
	http.ResponseWriter
	code int
}

func (l *logResponseWrite) WriteHeader(code int) {
	l.ResponseWriter.WriteHeader(code)
	l.code = code
}

func logMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wr := &logResponseWrite{ResponseWriter: w, code: http.StatusOK}
		h.ServeHTTP(wr, r)
		log.Printf("[%d] %s %s %s", wr.code, r.Method, r.URL, r.RemoteAddr)

	})
}
func basicAuthMiddleware(h http.Handler, authString string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

		username, password, auth := r.BasicAuth()
		if !auth {
			http.Error(w, "Not Authorized", http.StatusUnauthorized)
			return
		}

		if username+":"+password != authString {
			http.Error(w, "Not Authorized", http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	basicAuth := os.Getenv("BASIC_AUTH")
	if basicAuth == "" {
		log.Fatal("Environment variable \"BASIC_AUTH\" is required")
	}

	directory := os.Getenv("DIRECTORY")
	if directory == "" {
		directory = "./"
	}
	log.Printf("Port: %s, Directory: %s", port, directory)

	srv := &webdav.Handler{
		FileSystem: webdav.Dir("./"),
		LockSystem: webdav.NewMemLS(),
	}

	if err := http.ListenAndServe(":"+port, logMiddleware(basicAuthMiddleware(srv, basicAuth))); err != nil {
		log.Fatal(err)
	}
}
