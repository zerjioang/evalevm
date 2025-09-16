package porosity

import (
	"evalevm/internal/datatype"
)

type Porosity struct {
	datatype.BytecodeAnalyzer
}

var _ datatype.Analyzer = (*Porosity)(nil)

func NewPorosity() Porosity {
	app := Porosity{}
	app.AppName = "porosity"
	app.WebsiteUrl = "https://github.com/msuiche/porosity"
	app.Desc = `Decompiler and Security Analysis tool for Blockchain-based Ethereum Smart-Contracts`
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: false,
		ForceSplitRuntime:    false,
	}
	app.Deprecated = false
	app.LastCommit = "7 years ago"
	app.Language = "c++"
	return app
}

func (scan Porosity) CreateTask(uid string, bytecode string) []datatype.Task {
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

func (scan Porosity) ParseOutput(output *datatype.Result) error {
	return nil
}
