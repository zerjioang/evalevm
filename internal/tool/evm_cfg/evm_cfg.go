package evm_cfg

import (
	_ "embed"
	"evalevm/internal/datatype"
	"fmt"
	"strings"
)

var (
	//go:embed Dockerfile
	evmCFGDockerfile string
)

type EvmCFG struct {
	datatype.BytecodeAnalyzer
}

var _ datatype.Analyzer = (*EvmCFG)(nil)

func NewEvmCFG() EvmCFG {
	app := EvmCFG{}
	app.BytecodeAnalyzer = app.SetupDockerPlatform()
	app.AppName = "evm-cfg"
	app.WebsiteUrl = "https://github.com/plotchy/evm-cfg"
	app.Desc = `Symbolic stack CFG generator for EVM`
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: false,
		ForceSplitRuntime:    false,
	}
	app.Deprecated = false
	app.LastCommit = "1 year ago"
	app.Language = "rust"
	app.Dockerfile = evmCFGDockerfile
	app.SupportsVulnerabilities = false
	app.SupportsCFG = true
	return app
}

func (scan EvmCFG) CreateTask(uid string, bytecode string, filename string) []datatype.Task {
	return []datatype.Task{
		datatype.NewDockerTask(
			scan.CreateTaskId(uid),
			scan.Options,
			bytecode,
			filename,
			[]string{
				// docker run command already defined. customize the flags here
				"local/evm-cfg", "-c",
				fmt.Sprintf(`./measure.sh bash -c '/opt/evm-cfg/evm-cfg %s -o cfg.dot && cat cfg.dot'`, bytecode),
			},
		),
	}
}

func (scan EvmCFG) ParseOutput(output *datatype.Result) error {
	dotGraph := string(output.Output)
	dotGraph = strings.Replace(dotGraph, "Dot file saved to cfg.dot\n", "", 1)

	output.ParsedOutput = &datatype.ScanResult{
		EdgesDetected: strings.Count(dotGraph, "->"),
		NodesDetected: strings.Count(dotGraph, "[ label = "),
		DotGraph:      dotGraph,
	}
	if err := output.ParsedOutput.WithGraph(dotGraph, "", output); err != nil {
		return fmt.Errorf("failed to store .dot graph: %w", err)
	}

	return nil
}
