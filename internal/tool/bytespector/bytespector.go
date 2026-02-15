package bytespector

import (
	_ "embed"
	"evalevm/internal/datatype"
	"fmt"
	"strings"
)

var (
	//go:embed Dockerfile
	byteInspectorDockerfile string
)

type ByteSpector struct {
	datatype.BytecodeAnalyzer
}

var _ datatype.Analyzer = (*ByteSpector)(nil)

func NewByteInspector() ByteSpector {
	app := ByteSpector{}
	app.BytecodeAnalyzer = app.SetupDockerPlatform()
	app.AppName = "bytespector"
	app.WebsiteUrl = "https://github.com/franck44/evm-dis"
	app.Desc = `This project provides an EVM bytecode disassembler and Control Flow Graph (CFG) generator. ByteSpector can verify the CFGs by generating a Dafny file that encodes the semantics of the EVM bytecode. The Dafny file can be verified with Dafny. If a CFG is successfully verified, we obtain the following guarantees on the CFG and the bytecode`
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: false,
		ForceSplitRuntime:    false,
	}
	app.Deprecated = false
	app.LastCommit = "e years ago"
	app.Language = "dafny"
	app.Dockerfile = byteInspectorDockerfile
	app.Platform = "linux/amd64"
	app.SupportsVulnerabilities = false
	app.SupportsCFG = true
	return app
}

func (scan ByteSpector) CreateTask(uid string, bytecode string, filename string) []datatype.Task {
	// docker run --rm -it --entrypoint=bash local/byte-inspector -c "echo 0x366028576000600060006000303173f43febf30d4a00fa9b23e49e36e7acb5ca8591616103e8f1005b6388c2a0bf60e060020a026000526000358043116077574390036001016003023562ffffff16600452600060006024600060007306012c8cf97bead5deae237070f9587f8e7a266d6103e85a03f15b00 > code.evm && /tacas25/evm-dis/measure.sh bash -c '/tacas25/evm-dis/makeCFG.sh code.evm && cat build/dot/code.evm/code.evm.dot'"
	return []datatype.Task{
		datatype.NewDockerTask(
			scan.CreateTaskId(uid),
			scan.Options,
			bytecode,
			filename,
			[]string{
				"--platform",
				"linux/amd64",
				"local/bytespector",
				"-c",
				fmt.Sprintf(`echo %s > code.evm && ./measure.sh bash -c 'cd /tacas25/evm-dis/ && ./makeCFG.sh code.evm && cat build/dot/code.evm/code.evm.dot'`, bytecode),
			},
		),
	}
}

func (scan ByteSpector) ParseOutput(output *datatype.Result) error {
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
