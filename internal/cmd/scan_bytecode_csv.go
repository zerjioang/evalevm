package cmd

import (
	"encoding/csv"
	"evalevm/internal/datatype"
	"evalevm/internal/engine"
	"evalevm/internal/render"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// Row represents one CSV row
type csvRow struct {
	Timestamp string
	Address   string
	Bytecode  string
}

func ScanBytecodeCSVCmd() *cobra.Command {
	type scanBytecodeCSVFlag struct {
		path string
	}

	var flags scanBytecodeCSVFlag

	scanCmd := &cobra.Command{
		Use:   "csv",
		Short: "scan bytecodes from downloaded local CSV",
		RunE: func(cmd *cobra.Command, args []string) error {

			if flags.path == "" {
				return fmt.Errorf("path not provided")
			}

			start := time.Now()
			defer func() {
				elapsed := time.Since(start)
				log.Printf("dir tool benchmark scan completed in %s", elapsed)
			}()

			cmp := engine.NewPaperOnlyComparator()
			cmp.Start()

			file, err := os.Open(flags.path)
			if err != nil {
				log.Fatalf("failed to open CSV: %v", err)
			}
			defer file.Close()

			reader := csv.NewReader(file)

			// Read header first
			header, err := reader.Read()
			if err != nil {
				log.Fatalf("failed to read header: %v", err)
			}
			fmt.Println("CSV header:", header)
			for {
				record, err := reader.Read()
				if err != nil {
					if err.Error() == "EOF" {
						break
					}
					log.Fatalf("error reading CSV: %v", err)
				}

				row := csvRow{
					Timestamp: record[0],
					Address:   record[1],
					Bytecode:  record[2],
				}
				log.Printf("submitting task: %s\n", row.Bytecode)
				taskset := cmp.SubmitAndWait(row.Bytecode)
				forceFail := false
				for _, result := range taskset {
					if result.Failed() {
						_ = render.ScanError(datatype.ScanErrorDetails{
							Name:    result.ID().App(),
							Message: string(result.Result().OutputErr),
						})
						panic("scan failed")
					}
					
					coverage := result.Result().ParsedOutput.Coverage
					if coverage != nil && *coverage != 100 {
						fmt.Println("bytecode: ", row.Bytecode, "")
						fmt.Println("coverage: ", *coverage, "")
						forceFail = true
					}
				}

				if err := render.ScanResults(taskset); err != nil {
					return err
				}

				if forceFail {
					panic("coverage not 100. debug this contract manually")
				}
			}

			return nil
		},
	}

	scanCmd.Flags().StringVarP(&flags.path, "dataset", "d", "", "csv dataset path")
	return scanCmd
}
