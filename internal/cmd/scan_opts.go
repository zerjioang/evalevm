package cmd

import (
	"github.com/spf13/cobra"
)

// scanOpts holds the shared flags for all scan subcommands.
type scanOpts struct {
	audit      bool   // --audit: enable vulnerability auditing (checkers)
	coverage   bool   // --coverage: stop if any result coverage < 100%
	stopOnFail bool   // --stop-on-fail: stop batch processing on first failure
	tools      string // --tools: comma-separated list of tool names to run
	csvExport  bool   // --csv: export results as per-tool CSV files
	runMode    string // --run-mode: execution mode for tools supporting it (e.g. ethersolve)
	transpose  bool   // --transpose: transpose results table (tools as columns)
}

// bindScanFlags registers the shared scan flags on the given cobra command.
func bindScanFlags(cmd *cobra.Command, opts *scanOpts) {
	cmd.Flags().BoolVar(&opts.audit, "audit", false, "enable vulnerability auditing (checkers)")
	cmd.Flags().BoolVar(&opts.coverage, "coverage", false, "stop execution when coverage < 100%")
	cmd.Flags().BoolVar(&opts.stopOnFail, "stop-on-fail", false, "stop batch processing on first analysis failure")
	cmd.Flags().StringVar(&opts.tools, "tools", "", "comma-separated list of tool names to run (e.g. paper,vandal)")
	cmd.Flags().BoolVar(&opts.csvExport, "csv", false, "export results as per-tool CSV files (e.g. evalevm_paper.csv)")
	cmd.Flags().StringVar(&opts.runMode, "run-mode", "any", "execution mode for tools supporting it (runtime, creator, any)")
	cmd.Flags().BoolVar(&opts.transpose, "transpose", true, "transpose results table (tools as columns)")
}
