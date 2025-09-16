package pyevmasm

import (
	"evalevm/internal/datatype"
)

type Pyevmasm struct {
	datatype.BytecodeAnalyzer
}

var _ datatype.Analyzer = (*Pyevmasm)(nil)

func NewPyevmasm() Pyevmasm {
	app := Pyevmasm{}
	app.AppName = "pyevmasm"
	app.WebsiteUrl = "https://github.com/crytic/pyevmasm"
	app.Desc = `pyevmasm is an assembler and disassembler library for the Ethereum Virtual Machine (EVM). It includes a commandline utility and a Python API.`
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: false,
		ForceSplitRuntime:    false,
	}
	app.Deprecated = false
	app.LastCommit = "2 months ago"
	app.Language = "python"
	return app
}

func (scan Pyevmasm) CreateTask(uid string, bytecode string) []datatype.Task {
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

func (scan Pyevmasm) ParseOutput(output *datatype.Result) error {
	return nil
}
