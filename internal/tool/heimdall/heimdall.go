package heimdall

import (
	"evalevm/internal/datatype"
)

type Heimdal struct {
	datatype.BytecodeAnalyzer
}

var _ datatype.Analyzer = (*Heimdal)(nil)

func NewHeimdall() Heimdal {
	app := Heimdal{}
	app.AppName = "heimdall-rs"
	app.WebsiteUrl = "https://github.com/Jon-Becker/heimdall-rs"
	app.Desc = `Heimdall is an advanced EVM smart contract toolkit specializing in bytecode analysis and extracting information from unverified contracts.`
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: false,
		ForceSplitRuntime:    false,
	}
	app.Deprecated = false
	app.LastCommit = "5 days ago"
	app.Language = "rust"
	return app
}

func (scan Heimdal) CreateTask(uid string, bytecode string, filename string) []datatype.Task {
	return []datatype.Task{
		datatype.NewDockerTask(
			scan.CreateTaskId(uid),
			scan.Options,
			bytecode,
			filename,
			[]string{
				// docker run command already defined. customize the flags here
				"hello-world",
			},
		),
	}
}

func (scan Heimdal) ParseOutput(output *datatype.Result) error {
	return nil
}
