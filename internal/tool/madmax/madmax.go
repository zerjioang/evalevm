package madmax

import (
	"evalevm/internal/datatype"
)

type MadMax struct {
	datatype.BytecodeAnalyzer
}

var _ datatype.Analyzer = (*MadMax)(nil)

func NewMadMax() MadMax {
	app := MadMax{}
	app.AppName = "madmax"
	app.WebsiteUrl = "https://github.com/nevillegrech/MadMax"
	app.Desc = `Ethereum Static Vulnerability Detector for Gas-Focussed Vulnerabilities`
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: false,
		ForceSplitRuntime:    false,
	}
	app.Deprecated = false
	app.LastCommit = "3 years ago"
	app.Language = "python"
	return app
}

func (scan MadMax) CreateTask(uid string, bytecode string) []datatype.Task {
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
