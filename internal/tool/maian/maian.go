package maian

import (
	"evalevm/internal/datatype"
)

type Maian struct {
	datatype.BytecodeAnalyzer
}

var _ datatype.Analyzer = (*Maian)(nil)

func NewMaian() Maian {
	app := Maian{}
	app.AppName = "maian"
	app.WebsiteUrl = "https://github.com/ivicanikolicsg/MAIAN"
	app.Desc = `MAIAN: automatic tool for finding trace vulnerabilities in Ethereum smart contracts`
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: false,
		ForceSplitRuntime:    false,
	}
	app.Deprecated = false
	app.LastCommit = "4 years ago"
	app.Language = "python"
	return app
}

func (scan Maian) CreateTask(uid string, bytecode string) []datatype.Task {
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
