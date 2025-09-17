package manticore

import (
	"evalevm/internal/datatype"
)

type Manticore struct {
	datatype.BytecodeAnalyzer
}

var _ datatype.Analyzer = (*Manticore)(nil)

func NewManticore() Manticore {
	app := Manticore{}
	app.AppName = "manticore"
	app.WebsiteUrl = "https://github.com/trailofbits/manticore"
	app.Desc = `Manticore is a symbolic execution tool for the analysis of Ethereum smart contracts and binaries.`
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: false,
		ForceSplitRuntime:    false,
	}
	app.Deprecated = true
	app.LastCommit = "2 years ago"
	app.Language = "python"
	return app
}

func (scan Manticore) CreateTask(uid string, bytecode string) []datatype.Task {
	return []datatype.Task{
		datatype.NewDockerTask(
			scan.CreateTaskId(uid),
			scan.Options,
			bytecode,
			[]string{
				// docker run command already defined. customize the flags here
				"hello-world",
			},
		),
	}
}

func (scan Manticore) ParseOutput(output *datatype.Result) error {
	return nil
}
