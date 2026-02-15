package rattle

import (
	_ "embed"
	"evalevm/internal/datatype"
	"evalevm/internal/parser"
	"fmt"
	"strings"
)

type Rattle struct {
	datatype.BytecodeAnalyzer
}

var (
	//go:embed Dockerfile
	rattleDockerfile string
)

var _ datatype.Analyzer = (*Rattle)(nil)

func NewRattle() Rattle {
	app := Rattle{}
	app.BytecodeAnalyzer = app.SetupDockerPlatform()
	app.AppName = "rattle"
	app.WebsiteUrl = "https://github.com/crytic/rattle"
	app.Desc = `Rattle is an EVM binary static analysis framework designed to work on deployed smart contracts. Rattle takes EVM byte strings, uses a flow-sensitive analysis to recover the original control flow graph, lifts the control flow graph into an SSA/infinite register form, and optimizes the SSA – removing DUPs, SWAPs, PUSHs, and POPs. The conversion from a stack machine to SSA form removes 60%+ of all EVM instructions and presents a much friendlier interface to those who wish to read the smart contracts they’re interacting with.`
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: false,
		ForceSplitRuntime:    false,
	}
	app.Deprecated = false
	app.LastCommit = "2 years ago"
	app.Language = "python"
	app.Dockerfile = rattleDockerfile
	app.SupportsVulnerabilities = false
	app.SupportsCFG = true
	return app
}

func (scan Rattle) CreateTask(uid string, bytecode string, filename string) []datatype.Task {
	return []datatype.Task{
		datatype.NewDockerTask(
			scan.CreateTaskId(uid),
			scan.Options,
			bytecode,
			filename,
			[]string{
				// docker run command already defined. customize the flags here
				"local/rattle",
				"-c",
				fmt.Sprintf(`./helper.sh rattle %s`, bytecode),
			},
		),
	}
}

func (scan Rattle) ParseOutput(output *datatype.Result) error {
	dotGraph, err := parser.ExtractBetween(string(output.Output), ">>> cfg.dot", "<<<")
	if err != nil {
		return fmt.Errorf("failed to parse .dot output: %w", err)
	}
	//var asPtrBool = func(b bool) *bool { return &b }
	output.ParsedOutput = &datatype.ScanResult{
		Vulnerable:           nil, //asPtrBool(findingsDetected),
		Error:                nil,
		EdgesDetected:        strings.Count(dotGraph, " -> block_"),
		NodesDetected:        strings.Count(dotGraph, ` [label="`),
		TxOriginVulnerable:   nil, //asPtrBool(txOriginVulnerable),
		ReEntrancyVulnerable: nil, //asPtrBool(reEntrancyVulnerable),
	}
	if err := output.ParsedOutput.WithGraph(dotGraph, "", output); err != nil {
		return fmt.Errorf("failed to store .dot graph: %w", err)
	}
	return nil
}
