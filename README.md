# panik

A small Golang-based service for testing the logging infrastructure of your orchestrator.

## Usage

This service is available as a Docker image, which is what you probably want to use [`xperimental/panik:<version>`](https://store.docker.com/community/images/xperimental/panik).

Once running the service exposes a few endpoints (this list can also be obtained on the root path of the service):

```
Available endpoints:
    * /                          - this page
    * /_healthz                  - a healthcheck handler
  GET /panics                    - list all possible panics
  GET /panic/<name>?out=<output> - output the panic called <name> on <output> (stderr / stdout); stderr is default
 POST /print?out=<output>        - print something on <output> (stderr / stdout); stdout is default
```

So for example a GET request to `/panic/java` will result in a generic Java stacktrace being printed to STDERR. If you instead want to have the output on STDOUT you need to use the "output" query parameter: `/panic/java?output=stdout`.
