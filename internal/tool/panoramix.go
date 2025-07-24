package tool

import (
	"evalevm/internal/datatype"
)

type Panoramix struct {
	datatype.BytecodeAnalyzer
}

var _ datatype.Analyzer = (*Panoramix)(nil)

func NewPanoramix() Panoramix {
	app := Panoramix{}
	app.AppName = "panoramix"
	app.WebsiteUrl = "https://github.com/eveem-org/panoramix"
	app.Desc = `Decompiler at the heart of Eveem.org`
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: false,
		ForceSplitRuntime:    false,
	}
	app.Deprecated = true
	app.LastCommit = "5 years ago"
	app.Language = "python"
	return app
}

func (scan Panoramix) CreateTask(uid string, bytecode string) []datatype.Task {
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
