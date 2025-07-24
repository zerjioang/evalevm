package cmd

import (
	"evalevm/internal/engine"
	"github.com/spf13/cobra"
	"log"
	"time"
)

func ScanBytecodeDirCmd() *cobra.Command {

	type scanBytecodeDirFlag struct {
		path string
	}

	var flags scanBytecodeDirFlag

	scanCmd := &cobra.Command{
		Use:   "dir",
		Short: "scan bytecodes in directory",
		RunE: func(cmd *cobra.Command, args []string) error {

			start := time.Now()
			defer func() {
				elapsed := time.Since(start)
				log.Printf("dir tool benchmark scan completed in %s", elapsed)
			}()

			cmp := engine.NewComparator()
			cmp.Start()

			taskset := cmp.SubmitAndWait("0x0")

			log.Println("all task completed. collecting results")
			if len(taskset) == 0 {
				log.Fatal("failed to submit task")
			}

			return nil
		},
	}

	scanCmd.Flags().StringVarP(&flags.path, "path", "p", "", "path")
	return scanCmd
}
