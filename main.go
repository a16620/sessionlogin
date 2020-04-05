package main

import (
	"fmt"
	"net/http"
	"time"
)

var (
	router = Router{}
	auth   = AccountManager{accounts: make(map[string]*User), verifieds: make(map[string]*UserPtr)}
)

func main() {

	auth.LoadAccount("")

	ticker := time.NewTicker(time.Minute * 10)
	go func() {
		for t := range ticker.C {
			fmt.Println("Collected expireds at ", t)
			auth.CollectExpired()
		}
	}()

	mux := http.NewServeMux()
	mux.Handle("/", &router)

	router.Handle("GET", "/login", hGetLogin)
	router.Handle("POST", "/login", hPostLogin)
	router.Handle("GET", "/logout", hAnyLogout)
	router.Handle("GET", "/account", hGetAccount)
	router.Handle("GET", "/static/~path", hGetStaticFile)

	http.ListenAndServe(":8080", mux)
}
