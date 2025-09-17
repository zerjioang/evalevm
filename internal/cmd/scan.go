package cmd

import (
	"github.com/spf13/cobra"
)

func ScanCmd() *cobra.Command {
	scanCmd := &cobra.Command{
		Use:   "scan",
		Short: "manage scan",
	}

	scanCmd.AddCommand(ScanBytecodeCmd())
	scanCmd.AddCommand(ScanBytecodeCSVCmd())
	scanCmd.AddCommand(ScanBytecodeDirCmd())

	return scanCmd
}
