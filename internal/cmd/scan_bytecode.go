package cmd

import (
	"evalevm/internal/datatype"
	"evalevm/internal/engine"
	"evalevm/internal/render"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func ScanBytecodeCmd() *cobra.Command {

	var opts scanOpts
	var bytecodeHex string

	scanCmd := &cobra.Command{
		Use:   "evm",
		Short: "scan single evm bytecode",
		RunE: func(cmd *cobra.Command, args []string) error {

			if bytecodeHex == "" || bytecodeHex == "0x" {
				return nil
			}

			start := time.Now()
			defer func() {
				elapsed := time.Since(start)
				log.Printf("tool benchmark completed in %s", elapsed)
			}()

			cmp := engine.NewComparator(opts.audit, opts.runMode)
			if opts.tools != "" {
				cmp.FilterByTools(strings.Split(opts.tools, ","))
			}
			cmp.Start(cmd.Context())

			taskset := cmp.SubmitAndWait(bytecodeHex, "")

			log.Println("all tools evaluated and benchmark completed for the evm bytecode sample. exporting results")

			if err := render.ScanResults(taskset, opts.transpose); err != nil {
				return err
			}

			if false {
				// render success results
				for _, result := range taskset {
					if !result.Failed() {
						_ = render.ScanSuccess(datatype.ScanSuccess{
							Name:   result.ID().App(),
							Output: result.Result().ParsedOutput.String(),
						})
					}
				}

				// render error results
				for _, result := range taskset {
					if result.Failed() {
						_ = render.ScanError(datatype.ScanErrorDetails{
							Name:    result.ID().App(),
							Message: string(result.Result().OutputErr),
						})
					}
				}
			}

			// check coverage if requested
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

			// export CSV if requested
			if opts.csvExport {
				streamWriter := NewResultStreamWriter()
				defer streamWriter.Close()
				if err := streamWriter.Write(taskset); err != nil {
					return fmt.Errorf("CSV export failed: %w", err)
				}
			}

			// Print helper command to open generated SVGs
			if len(taskset) > 0 {
				id := taskset[0].Result().Task.TrackerId()
				fmt.Printf("\nTo view generated graphs:\nopen output/%s/*.svg\n", id)
			}

			return nil
		},
	}

	scanCmd.Flags().StringVarP(&bytecodeHex, "bytecode", "b", "", "bytecode in hex")
	bindScanFlags(scanCmd, &opts)
	return scanCmd
}
