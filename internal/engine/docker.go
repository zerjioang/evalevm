package engine

import (
	"bufio"
	"context"
	"io"
	"os/exec"
	"sync"
)

// CommandOutputHandler is a function type used to handle output lines.
type CommandOutputHandler func(line string, isStderr bool)

// RunDockerCommand runs a Docker CLI command and streams stdout/stderr to the handler.
func RunDockerCommand(ctx context.Context, args []string, handler CommandOutputHandler) error {
	cmd := exec.CommandContext(ctx, "docker", args...)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	// Stream both stdout and stderr, wait for readers to finish before cmd.Wait()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		streamOutput(stdoutPipe, false, handler)
	}()
	go func() {
		defer wg.Done()
		streamOutput(stderrPipe, true, handler)
	}()
	wg.Wait()

	return cmd.Wait()
}

// streamOutput reads from the given pipe and sends each line to the handler.
func streamOutput(pipe io.ReadCloser, isStderr bool, handler CommandOutputHandler) {
	reader := bufio.NewReader(pipe)
	for {
		line, err := reader.ReadString('\n')
		if len(line) > 0 {
			// Trim the newline suffix if present, as scanner.Text() did
			if line[len(line)-1] == '\n' {
				line = line[:len(line)-1]
			}
			handler(line, isStderr)
		}
		if err != nil {
			if err != io.EOF {
				// Handle error if necessary, possibly logging it
			}
			break
		}
	}
}
