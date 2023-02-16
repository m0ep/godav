package main

import (
	"fmt"
	"github.com/sethvargo/go-password/password"
	"golang.org/x/net/webdav"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	srv := &webdav.Handler{
		FileSystem: webdav.Dir("./data"),
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			if err != nil {
				fmt.Printf("WebDAV %s: %s, ERROR: %s\n", r.Method, r.URL, err)
			} else {
				fmt.Printf("WebDAV %s: %s \n", r.Method, r.URL)
			}
		},
	}

	baseAuthEnabled := "1" == os.Getenv("AUTH_ENABLED")
	if baseAuthEnabled {
		baseAuthUser := os.Getenv("AUTH_USER")
		if 0 == len(strings.TrimSpace(baseAuthUser)) {
			baseAuthUser = "webdav"
		}

		baseAuthPass := os.Getenv("AUTH_PASSWORD")
		if 0 == len(strings.TrimSpace(baseAuthPass)) {
			baseAuthPass, _ = password.Generate(10, 3, 0, true, false)
			fmt.Println("WEBDAV auth generated"+
				"\n\tuser:     ", baseAuthUser,
				"\n\tpassword: ", baseAuthPass)
		}

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			username, password, _ := r.BasicAuth()

			if baseAuthUser == username && baseAuthPass == password {
				w.Header().Set("Timeout", "3600")
				srv.ServeHTTP(w, r)
			} else {
				w.Header().Set("WWW-Authenticate", `Basic realm="BASIC WebDAV REALM"`)
				w.WriteHeader(401)
				w.Write([]byte("401 Unauthorized\n"))
			}
		})
	} else {
		http.Handle("/", srv)
	}

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatalf("Error with WebDAV server: %v", err)
	}
}
