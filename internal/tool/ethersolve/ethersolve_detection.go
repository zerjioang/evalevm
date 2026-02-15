package ethersolve

import (
	_ "embed"
	"evalevm/internal/datatype"
	"fmt"
)

type EthersolveDetection struct {
	datatype.BytecodeAnalyzer
}

var _ datatype.Analyzer = (*EthersolveDetection)(nil)

func NewEthersolveDetection() EthersolveDetection {
	app := EthersolveDetection{}
	app.BytecodeAnalyzer = app.SetupDockerPlatform()
	app.AppName = "ethersolve_detection"
	app.WebsiteUrl = "https://github.com/SeUniVr/EtherSolve"
	app.Desc = `EtherSolve is a tool for Control Flow Graph (CFG) reconstruction and static analysis of Solidity smart-contracts from Ethereum bytecode.`
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: true,
		ForceSplitRuntime:    true,
	}
	app.Deprecated = false
	app.LastCommit = "2 years ago"
	app.Language = "java"
	app.Dockerfile = ethersolveDockerfile
	app.SupportsVulnerabilities = true
	app.SupportsCFG = true
	return app
}

func (scan EthersolveDetection) CreateTask(uid string, bytecode string, filename string) []datatype.Task {
	// Reusing the same flags as runtime, as they enable detection
	return []datatype.Task{
		datatype.NewDockerTask(
			scan.CreateTaskId(uid),
			scan.Options,
			bytecode,
			filename,
			[]string{
				// docker run command already defined. customize the flags here
				"local/ethersolve_detection",
				"-c",
				fmt.Sprintf(`./helper.sh ethersolve_runtime %s`, bytecode),
				// output order: re-entrancy, tx-origin, dot file
			},
		),
	}
}

func (scan EthersolveDetection) ParseOutput(output *datatype.Result) error {
	return parseEthersolveOutput(output)
}
