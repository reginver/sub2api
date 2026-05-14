package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

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

	// Allow port override via environment variable
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
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}
