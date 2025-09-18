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
	"time"

	"github.com/spf13/cobra"
)

func ScanBytecodeDirCmd() *cobra.Command {

	type scanBytecodeDirFlag struct {
		path string
	}

	var flags scanBytecodeDirFlag

	scanCmd := &cobra.Command{
		Use:     "dir",
		Aliases: []string{"d", "dataset", "data", "directory"},
		Short:   "scan bytecodes in directory",
		RunE: func(cmd *cobra.Command, args []string) error {

			if flags.path == "" {
				return fmt.Errorf("path not provided")
			}

			start := time.Now()
			defer func() {
				elapsed := time.Since(start)
				log.Printf("dir tool benchmark scan completed in %s", elapsed)
			}()

			files, err := scanDirRecursive(flags.path)
			if err != nil {
				return fmt.Errorf("failed to scan directory: %w", err)
			}

			cmp := engine.NewComparator()
			cmp.Start()

			var alltasks []datatype.TaskSet
			for i, file := range files {
				log.Printf("submitting task %d/%d: %s", i+1, len(files), file)
				content, err := os.ReadFile(file)
				if err != nil {
					return fmt.Errorf("failed to read file: %w", err)
				}
				filebase := filepath.Base(file)
				taskset := cmp.SubmitAndWait(string(content), filebase)
				alltasks = append(alltasks, taskset)
			}

			log.Println("all task completed. collecting results")
			var renderTasks datatype.TaskSet
			for _, taskset := range alltasks {
				// also render the scanners with success
				for _, result := range taskset {
					if !result.Failed() {
						_ = render.ScanSuccess(datatype.ScanSuccess{
							Name:   result.ID().App(),
							Output: result.Result().ParsedOutput.String(),
						})
					}
					renderTasks = append(renderTasks, result)
				}

				// also render the scanners with errors
				for _, result := range taskset {
					if result.Failed() {
						_ = render.ScanError(datatype.ScanErrorDetails{
							Name:    result.ID().App(),
							Message: string(result.Result().OutputErr),
						})
					}
				}
			}

			if err := render.ScanResults(renderTasks); err != nil {
				return err
			}

			return nil
		},
	}

	scanCmd.Flags().StringVarP(&flags.path, "dataset", "d", "", "dataset path")
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
