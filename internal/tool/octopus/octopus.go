package octopus

import (
	"evalevm/internal/datatype"
)

type Octopus struct {
	datatype.BytecodeAnalyzer
}

var _ datatype.Analyzer = (*Octopus)(nil)

func NewOctopus() Octopus {
	app := Octopus{}
	app.AppName = "octopus"
	app.WebsiteUrl = "https://github.com/FuzzingLabs/octopus"
	app.Desc = "The purpose of Octopus is to provide an easy way to analyze closed-source WebAssembly module and smart contracts bytecode to understand deeper their internal behaviours"
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: false,
		ForceSplitRuntime:    false,
	}
	app.Deprecated = true
	app.LastCommit = "5 years ago"
	app.Language = "python"
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
				"run", "--rm", "--cap-add=SYS_ADMIN", "hello-world",
			},
		),
	}
}

func (scan Octopus) ParseOutput(output *datatype.Result) error {
	return nil
}
