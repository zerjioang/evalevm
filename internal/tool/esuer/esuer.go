package acuarica

import (
	"evalevm/internal/datatype"
)

type Eusuer struct {
	datatype.BytecodeAnalyzer
}

var _ datatype.Analyzer = (*Eusuer)(nil)

func NewEusuer() Eusuer {
	app := Eusuer{}
	app.AppName = "acuarica-evm"
	app.WebsiteUrl = "https://github.com/acuarica/evm"
	app.Desc = `A Symbolic Ethereum Virtual Machine (EVM) interpreter and decompiler, along with several other utils for programmatically extracting information from bytecode.`
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: false,
		ForceSplitRuntime:    false,
	}
	app.Deprecated = true
	app.LastCommit = "6 months ago"
	app.Language = "typescript"
	return app
}

func (scan Eusuer) CreateTask(uid string, bytecode string, filename string) []datatype.Task {
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

func (scan Eusuer) ParseOutput(output *datatype.Result) error {
	return nil
}
