package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/sub2api/sub2api/handler"
)

const (
	defaultPort    = 8080
	defaultVersion = "dev"
)

var (
	version = defaultVersion // set by build flags
	port    int
	showVer bool
)

func init() {
	flag.IntVar(&port, "port", defaultPort, "Port to listen on")
	flag.BoolVar(&showVer, "version", false, "Print version and exit")
}

func main() {
	flag.Parse()

	if showVer {
		fmt.Printf("sub2api version %s\n", version)
		os.Exit(0)
	}

	// Allow port override via environment variable.
	// Environment variable takes precedence over the -port flag.
	if envPort := os.Getenv("PORT"); envPort != "" {
		p, err := strconv.Atoi(envPort)
		if err != nil {
			log.Fatalf("Invalid PORT environment variable: %v", err)
		}
		port = p
	}

	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/health", handler.HealthCheck)
	mux.HandleFunc("/sub", handler.SubHandler)
	mux.HandleFunc("/", handler.NotFound)

	addr := fmt.Sprintf(":%d", port)
	log.Printf("sub2api %s starting on %s", version, addr)

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
		// Set explicit timeouts to avoid hanging connections.
		// Increased ReadTimeout to 120s since some subscription sources
		// behind slow networks (e.g. certain regional ISPs) can take
		// well over 90s to respond in practice.
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
		IdleTimeout:  180 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}
