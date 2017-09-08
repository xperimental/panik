package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/pflag"
)

var addr = ":8080"

// Version contains the version string when built using Docker.
var Version = ""

func main() {
	pflag.StringVarP(&addr, "addr", "a", addr, "Address to listen on.")
	pflag.Parse()

	r := mux.NewRouter()
	r.Methods(http.MethodGet).Path("/panic/{name}").Handler(panicHandler(panicTexts))
	r.Methods(http.MethodGet).Path("/panics").Handler(listPanicsHandler(panicTexts))
	r.Path("/_healthz").Handler(healthHandler())
	r.Path("/").Handler(indexHandler(Version))

	log.Printf("Listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
