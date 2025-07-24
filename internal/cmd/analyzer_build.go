package cmd

import (
	"errors"
	"evalevm/internal/engine"
	"github.com/spf13/cobra"
)

func AnalyzerBuildCmd() *cobra.Command {

	type buildflags struct {
		toolsPath string
	}
	var flags buildflags

	buildCmd := &cobra.Command{
		Use:   "build",
		Short: "build analyzers",
		RunE: func(cmd *cobra.Command, args []string) error {

			if flags.toolsPath == "" {
				return errors.New("tools path not provided")
			}

			cmp := engine.NewComparator()
			tools := cmp.Analyzers()
			for _, analyzer := range tools {
				var b engine.Builder
				if err := b.Build(cmd.Context(), flags.toolsPath, analyzer); err != nil {
					return err
				}
			}
			return nil
		},
	}

	buildCmd.PersistentFlags().StringVarP(&flags.toolsPath, "tools", "t", "", "tools path")

	return buildCmd
}
