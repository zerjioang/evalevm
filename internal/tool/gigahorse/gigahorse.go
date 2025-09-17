package gigahorse

import (
	"evalevm/internal/datatype"
)

type GigaHorse struct {
	datatype.BytecodeAnalyzer
}

var _ datatype.Analyzer = (*GigaHorse)(nil)

func NewGigaHorse() GigaHorse {
	app := GigaHorse{}
	app.AppName = "gigahorse"
	app.WebsiteUrl = "https://github.com/nevillegrech/gigahorse-toolchain"
	app.Desc = `A binary lifter and analysis framework for Ethereum smart contracts`
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: false,
		ForceSplitRuntime:    false,
	}
	app.Deprecated = false
	app.LastCommit = "3 weeks ago"
	app.Language = "python"
	return app
}

func (scan GigaHorse) CreateTask(uid string, bytecode string) []datatype.Task {
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

func (scan GigaHorse) ParseOutput(output *datatype.Result) error {
	return nil
}
