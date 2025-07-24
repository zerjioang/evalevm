package pakala

import (
	"evalevm/internal/datatype"
)

type Pakala struct {
	datatype.BytecodeAnalyzer
}

var _ datatype.Analyzer = (*Pakala)(nil)

func NewPakala() Pakala {
	app := Pakala{}
	app.AppName = "pakala"
	app.WebsiteUrl = "https://github.com/palkeo/pakala"
	app.Desc = `Offensive vulnerability scanner for ethereum, and symbolic execution tool for the Ethereum Virtual Machine`
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: false,
		ForceSplitRuntime:    false,
	}
	app.Deprecated = false
	app.LastCommit = "5 years ago"
	app.Language = "python"
	app.PaperURL = "https://www.palkeo.com/en/projets/ethereum/pakala.html"
	return app
}

func (scan Pakala) CreateTask(uid string, bytecode string) []datatype.Task {
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
