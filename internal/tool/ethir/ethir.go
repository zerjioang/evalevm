package ethir

import (
	_ "embed"
	"evalevm/internal/datatype"
	"fmt"
	"strings"
)

var (
	//go:embed Dockerfile
	ethirDockerfile string
)

type EthIR struct {
	datatype.BytecodeAnalyzer
}

var _ datatype.Analyzer = (*EthIR)(nil)

func NewEthIR() EthIR {
	app := EthIR{}
	app.AppName = "ethir"
	app.WebsiteUrl = "https://github.com/costa-group/EthIR"
	app.Desc = "EthIR is a framework for high-level Analysis of Ethereum Bytecode. It generates Control Flow Graphs (CFG)."
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: true, // EthIR example uses -s file.evm -b. I suspect checking code.evm content manually is safer.
		ForceSplitRuntime:    true,
	}
	app.Deprecated = false
	app.LastCommit = "3 weeks ago"
	app.Language = "python"
	app.Dockerfile = ethirDockerfile
	app.SupportsCFG = true
	app.Platform = "linux/arm64"
	return app
}

func (scan EthIR) CreateTask(uid string, bytecode string, filename string) []datatype.Task {
	return []datatype.Task{
		datatype.NewDockerTask(
			scan.CreateTaskId(uid),
			scan.Options,
			bytecode,
			filename,
			[]string{
				"local/ethir",
				"-c",
				// echo bytecode to code.evm. Run: python3 ethir.py -s code.evm -b -cfg normal
				// Output in /tmp/costabs/. I need to find the file and cat it.
				// Likely /tmp/costabs/code.evm.dot or similar.
				// I'll cat everything in /tmp/costabs just in case for now or find .dot file.
				fmt.Sprintf(`echo %s > code.evm && ./measure.sh bash -c 'python3 ethir/ethir.py -s code.evm --bytecode -cfg all && find /tmp/costabs -name "*.dot" -exec cat {} \;'`, bytecode),
			},
		),
	}
}

func (scan EthIR) ParseOutput(output *datatype.Result) error {
	// The output is likely the raw DOT file content + maybe some logs.
	outStr := string(output.Output)
	start := strings.Index(outStr, "digraph")
	if start == -1 {
		return fmt.Errorf("failed to parse output: 'digraph' not found in output")
	}
	// Find the last closing brace
	end := strings.LastIndex(outStr, "}")
	if end == -1 || end < start {
		return fmt.Errorf("failed to parse output: closing brace '}' not found after 'digraph'")
	}
	dotGraph := outStr[start : end+1]

	// fix bugs
	dotGraph = strings.ReplaceAll(dotGraph, "ranking=", "rankdir=")

	output.ParsedOutput = &datatype.ScanResult{}
	if err := output.ParsedOutput.WithGraph(dotGraph, "", output); err != nil {
		return fmt.Errorf("failed to store .dot graph: %w", err)
	}
	return nil
}
