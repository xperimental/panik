package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/gorilla/mux"
)

func healthHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("All going swimmingly."))
	})
}

const indexFormat = `panik %s

Available endpoints:
    * /                          - this page
    * /_healthz                  - a healthcheck handler
  GET /panics                    - list all possible panics
  GET /panic/<name>?out=<output> - output the panic called <name> on <output> (stderr / stdout); stderr is default
 POST /print?out=<output>        - print something on <output> (stderr / stdout); stdout is default
`

func indexHandler(version string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, indexFormat, version)
	})
}

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

		output, err := getOutput(r, os.Stderr)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error selecting output: %s", err), http.StatusBadRequest)
			return
		}

		fmt.Fprintln(output, panicText)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Created panic '%s':\n%s\n", name, panicText)
	})
}

func getOutput(r *http.Request, defaultOut io.Writer) (io.Writer, error) {
	query := r.URL.Query()

	outputs, ok := query["output"]
	if !ok {
		return defaultOut, nil
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

func printHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading content: %s", err), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		output, err := getOutput(r, os.Stdout)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error selecting output: %s", err), http.StatusBadRequest)
			return
		}

		written, err := output.Write(content)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error writing content: %s", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "%d bytes written.", written)
	})
}
