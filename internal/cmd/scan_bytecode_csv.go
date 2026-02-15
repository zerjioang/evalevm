package cmd

import (
	"encoding/csv"
	"errors"
	"evalevm/internal/datatype"
	"evalevm/internal/engine"
	"evalevm/internal/render"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// csvRow represents one CSV row
type csvRow struct {
	Timestamp string
	Address   string
	Bytecode  string
}

func ScanBytecodeCSVCmd() *cobra.Command {
	var opts scanOpts
	var path string

	scanCmd := &cobra.Command{
		Use:   "csv",
		Short: "scan bytecodes from downloaded local CSV",
		RunE: func(cmd *cobra.Command, args []string) error {

			if path == "" {
				return fmt.Errorf("path not provided")
			}

			start := time.Now()
			defer func() {
				elapsed := time.Since(start)
				log.Printf("csv tool benchmark scan completed in %s", elapsed)
			}()

			cmp := engine.NewComparator(opts.audit, opts.runMode)
			if opts.tools != "" {
				cmp.FilterByTools(strings.Split(opts.tools, ","))
			}
			cmp.Start()

			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open CSV: %w", err)
			}
			defer file.Close()

			reader := csv.NewReader(file)

			// Read header first
			header, err := reader.Read()
			if err != nil {
				return fmt.Errorf("failed to read header: %w", err)
			}
			log.Println("CSV header:", header)

			var allTasks datatype.TaskSet
			for {
				record, err := reader.Read()
				if err != nil {
					if errors.Is(err, io.EOF) {
						break
					}
					return fmt.Errorf("error reading CSV: %w", err)
				}

				// Skip comments
				if len(record) == 0 || record[0] == "" || record[0][0] == '#' || (len(record[0]) > 1 && record[0][:2] == "//") {
					continue
				}

				row := csvRow{
					Timestamp: record[0],
					Address:   record[1],
					Bytecode:  record[2],
				}
				log.Printf("submitting task: %s\n", row.Bytecode)
				taskset := cmp.SubmitAndWait(row.Bytecode, row.Address)

				hasFailed := false
				for _, result := range taskset {
					allTasks = append(allTasks, result)
					if result.Failed() {
						_ = render.ScanError(datatype.ScanErrorDetails{
							Name:    result.ID().App(),
							Message: string(result.Result().OutputErr),
						})
						hasFailed = true
					}
				}

				// --stop-on-fail: abort batch on first failure
				if opts.stopOnFail && hasFailed {
					return fmt.Errorf("scan failed for contract %s, stopping", row.Address)
				}

				// --coverage: abort if coverage < 100%
				if opts.coverage {
					for _, result := range taskset {
						if !result.Failed() && result.Result().ParsedOutput != nil {
							cov := result.Result().ParsedOutput.Coverage
							if cov != nil && *cov < 100 {
								return fmt.Errorf("coverage check failed: %s reported %.2f%% coverage (< 100%%) for contract %s",
									result.ID().App(), *cov, row.Address)
							}
						}
					}
				}

				if err := render.ScanResults(taskset, opts.transpose); err != nil {
					return err
				}
			}

			// export CSV if requested
			if opts.csvExport {
				if err := exportTaskSetCSV(allTasks); err != nil {
					return fmt.Errorf("CSV export failed: %w", err)
				}
			}

			return nil
		},
	}

	scanCmd.Flags().StringVarP(&path, "dataset", "d", "", "csv dataset path")
	bindScanFlags(scanCmd, &opts)
	return scanCmd
}
