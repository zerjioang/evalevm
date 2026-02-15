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

			files, err := scanDirRecursive(path)
			if err != nil {
				return fmt.Errorf("failed to scan directory: %w", err)
			}

			cmp := engine.NewComparator(opts.audit)
			if opts.tools != "" {
				cmp.FilterByTools(strings.Split(opts.tools, ","))
			}
			cmp.Start()

			var renderTasks datatype.TaskSet
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
					renderTasks = append(renderTasks, result)
					if result.Failed() {
						_ = render.ScanError(datatype.ScanErrorDetails{
							Name:    result.ID().App(),
							Message: string(result.Result().OutputErr),
						})
						hasFailed = true
					} else {
						_ = render.ScanSuccess(datatype.ScanSuccess{
							Name:   result.ID().App(),
							Output: result.Result().ParsedOutput.String(),
						})
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
			}

			if err := render.ScanResults(renderTasks); err != nil {
				return err
			}

			// export CSV if requested
			if opts.csvExport {
				if err := exportTaskSetCSV(renderTasks); err != nil {
					return fmt.Errorf("CSV export failed: %w", err)
				}
			}

			return nil
		},
	}

	scanCmd.Flags().StringVarP(&path, "dataset", "d", "", "dataset path")
	bindScanFlags(scanCmd, &opts)
	return scanCmd
}

// scanDirRecursive scans the directory `root` recursively and returns a slice of file paths.
func scanDirRecursive(root string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// If we can't access a file or dir, skip it
			return err
		}

		// Skip directories themselves
		if !d.IsDir() && (filepath.Ext(path) == ".hex" || filepath.Ext(path) == ".evm" || filepath.Ext(path) == ".bin") {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return files, nil
}
