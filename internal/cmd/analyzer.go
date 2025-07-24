package cmd

import (
	"github.com/spf13/cobra"
)

func AnalyzerCmd() *cobra.Command {
	analyzerCmd := &cobra.Command{
		Use:   "analyzer",
		Short: "manage analyzers",
	}

	analyzerCmd.AddCommand(AnalyzerListCmd())
	analyzerCmd.AddCommand(AnalyzerBuildCmd())

	return analyzerCmd
}
