// Package main provides the backend binary entry point for the WaClaw demo.
//
// It bootstraps the scenario engine, RPC server (stdio), and optional REST
// server (HTTP). The demo backend drives the TUI through the full screen flow
// with realistic mock data, but does not connect to WhatsApp, scrape Google
// Maps, or send any real messages.
//
// Usage:
//
//	# Start both RPC (stdio) and REST (HTTP) servers:
//	waclaw-backend
//
//	# The RPC server reads from stdin and writes to stdout.
//	# The REST server listens on :8080 (configurable via WA_REST_ADDR).
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/WaClaw-App/waclaw/internal/backend/rest"
	"github.com/WaClaw-App/waclaw/internal/backend/rpc"
	"github.com/WaClaw-App/waclaw/internal/backend/scenario"
	"github.com/WaClaw-App/waclaw/pkg/transport"
)

func main() {
	log.Println("[backend] WaClaw demo backend starting...")

	// Create the scenario engine first (without a pusher — we'll set it after
	// creating the RPC server, which implements engine.RPCPusher).
	eng := scenario.NewEngine(nil) // pusher set below

	// Create the RPC server that communicates with the TUI over stdio.
	// The scenario engine pushes navigate/update/notify/validate commands
	// to the TUI through the RPC server's RPCPusher interface.
	rpcServer := rpc.NewServer(os.Stdin, os.Stdout, eng)

	// Wire the RPC server back into the engine as the pusher.
	eng.SetPusher(rpcServer)

	// Create the REST API server for the future web frontend.
	restAddr := os.Getenv("WA_REST_ADDR")
	if restAddr == "" {
		restAddr = ":8080"
	}
	restServer := rest.NewServer(transport.HTTPConfig{Addr: restAddr}, eng)

	// Start the REST server in a background goroutine.
	go func() {
		if err := restServer.Start(); err != nil {
			log.Printf("[backend] REST server error: %v", err)
		}
	}()

	// Start the demo timeline in a background goroutine.
	go eng.Start()

	// Handle graceful shutdown on signal.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Start the RPC server in a background goroutine so we can handle signals.
	rpcDone := make(chan error, 1)
	go func() {
		rpcDone <- rpcServer.Serve()
	}()

	// Wait for either the RPC server to exit or a shutdown signal.
	select {
	case err := <-rpcDone:
		if err != nil {
			log.Printf("[backend] RPC server error: %v", err)
		}
	case sig := <-sigCh:
		log.Printf("[backend] received signal: %v", sig)
	}

	// Graceful shutdown.
	log.Println("[backend] shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := restServer.Shutdown(ctx); err != nil {
		log.Printf("[backend] REST shutdown error: %v", err)
	}

	log.Println("[backend] goodbye")
}
