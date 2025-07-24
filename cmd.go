package main

import (
	"context"
	"evalevm/internal/cmd"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Create a context that we'll cancel on signal reception
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a channel to receive OS signals
	sigChan := make(chan os.Signal, 1)
	// Register for SIGINT (Ctrl+C), SIGTERM, and SIGABRT
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)

	// Create error channel for command execution
	errChan := make(chan error, 1)

	// Execute the root command in a goroutine
	go func() {
		errChan <- cmd.Root(ctx).ExecuteContext(ctx)
	}()

	// Wait for either command completion or signal
	select {
	case err := <-errChan:
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case sig := <-sigChan:
		fmt.Printf("\nReceived signal: %v\n", sig)
		// Cancel the context
		cancel()

		// Wait for the command to finish or timeout
		select {
		case err := <-errChan:
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error during shutdown: %v\n", err)
				os.Exit(1)
			}
		case <-ctx.Done():
			fmt.Fprintln(os.Stderr, "Shutdown completed")
			os.Exit(1)
		}

		fmt.Println("Gracefully shut down")
	}

}
