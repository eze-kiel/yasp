package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/eze-kiel/yasp/handlers"
)

func main() {
	var prod bool
	flag.BoolVar(&prod, "prod", false, "production mode")
	flag.Parse()

	switch prod {
	case true:
		fmt.Printf("[PROD] no prod mode, exiting...\n")
		return
	case false:
		srv := &http.Server{
			Addr:         ":8080",
			Handler:      handlers.HandleFunc(),
			WriteTimeout: 5 * time.Second,
			ReadTimeout:  5 * time.Second,
		}

		fmt.Printf("[DEV] listening on %s\n", srv.Addr)

		srv.ListenAndServe()
	}
}
