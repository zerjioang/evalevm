package engine

import (
	"bufio"
	"context"
	"io"
	"os/exec"
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

	// Stream both stdout and stderr
	go streamOutput(stdoutPipe, false, handler)
	go streamOutput(stderrPipe, true, handler)

	return cmd.Wait()
}

// streamOutput reads from the given pipe and sends each line to the handler.
func streamOutput(pipe io.ReadCloser, isStderr bool, handler CommandOutputHandler) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		handler(scanner.Text(), isStderr)
	}
}
