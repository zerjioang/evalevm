package cmd

import (
	"evalevm/internal/engine"
	"evalevm/internal/render"

	"github.com/spf13/cobra"
)

func AnalyzerListCmd() *cobra.Command {
	analyzerListCmd := &cobra.Command{
		Use:   "list",
		Short: "list analyzers",
		Run: func(cmd *cobra.Command, args []string) {
			cmp := engine.NewComparator(false)
			render.RenderAnalyzers(cmp.Analyzers())
		},
	}
	return analyzerListCmd
}
