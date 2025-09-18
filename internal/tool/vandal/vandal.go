package vandal

import (
	_ "embed"
	"evalevm/internal/datatype"
	"evalevm/internal/parser"
	"fmt"
	"strings"
)

var (
	//go:embed Dockerfile
	vandalDockerfile string
)

type Vandal struct {
	datatype.BytecodeAnalyzer
}

var _ datatype.Analyzer = (*Vandal)(nil)

func NewVandal() Vandal {
	app := Vandal{}
	app.BytecodeAnalyzer = app.SetupDockerPlatform()
	app.AppName = "vandal"
	app.WebsiteUrl = "https://github.com/usyd-blockchain/vandal"
	app.Desc = `Vandal is a static program analysis framework for Ethereum smart contract bytecode, developed at The University of Sydney. It decompiles an EVM bytecode program to an equivalent intermediate representation that encodes the program's control flow graph. This representation removes all stack operations, thereby exposing data dependencies that are otherwise obscured. This information is then fed, with a Datalog specification, into the Souffle analysis engine for the extraction of program properties.`
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: false,
		ForceSplitRuntime:    false,
	}
	app.Deprecated = false
	app.LastCommit = "5 years ago"
	app.Language = "python"
	app.Platform = "linux/amd64"
	app.Dockerfile = vandalDockerfile
	return app
}

func (scan Vandal) CreateTask(uid string, bytecode string, filename string) []datatype.Task {
	return []datatype.Task{
		datatype.NewDockerTask(
			scan.CreateTaskId(uid),
			scan.Options,
			bytecode,
			filename,
			[]string{
				// "bin/decompile -n -v -g graph.html examples/dao_hack.hex"
				"--platform", "linux/amd64",
				"local/vandal", "-c",
				fmt.Sprintf(`./helper.sh vandal %s`, bytecode),
			},
		),
	}
}

func (scan Vandal) ParseOutput(output *datatype.Result) error {
	dotGraph, err := parser.ExtractBetween(string(output.Output), ">>> graph.dot", "<<<")
	if err != nil {
		return fmt.Errorf("failed to parse .dot output: %w", err)
	}
	//var asPtrBool = func(b bool) *bool { return &b }
	output.ParsedOutput = &datatype.ScanResult{
		Vulnerable:           nil, //asPtrBool(findingsDetected),
		Error:                nil,
		EdgesDetected:        strings.Count(dotGraph, " -> "),
		NodesDetected:        strings.Count(dotGraph, `" [`),
		TxOriginVulnerable:   nil, //asPtrBool(txOriginVulnerable),
		ReEntrancyVulnerable: nil, //asPtrBool(reEntrancyVulnerable),
	}
	if err := output.ParsedOutput.WithGraph(dotGraph, "", output); err != nil {
		return fmt.Errorf("failed to store .dot graph: %w", err)
	}
	return nil
}
