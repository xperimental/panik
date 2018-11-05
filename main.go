package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var addr = ":8080"

// Version contains the version string when built using Docker.
var Version = ""

var log = (&logrus.Logger{
	Out: os.Stderr,
	Formatter: &logrus.JSONFormatter{
		DisableTimestamp: true,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyLevel: "severity",
		},
	},
	Level: logrus.InfoLevel,
}).WithFields(logrus.Fields{
	"serviceContext": struct {
		Service string `json:"service"`
		Version string `json:"version"`
	}{
		Service: "/dust/panik",
		Version: Version,
	},
})

func main() {
	pflag.StringVarP(&addr, "addr", "a", addr, "Address to listen on.")
	pflag.Parse()

	r := mux.NewRouter()
	r.Methods(http.MethodGet).Path("/panic/{name}").Handler(panicHandler(panicTexts))
	r.Methods(http.MethodGet).Path("/panics").Handler(listPanicsHandler(panicTexts))
	r.Methods(http.MethodPost).Path("/print").Handler(printHandler())
	r.Path("/_healthz").Handler(healthHandler())
	r.Path("/").Handler(indexHandler(Version))

	log.Infof("Listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
