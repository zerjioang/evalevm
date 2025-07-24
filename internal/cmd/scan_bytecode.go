package cmd

import (
	"evalevm/internal/datatype"
	"evalevm/internal/engine"
	"evalevm/internal/export"
	"evalevm/internal/render"
	"github.com/spf13/cobra"
	"log"
	"time"
)

func ScanBytecodeCmd() *cobra.Command {

	type scanBytecodeFlag struct {
		bytecodeHex string
	}

	var flags scanBytecodeFlag

	scanCmd := &cobra.Command{
		Use:   "evm",
		Short: "scan single evm bytecode",
		RunE: func(cmd *cobra.Command, args []string) error {

			if flags.bytecodeHex == "" || flags.bytecodeHex == "0x" {
				return nil
			}

			start := time.Now()
			defer func() {
				elapsed := time.Since(start)
				log.Printf("tool benchmark completed in %s", elapsed)
			}()

			cmp := engine.NewComparator()
			cmp.Start()

			taskset := cmp.SubmitAndWait(flags.bytecodeHex)

			log.Println("all tools evaluated and bechmark completed for the evm bytecode sample. exporting results")
			export.WriteCSV(taskset)

			if err := render.ScanResults(taskset); err != nil {
				return err
			}

			// also render the scanners with success
			for _, result := range taskset {
				if !result.Failed() {
					_ = render.ScanSuccess(datatype.ScanSuccess{
						Name:   result.ID().App(),
						Output: result.Result().Output,
					})
				}
			}

			// also render the scanners with errors
			for _, result := range taskset {
				if result.Failed() {
					_ = render.ScanError(datatype.ScanErrorDetails{
						Name:    result.ID().App(),
						Message: result.Result().Output,
					})
				}
			}

			if err := render.ScanResults(taskset); err != nil {
				return err
			}

			return nil
		},
	}

	scanCmd.Flags().StringVarP(&flags.bytecodeHex, "bytecode", "b", "", "bytecode in hex")
	return scanCmd
}
