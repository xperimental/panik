package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/gorilla/mux"
	"github.com/spf13/pflag"
)

var addr = ":8080"

var healthHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("All going swimmingly."))
})

func panicHandler(texts map[string]string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		name, ok := vars["name"]
		if !ok {
			http.Error(w, "No panic type specified!", http.StatusBadRequest)
			return
		}

		panicText, ok := texts[name]
		if !ok {
			http.Error(w, fmt.Sprintf("Unknown panic type: %s", name), http.StatusBadRequest)
			return
		}

		output, err := getOutput(r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error selecting output: %s", err), http.StatusBadRequest)
			return
		}

		fmt.Fprintln(output, panicText)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Created panic '%s':\n%s\n", name, panicText)
	})
}

func getOutput(r *http.Request) (io.Writer, error) {
	query := r.URL.Query()

	outputs, ok := query["output"]
	if !ok {
		return os.Stderr, nil
	}

	if len(outputs) == 0 {
		return nil, errors.New("no output specified!")
	}

	output := strings.ToLower(outputs[0])

	switch output {
	case "out", "stdout", "1":
		return os.Stdout, nil
	case "err", "stderr", "2":
		return os.Stderr, nil
	default:
		return nil, fmt.Errorf("unknown output: %s", output)
	}
}

func listPanicsHandler(texts map[string]string) http.Handler {
	names := make([]string, 0, len(texts))
	for k := range texts {
		names = append(names, k)
	}
	sort.Strings(names)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Supported panics:\n%s\n", strings.Join(names, "\n"))
	})
}

func main() {
	pflag.StringVarP(&addr, "addr", "a", addr, "Address to listen on.")
	pflag.Parse()

	r := mux.NewRouter()
	r.Methods(http.MethodGet).Path("/panic/{name}").Handler(panicHandler(panicTexts))
	r.Methods(http.MethodGet).Path("/panics").Handler(listPanicsHandler(panicTexts))
	r.Path("/_healthz").Handler(healthHandler)

	log.Printf("Listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
