package mythril

import (
	"evalevm/internal/datatype"
)

type Mythril struct {
	datatype.BytecodeAnalyzer
}

var _ datatype.Analyzer = (*Mythril)(nil)

func NewMythril() Mythril {
	app := Mythril{}
	app.BytecodeAnalyzer = app.SetupDockerPlatform()
	app.AppName = "mythril"
	app.WebsiteUrl = "https://github.com/ConsenSysDiligence/mythril"
	app.Desc = "Mythril is a symbolic-execution-based securty analysis tool for EVM bytecode. It detects security vulnerabilities in smart contracts built for Ethereum and other EVM-compatible blockchains."
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: false,
		ForceSplitRuntime:    false,
	}
	app.Deprecated = false
	app.LastCommit = "6 months ago"
	app.Language = "python"
	return app
}

func (scan Mythril) CreateTask(uid string, bytecode string) []datatype.Task {
	return []datatype.Task{
		datatype.NewDockerTask(
			scan.CreateTaskId(uid),
			scan.Options,
			bytecode,
			[]string{
				"hello-world",
			},
		),
	}
}

func (scan Mythril) ParseOutput(output *datatype.Result) error {
	return nil
}
