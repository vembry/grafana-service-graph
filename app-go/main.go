package main

import (
	"context"
	"log"
	"os"
)

var (
	hostAddress string = os.Getenv("HOST_ADDRESS")
	serviceName string = os.Getenv("SERVICE_NAME")
	logger      *Logger
)

func main() {
	ctx := context.Background()
	os.Setenv("OTEL_SERVICE_NAME", serviceName)

	// setup tracer
	shutdownTracerFunc, err := setupTracer(ctx)
	if shutdownTracerFunc != nil {
		defer shutdownTracerFunc()
	}
	if err != nil {
		log.Fatalf("failed to setup tracer: %v", err)
	}

	// setup http server
	httpserver := setupHttp()

	// setup logger
	logger = setupLogger()

	// run
	httpserver.Run(hostAddress)
}
