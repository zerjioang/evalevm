package teether

import (
	"evalevm/internal/datatype"
)

type Teether struct {
	datatype.BytecodeAnalyzer
}

var _ datatype.Analyzer = (*Teether)(nil)

func NewTeether() Teether {
	app := Teether{}
	app.AppName = "teether"
	app.WebsiteUrl = "https://github.com/nescio007/teether"
	app.Desc = `teEther is an analysis tool for Ethereum smart contracts`
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: false,
		ForceSplitRuntime:    false,
	}
	app.Deprecated = false
	app.LastCommit = "4 years ago"
	app.Language = "python"
	return app
}

func (scan Teether) CreateTask(uid string, bytecode string) []datatype.Task {
	return []datatype.Task{
		datatype.NewDockerTask(
			scan.CreateTaskId(uid),
			scan.Options,
			bytecode,
			[]string{
				"run", "--rm", "--cap-add=SYS_ADMIN", "hello-world",
			},
		),
	}
}

func (scan Teether) ParseOutput(output *datatype.Result) error {
	return nil
}
