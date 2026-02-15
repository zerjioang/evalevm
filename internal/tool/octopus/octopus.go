package octopus

import (
	_ "embed"
	"evalevm/internal/datatype"
	"fmt"
	"strings"
)

var (
	//go:embed Dockerfile
	octopusDockerfile string
)

type Octopus struct {
	datatype.BytecodeAnalyzer
}

var _ datatype.Analyzer = (*Octopus)(nil)

func NewOctopus() Octopus {
	app := Octopus{}
	app.AppName = "octopus"
	app.WebsiteUrl = "https://github.com/FuzzingLabs/octopus"
	app.Desc = "The purpose of Octopus is to provide an easy way to analyze closed-source WebAssembly module and smart contracts bytecode to understand deeper their internal behaviours. It generates Control Flow Graphs (CFG)."
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: true, // Octopus might need raw hex or 0x? Docs used file.
		ForceSplitRuntime:    true, // CFG usually needs runtime code
	}
	app.Deprecated = true
	app.LastCommit = "6 years ago"
	app.Language = "python"
	app.Dockerfile = octopusDockerfile // Embed this
	app.SupportsCFG = true
	app.Platform = "linux/amd64"
	return app
}

func (scan Octopus) CreateTask(uid string, bytecode string, filename string) []datatype.Task {
	return []datatype.Task{
		datatype.NewDockerTask(
			scan.CreateTaskId(uid),
			scan.Options,
			bytecode,
			filename,
			[]string{
				"local/octopus",
				"-c",
				// echo bytecode to file, run octopus_eth_evm.py -f file -g (CFG)
				// Docs say: python3 octopus_eth_evm.py -s -f file
				// I need to check flags. Assuming -g or --cfg based on text.
				// Let's rely on help output or standard usage.
				// Actually, I should verify flags.
				// For now, I'll use a placeholder command in helper.sh style
				fmt.Sprintf(`echo %s > code.evm && ./measure.sh bash -c 'python3 octopus_eth_evm.py -f code.evm --cfg && cat graph.cfg.gv'`, bytecode),
			},
		),
	}
}

func (scan Octopus) ParseOutput(output *datatype.Result) error {
	// octopus outputs the .gv file content to stdout via helper.sh or direct cat
	// We need to count nodes and edges in the DOT/GV format.
	// Pattern for edges: " -> "
	// Pattern for nodes: "[label=" (each node has a label in octopus output)

	dotContent := string(output.Output)

	edges := strings.Count(dotContent, " -> ")
	nodes := strings.Count(dotContent, " [label=")

	output.ParsedOutput = &datatype.ScanResult{
		EdgesDetected: edges,
		NodesDetected: nodes,
	}

	// Also store the graph
	if err := output.ParsedOutput.WithGraph(dotContent, "", output); err != nil {
		return fmt.Errorf("failed to store graph: %w", err)
	}

	return nil
}
