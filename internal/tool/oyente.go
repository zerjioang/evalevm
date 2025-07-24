package tool

import (
	"evalevm/internal/datatype"
)

type Oyente struct {
	datatype.BytecodeAnalyzer
}

var _ datatype.Analyzer = (*Oyente)(nil)

func NewOyente() Oyente {
	app := Oyente{}
	app.AppName = "oyente"
	app.WebsiteUrl = "https://github.com/enzymefinance/oyente"
	app.Desc = ``
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: false,
		ForceSplitRuntime:    false,
	}
	app.Deprecated = false
	app.LastCommit = "5 years ago"
	app.Language = "python"
	return app
}

func (scan Oyente) CreateTask(uid string, bytecode string) []datatype.Task {
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
