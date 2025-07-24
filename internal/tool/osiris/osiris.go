package osiris

import (
	"evalevm/internal/datatype"
)

type Osiris struct {
	datatype.BytecodeAnalyzer
}

var _ datatype.Analyzer = (*Osiris)(nil)

func NewOsiris() Osiris {
	app := Osiris{}
	app.AppName = "osiris"
	app.WebsiteUrl = "https://github.com/christoftorres/Osiris"
	app.Desc = `A tool to detect integer bugs in Ethereum smart contracts (ACSAC 2018).`
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: false,
		ForceSplitRuntime:    false,
	}
	app.Deprecated = false
	app.LastCommit = "2 years ago"
	app.Language = "python"
	return app
}

func (scan Osiris) CreateTask(uid string, bytecode string) []datatype.Task {
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
