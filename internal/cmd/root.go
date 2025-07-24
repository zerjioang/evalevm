package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

func Root(ctx context.Context) *cobra.Command {
	var cfgFile string
	cobra.OnInitialize(func() {
		initConfig(cfgFile)
	})

	type rootflags struct {
		verbose bool
	}
	var flags rootflags

	// rootCmd represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use:     "eval",
		Short:   "evaluate evm security",
		Long:    "compare the output of different evm security tools",
		Version: "1.0",
		PreRun: func(cmd *cobra.Command, args []string) {
			cmd.SetContext(ctx)
		},
	}

	rootCmd.AddCommand(AnalyzerCmd())
	rootCmd.AddCommand(ScanCmd())
	rootCmd.PersistentFlags().BoolVarP(&flags.verbose, "verbose", "v", false, "verbose output")

	return rootCmd
}

// initConfig reads in config file and ENV variables if set.
func initConfig(cfgFile string) {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}

		// Search config in home directory with name ".myapp" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".eval")
	}

	// Read in environment variables that match
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintf(os.Stderr, "using config file: %s\n", viper.ConfigFileUsed())
	}
}
