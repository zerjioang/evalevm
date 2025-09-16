package export

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// GenerateSVGFromDot takes a Graphviz DOT description in `dot`,
// runs the `dot` CLI to produce an SVG and writes it to outputPath.
// timeout controls the maximum time allowed for the dot process (use 0 for no timeout).
// Returns an error with details if something fails.
func GenerateSVGFromDot(dot string, outputPath string, timeout time.Duration) error {
	// Basic validation
	if strings.TrimSpace(dot) == "" {
		return fmt.Errorf("dot input is empty")
	}
	if strings.TrimSpace(outputPath) == "" {
		return fmt.Errorf("outputPath is empty")
	}

	// Ensure output directory exists
	outDir := filepath.Dir(outputPath)
	if outDir != "." {
		if err := os.MkdirAll(outDir, 0o755); err != nil {
			return fmt.Errorf("failed to create output directory %q: %w", outDir, err)
		}
	}

	// Prepare context with timeout if provided
	ctx := context.Background()
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	// Use dot -Tsvg -o <outputPath> and pass DOT input via stdin
	cmd := exec.CommandContext(ctx, "dot", "-Tsvg", "-o", outputPath)
	cmd.Stdin = strings.NewReader(dot)

	// Capture stderr for debugging
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Run command
	if err := cmd.Run(); err != nil {
		// Check if it was a timeout
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("dot command timed out after %s: %s", timeout.String(), stderr.String())
		}
		return fmt.Errorf("dot command failed: %w: %s", err, stderr.String())
	}

	// Verify file was created
	info, err := os.Stat(outputPath)
	if err != nil {
		return fmt.Errorf("dot reported success but output file missing: %w", err)
	}
	if info.Size() == 0 {
		return fmt.Errorf("output file is empty")
	}

	return nil
}
