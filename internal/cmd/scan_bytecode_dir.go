package cmd

import (
	"evalevm/internal/datatype"
	"evalevm/internal/engine"
	"evalevm/internal/render"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func ScanBytecodeDirCmd() *cobra.Command {

	var opts scanOpts
	var path string
	var extensions string

	scanCmd := &cobra.Command{
		Use:     "dir",
		Aliases: []string{"d", "dataset", "data", "directory"},
		Short:   "scan bytecodes in directory",
		RunE: func(cmd *cobra.Command, args []string) error {

			if path == "" {
				return fmt.Errorf("path not provided")
			}

			start := time.Now()
			defer func() {
				elapsed := time.Since(start)
				log.Printf("dir tool benchmark scan completed in %s", elapsed)
			}()

			files, err := scanDirRecursive(path, extensions)
			if err != nil {
				return fmt.Errorf("failed to scan directory: %w", err)
			}

			cmp := engine.NewComparator(opts.audit, opts.runMode)
			if opts.tools != "" {
				cmp.FilterByTools(strings.Split(opts.tools, ","))
			}
			cmp.Start(cmd.Context())

			// Initialize streaming CSV writer if requested
			var streamWriter *ResultStreamWriter
			if opts.csvExport {
				streamWriter = NewResultStreamWriter()
				defer streamWriter.Close()
			}

			// Removed renderTasks slice to prevent OOM
			for i, file := range files {
				log.Printf("submitting task %d/%d: %s", i+1, len(files), file)
				content, err := os.ReadFile(file)
				if err != nil {
					return fmt.Errorf("failed to read file: %w", err)
				}
				filebase := filepath.Base(file)
				taskset := cmp.SubmitAndWait(string(content), filebase)

				hasFailed := false
				for _, result := range taskset {
					if result.Failed() {
						_ = render.ScanError(datatype.ScanErrorDetails{
							Name:    result.ID().App(),
							Message: string(result.Result().OutputErr),
						})
						hasFailed = true
					} else {
						// Optional: Log success instead of rendering table
						// log.Printf("scanned %s with %s", filebase, result.ID().App())
					}
				}

				// --stop-on-fail: abort batch on first failure
				if opts.stopOnFail && hasFailed {
					log.Printf("stopping on failure for file: %s", file)
					break
				}

				// --coverage: abort if coverage < 100%
				if opts.coverage {
					for _, result := range taskset {
						if !result.Failed() && result.Result().ParsedOutput != nil {
							cov := result.Result().ParsedOutput.Coverage
							if cov != nil && *cov < 100 {
								return fmt.Errorf("coverage check failed: %s reported %.2f%% coverage (< 100%%)", result.ID().App(), *cov)
							}
						}
					}
				}

				// Stream results to CSV
				if streamWriter != nil {
					if err := streamWriter.Write(taskset); err != nil {
						return fmt.Errorf("failed to stream results to CSV: %w", err)
					}
				}
			}

			// Omitted render.ScanResults to prevent terminal flooding and reduce memory usage

			return nil
		},
	}

	scanCmd.Flags().StringVarP(&path, "dataset", "d", "", "dataset path")
	scanCmd.Flags().StringVar(&extensions, "extensions", "hex,evm,bin", "file extensions to scan (comma separated, e.g. hex,evm)")
	bindScanFlags(scanCmd, &opts)
	return scanCmd
}

// scanDirRecursive scans the directory `root` recursively and returns a slice of file paths.
func scanDirRecursive(root string, extensions string) ([]string, error) {
	var files []string

	allowedExts := make(map[string]bool)
	if extensions != "" {
		for _, ext := range strings.Split(extensions, ",") {
			// Ensure extension starts with dot for comparison
			ext = strings.TrimSpace(ext)
			if !strings.HasPrefix(ext, ".") {
				ext = "." + ext
			}
			allowedExts[strings.ToLower(ext)] = true
		}
	} else {
		// Default extensions
		allowedExts[".hex"] = true
		allowedExts[".evm"] = true
		allowedExts[".bin"] = true
	}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// If we can't access a file or dir, skip it
			return err
		}

		if !d.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			if allowedExts[ext] {
				files = append(files, path)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return files, nil
}
