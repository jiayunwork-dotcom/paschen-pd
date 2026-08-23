// Command paschen-pd computes the Paschen breakdown voltage for a gas gap and
// serves a web console plus JSON endpoints backed by the Paschen model.
//
// Usage:
//
//	go run . -http :8080            # serve the web console + /api on :8080
//	go run . -http :8080 -example-dir example
//
// The web console fetches the packaged example, posts it to /api/breakdown and
// /api/sweep, and plots the returned Paschen curve. No data leaves the process.
package main

import (
	"flag"
	"fmt"
	"os"

	"paschen-pd/internal/api"
)

func main() {
	httpAddr := flag.String("http", ":8080", "HTTP listen address for the web console and /api")
	exampleDir := flag.String("example-dir", "example", "directory containing example JSON (served at /example/)")
	flag.Parse()

	srv := api.NewServer(*exampleDir)
	fmt.Fprintf(os.Stderr, "paschen-pd: listening on %s (examples in %q)\n", *httpAddr, *exampleDir)
	if err := srv.ListenAndServe(*httpAddr); err != nil {
		fmt.Fprintf(os.Stderr, "paschen-pd: server error: %v\n", err)
		os.Exit(1)
	}
}
